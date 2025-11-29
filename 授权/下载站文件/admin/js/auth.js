(function (window) {
    const AUTH_API = '../api/admin/auth.php';
    const DASHBOARD_PAGE = 'dashboard.html';
    const LOGIN_PAGE = 'login.html';
    const TOKEN_KEY = 'adminAuthToken';
    const USER_KEY = 'adminUsername';
    const EXPIRE_KEY = 'adminTokenExpire';
    const VERIFY_CACHE_KEY = 'adminTokenVerifiedAt';

    const AdminAuth = {
        async login(username, password) {
            const response = await fetch(`${AUTH_API}?action=login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            const result = await response.json().catch(() => ({ message: '登录失败' }));
            if (!response.ok || result.code !== 200) {
                throw new Error(result.message || '登录失败');
            }

            this.setSession(result.data || {});
            return result;
        },

        async logout(redirect = true) {
            const token = this.getToken();
            this.clear();
            if (token) {
                try {
                    await fetch(`${AUTH_API}?action=logout`, {
                        method: 'POST',
                        headers: {
                            'Authorization': 'Bearer ' + token
                        }
                    });
                } catch (e) {
                    // ignore network error when logging out
                }
            }
            if (redirect) {
                this.redirectToLogin();
            }
        },

        getToken() {
            return localStorage.getItem(TOKEN_KEY) || '';
        },

        getUsername() {
            return localStorage.getItem(USER_KEY) || 'admin';
        },

        getExpire() {
            return parseInt(localStorage.getItem(EXPIRE_KEY) || '0', 10);
        },

        setSession(data) {
            if (data.token) {
                localStorage.setItem(TOKEN_KEY, data.token);
            }
            if (data.username) {
                localStorage.setItem(USER_KEY, data.username);
            }

            const expireAt = data.expire_at ? Date.parse(data.expire_at) : 0;
            if (expireAt && !Number.isNaN(expireAt)) {
                localStorage.setItem(EXPIRE_KEY, expireAt.toString());
            } else {
                localStorage.removeItem(EXPIRE_KEY);
            }
            sessionStorage.setItem(VERIFY_CACHE_KEY, Date.now().toString());
            this.applyAxiosToken();
        },

        clear() {
            localStorage.removeItem(TOKEN_KEY);
            localStorage.removeItem(USER_KEY);
            localStorage.removeItem(EXPIRE_KEY);
            sessionStorage.removeItem(VERIFY_CACHE_KEY);
            this.applyAxiosToken();
        },

        hasValidSession() {
            const token = this.getToken();
            if (!token) {
                return false;
            }
            const expire = this.getExpire();
            if (!expire) {
                return true;
            }
            return Date.now() < expire;
        },

        async verifyToken(options = {}) {
            const token = this.getToken();
            if (!token) {
                throw new Error('missing token');
            }

            const cachedAt = parseInt(sessionStorage.getItem(VERIFY_CACHE_KEY) || '0', 10);
            if (!options.force && cachedAt && Date.now() - cachedAt < 60000) {
                return true;
            }

            const response = await fetch(`${AUTH_API}?action=verify`, {
                headers: {
                    'Authorization': 'Bearer ' + token
                }
            });

            if (!response.ok) {
                throw new Error('verify failed');
            }

            const result = await response.json().catch(() => ({ data: {} }));
            sessionStorage.setItem(VERIFY_CACHE_KEY, Date.now().toString());
            if (result && result.data) {
                if (result.data.username) {
                    localStorage.setItem(USER_KEY, result.data.username);
                }
                if (result.data.expire_at) {
                    const ts = Date.parse(result.data.expire_at);
                    if (!Number.isNaN(ts)) {
                        localStorage.setItem(EXPIRE_KEY, ts.toString());
                    }
                }
            }
            this.applyAxiosToken();
            return true;
        },

        requireAuth() {
            if (!this.hasValidSession()) {
                this.handleUnauthorized();
                return false;
            }
            this.verifyToken().catch(() => this.handleUnauthorized());
            return true;
        },

        redirectToLogin() {
            if (window.location.pathname.endsWith(`/${LOGIN_PAGE}`)) {
                return;
            }
            window.location.href = LOGIN_PAGE;
        },

        redirectToDashboard() {
            if (window.location.pathname.endsWith(`/${DASHBOARD_PAGE}`)) {
                return;
            }
            window.location.href = DASHBOARD_PAGE;
        },

        handleUnauthorized() {
            this.clear();
            const body = document.body;
            if (body && body.dataset && body.dataset.authPage === 'login') {
                return;
            }
            this.redirectToLogin();
        },

        authFetch(url, options = {}) {
            const token = this.getToken();
            const config = Object.assign({}, options);
            config.headers = Object.assign({}, options.headers || {});
            if (token) {
                config.headers['Authorization'] = 'Bearer ' + token;
            }
            return fetch(url, config);
        },

        applyAxiosToken() {
            if (typeof axios === 'undefined') {
                return;
            }
            const token = this.getToken();
            if (token) {
                axios.defaults.headers.common['Authorization'] = 'Bearer ' + token;
            } else {
                delete axios.defaults.headers.common['Authorization'];
            }
        },

        attachAxiosInterceptor() {
            if (typeof axios === 'undefined' || axios.__ADMIN_INTERCEPTOR__) {
                return;
            }
            axios.__ADMIN_INTERCEPTOR__ = true;
            axios.interceptors.response.use(
                (response) => response,
                (error) => {
                    if (error && error.response && error.response.status === 401) {
                        AdminAuth.handleUnauthorized();
                    }
                    return Promise.reject(error);
                }
            );
        }
    };

    function bootstrap() {
        AdminAuth.applyAxiosToken();
        AdminAuth.attachAxiosInterceptor();

        document.addEventListener('DOMContentLoaded', () => {
            const body = document.body;
            if (!body || !body.dataset) {
                return;
            }
            const requiresAuth = body.dataset.requireAuth === 'true';
            const isLoginPage = body.dataset.authPage === 'login';
            const isGateway = body.dataset.authPage === 'gateway';

            if (requiresAuth) {
                AdminAuth.requireAuth();
            } else if (isLoginPage) {
                if (AdminAuth.hasValidSession()) {
                    AdminAuth.verifyToken().then(() => {
                        AdminAuth.redirectToDashboard();
                    }).catch(() => AdminAuth.clear());
                }
            } else if (isGateway) {
                if (AdminAuth.hasValidSession()) {
                    AdminAuth.redirectToDashboard();
                } else {
                    AdminAuth.redirectToLogin();
                }
            }
        });
    }

    bootstrap();
    window.AdminAuth = AdminAuth;
})(window);
