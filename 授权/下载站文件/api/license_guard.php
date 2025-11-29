<?php
/**
 * LicenseGuard
 * 将请求中的 license_code / system_token / domain 做校验，确保授权有效再继续后续操作
 */

require_once __DIR__ . '/../utils.php';

class LicenseGuard
{
    /**
     * 从请求头或参数中提取许可证信息
     * @return array
     */
    public static function getCredentialsFromRequest(): array
    {
        $licenseCode = self::valueFromRequest(
            ['HTTP_X_LICENSE_CODE', 'HTTP_AUTHORIZATION_LICENSE'],
            ['license_code', 'license', 'lic']
        );
        $systemToken = self::valueFromRequest(
            ['HTTP_X_SYSTEM_TOKEN'],
            ['system_token', 'token']
        );
        $domain = self::valueFromRequest(
            ['HTTP_X_LICENSE_DOMAIN'],
            ['domain', 'host']
        );

        if (empty($domain)) {
            if (!empty($_SERVER['HTTP_HOST'])) {
                $domain = $_SERVER['HTTP_HOST'];
            } elseif (!empty($_SERVER['SERVER_NAME'])) {
                $domain = $_SERVER['SERVER_NAME'];
            }
        }

        return [
            'license_code' => trim($licenseCode),
            'system_token' => trim($systemToken),
            'domain'       => trim($domain)
        ];
    }

    /**
     * 校验授权信息
     * @param PDO $db
     * @param string $licenseCode
     * @param string $systemToken
     * @param string $domain
     * @param string $context
     * @return array
     */
    public static function verify(PDO $db, string $licenseCode, string $systemToken, string $domain = '', string $context = 'access'): array
    {
        $licenseCode = trim($licenseCode);
        $systemToken = trim($systemToken);
        $domain = trim(strtolower($domain));

        if (empty($licenseCode)) {
            Utils::jsonResponse(401, 'license_code required');
        }

        if (!Utils::validateLicenseFormat($licenseCode)) {
            Utils::jsonResponse(400, 'invalid license_code');
        }

        if (empty($systemToken)) {
            Utils::jsonResponse(401, 'system_token required');
        }

        $stmt = $db->prepare("SELECT * FROM licenses WHERE license_code = :code LIMIT 1");
        $stmt->execute(['code' => $licenseCode]);
        $license = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$license) {
            Utils::log("授权失败: 许可证不存在 - $licenseCode", 'WARNING');
            Utils::jsonResponse(403, 'license not found');
        }

        if ((int)$license['status'] !== 1) {
            self::deny($db, $license, $systemToken, $domain, 'license disabled');
        }

        if (!empty($license['expire_time']) && strtotime($license['expire_time']) < time()) {
            self::deny($db, $license, $systemToken, $domain, 'license expired');
        }

        // 域名校验
        if (!empty($license['allowed_domains'])) {
            $allowedDomains = json_decode($license['allowed_domains'], true);
            if (!is_array($allowedDomains)) {
                $allowedDomains = [];
            }
            $allowedDomains = array_filter(array_map('strtolower', array_map('trim', $allowedDomains)));

            if (!empty($allowedDomains) && !in_array('*', $allowedDomains, true)) {
                if (empty($domain)) {
                    self::deny($db, $license, $systemToken, $domain, 'domain required');
                }

                $matched = false;
                foreach ($allowedDomains as $pattern) {
                    if (self::matchDomain($pattern, $domain)) {
                        $matched = true;
                        break;
                    }
                }

                if (!$matched) {
                    self::deny($db, $license, $systemToken, $domain, 'domain not authorized');
                }
            }
        }

        // 绑定系统 Token
        $boundTokens = [];
        if (!empty($license['bound_system_tokens'])) {
            $decodedTokens = json_decode($license['bound_system_tokens'], true);
            if (is_array($decodedTokens)) {
                $boundTokens = array_values(array_filter(array_map('strval', $decodedTokens)));
            }
        }

        if (!in_array($systemToken, $boundTokens, true)) {
            $maxDevices = isset($license['max_devices']) ? (int)$license['max_devices'] : 1;
            if ($maxDevices > 0 && count($boundTokens) >= $maxDevices) {
                self::deny($db, $license, $systemToken, $domain, 'license device limit reached');
            }

            $boundTokens[] = $systemToken;
            $updateStmt = $db->prepare("
                UPDATE licenses
                SET bound_system_tokens = :tokens
                WHERE id = :id
            ");
            $updateStmt->execute([
                'tokens' => json_encode(array_values(array_unique($boundTokens)), JSON_UNESCAPED_UNICODE),
                'id'     => $license['id']
            ]);

            $license['bound_system_tokens'] = json_encode($boundTokens);
        }

        // 更新 last_check_time & check_count
        $db->prepare("UPDATE licenses SET last_check_time = NOW(), check_count = check_count + 1 WHERE id = :id")
           ->execute(['id' => $license['id']]);

        self::writeLog($db, $license['id'], $systemToken, $domain, 'success', $context);

        return $license;
    }

    /**
     * 根据 Header / Query / Body 获取值
     */
    private static function valueFromRequest(array $serverKeys, array $paramKeys): string
    {
        foreach ($serverKeys as $key) {
            if (!empty($_SERVER[$key])) {
                return trim($_SERVER[$key]);
            }
        }

        foreach ($paramKeys as $key) {
            if (isset($_GET[$key]) && $_GET[$key] !== '') {
                return trim($_GET[$key]);
            }
            if (isset($_POST[$key]) && $_POST[$key] !== '') {
                return trim($_POST[$key]);
            }
        }

        return '';
    }

    /**
     * 领域匹配，支持 * 通配
     */
    private static function matchDomain(string $pattern, string $domain): bool
    {
        if ($pattern === '*') {
            return true;
        }

        if ($pattern === $domain) {
            return true;
        }

        if (strpos($pattern, '*') !== false) {
            $regex = '/^' . str_replace('\*', '.*', preg_quote($pattern, '/')) . '$/';
            return (bool)preg_match($regex, $domain);
        }

        return false;
    }

    /**
     * 拒绝访问并记录日志
     */
    private static function deny(?PDO $db, array $license, string $systemToken, string $domain, string $message): void
    {
        $code = $license['license_code'] ?? 'unknown';
        Utils::log("授权失败: {$message} - License: {$code} - Token: {$systemToken} - Domain: {$domain}", 'WARNING');

        if ($db && !empty($license['id'])) {
            self::writeLog($db, (int)$license['id'], $systemToken, $domain, 'failed', $message);
        }

        Utils::jsonResponse(403, $message);
    }

    /**
     * 写入授权日志
     */
    private static function writeLog(PDO $db, int $licenseId, string $systemToken, string $domain, string $status, string $message = ''): void
    {
        try {
            $stmt = $db->prepare("
                INSERT INTO authorization_logs
                (license_id, system_token, domain, ip_address, status, error_message, user_agent, created_at)
                VALUES (:license_id, :system_token, :domain, :ip, :status, :error, :ua, NOW())
            ");
            $stmt->execute([
                'license_id'  => $licenseId,
                'system_token'=> $systemToken,
                'domain'      => $domain ?: '-',
                'ip'          => Utils::getClientIP(),
                'status'      => $status,
                'error'       => $message,
                'ua'          => $_SERVER['HTTP_USER_AGENT'] ?? ''
            ]);
        } catch (Exception $e) {
            Utils::log("授权日志写入失败: " . $e->getMessage(), 'ERROR');
        }
    }
}
