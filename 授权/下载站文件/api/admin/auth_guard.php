<?php
/**
 * 管理员认证守卫
 */

$rootPath = dirname(__DIR__, 2);
require_once $rootPath . '/config.php';
require_once $rootPath . '/utils.php';

class AdminAuth
{
    /**
     * 登录并返回 token
     */
    public static function login(string $username, string $password): array
    {
        $username = trim($username);
        $hash = defined('ADMIN_PASSWORD_HASH') ? ADMIN_PASSWORD_HASH : '';

        if ($username !== ADMIN_USERNAME) {
            Utils::jsonResponse(401, '账号或密码错误');
        }

        if (!empty($hash)) {
            if (!password_verify($password, $hash)) {
                Utils::jsonResponse(401, '账号或密码错误');
            }
        } elseif ($password !== ADMIN_PASSWORD) {
            Utils::jsonResponse(401, '账号或密码错误');
        }

        return self::issueToken($username);
    }

    /**
     * 校验请求头中的 token
     */
    public static function requireAuth(): array
    {
        $token = self::extractToken();
        $payload = self::validateToken($token);
        if (!$payload) {
            Utils::jsonResponse(401, '未授权或登录已过期');
        }
        return $payload;
    }

    /**
     * 生成 token
     */
    private static function issueToken(string $username): array
    {
        $ttl = defined('ADMIN_TOKEN_TTL') ? (int)ADMIN_TOKEN_TTL : 7200;
        if ($ttl <= 0) {
            $ttl = 7200;
        }

        $now = time();
        $payload = [
            'sub' => $username,
            'iat' => $now,
            'exp' => $now + $ttl,
            'jti' => bin2hex(random_bytes(8))
        ];

        $body = rtrim(strtr(base64_encode(json_encode($payload, JSON_UNESCAPED_UNICODE)), '+/', '-_'), '=');
        $signature = hash_hmac('sha256', $body, API_SECRET_KEY);

        return [
            'token' => $body . '.' . $signature,
            'expire_at' => date('c', $payload['exp']),
            'username' => $username
        ];
    }

    /**
     * 验证 token
     */
    private static function validateToken(?string $token): ?array
    {
        if (empty($token)) {
            return null;
        }

        $parts = explode('.', $token);
        if (count($parts) !== 2) {
            return null;
        }

        [$body, $signature] = $parts;
        $expected = hash_hmac('sha256', $body, API_SECRET_KEY);
        if (!hash_equals($expected, $signature)) {
            return null;
        }

        $payload = json_decode(base64_decode(strtr($body, '-_', '+/')), true);
        if (!$payload || !isset($payload['exp']) || $payload['exp'] < time()) {
            return null;
        }

        return $payload;
    }

    /**
     * 从请求中获取 token
     */
    private static function extractToken(): ?string
    {
        $auth = $_SERVER['HTTP_AUTHORIZATION'] ?? '';
        if (stripos($auth, 'bearer ') === 0) {
            return trim(substr($auth, 7));
        }

        if (!empty($_SERVER['HTTP_X_ADMIN_TOKEN'])) {
            return trim($_SERVER['HTTP_X_ADMIN_TOKEN']);
        }

        if (!empty($_GET['admin_token'])) {
            return trim($_GET['admin_token']);
        }

        if (!empty($_POST['admin_token'])) {
            return trim($_POST['admin_token']);
        }

        return null;
    }
}
