package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/iwind/TeaGo/actions"
	"net"
	"net/http"
)

type UserShouldAuth struct {
	action *actions.ActionObject
}

func (this *UserShouldAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	if teaconst.IsRecoverMode {
		actionPtr.Object().RedirectURL("/recover")
		return false
	}

	this.action = actionPtr.Object()

	// 安全相关
	var action = this.action
	securityConfig, _ := configloaders.LoadSecurityConfig()
	if securityConfig == nil {
		action.AddHeader("X-Frame-Options", "SAMEORIGIN")
	} else if len(securityConfig.Frame) > 0 {
		action.AddHeader("X-Frame-Options", securityConfig.Frame)
	}
	action.AddHeader("Content-Security-Policy", "default-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")

	// 检查IP
	if !checkIP(securityConfig, action.RequestRemoteIP()) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}
	remoteAddr, _, _ := net.SplitHostPort(action.Request.RemoteAddr)
	if len(remoteAddr) > 0 && remoteAddr != action.RequestRemoteIP() && !checkIP(securityConfig, remoteAddr) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}

	// 检查请求
	if !checkRequestSecurity(securityConfig, action.Request) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

// StoreAdmin 存储用户名到SESSION
func (this *UserShouldAuth) StoreAdmin(adminId int64, remember bool) {
	// admin登录不标记user态
	this.storeSession(adminId, "", false, remember)
}

// StoreUser 存储用户端登录信息到SESSION
func (this *UserShouldAuth) StoreUser(userId int64, username string, remember bool) {
	this.storeSession(userId, username, true, remember)
}

func (this *UserShouldAuth) storeSession(id int64, username string, isUser bool, remember bool) {
	// 修改sid的时间
	if remember {
		cookie := &http.Cookie{
			Name:     teaconst.CookieSID,
			Value:    this.action.Session().Sid,
			Path:     "/",
			MaxAge:   14 * 86400,
			HttpOnly: true,
		}
		if this.action.Request.TLS != nil {
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Secure = true
		}
		this.action.AddCookie(cookie)
	} else {
		cookie := &http.Cookie{
			Name:     teaconst.CookieSID,
			Value:    this.action.Session().Sid,
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		}
		if this.action.Request.TLS != nil {
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Secure = true
		}
		this.action.AddCookie(cookie)
	}
	this.action.Session().Write("adminId", numberutils.FormatInt64(id))
	this.action.Session().Write("userId", numberutils.FormatInt64(id))
	if len(username) > 0 {
		this.action.Session().Write("username", username)
	}
	if isUser {
		this.action.Session().WriteInt("isUser", 1)
	} else {
		this.action.Session().WriteInt("isUser", 0)
	}
}

func (this *UserShouldAuth) IsUser() bool {
	// 管理端登录判断，用户端会话不复用
	if this.IsUserPortal() {
		return false
	}
	return this.action.Session().GetInt("adminId") > 0
}

// IsUserPortal 是否为用户端会话
func (this *UserShouldAuth) IsUserPortal() bool {
	return this.action.Session().GetInt("isUser") > 0 && this.UserId() > 0
}

func (this *UserShouldAuth) AdminId() int {
	return this.action.Session().GetInt("adminId")
}

func (this *UserShouldAuth) UserId() int64 {
	userId := this.action.Session().GetInt64("userId")
	if userId > 0 {
		return userId
	}
	return int64(this.action.Session().GetInt("adminId"))
}

func (this *UserShouldAuth) Username() string {
	return this.action.Session().GetString("username")
}

func (this *UserShouldAuth) Logout() {
	this.action.Session().Delete()
}
