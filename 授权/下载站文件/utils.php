<?php
/**
 * 工具类
 * 提供通用的工具方法
 */

class Utils {
    /**
     * 返回 JSON 响应
     * @param int $code HTTP 状态码
     * @param string $message 响应消息
     * @param mixed $data 附加数据
     */
    public static function jsonResponse($code, $message, $data = null) {
        http_response_code($code);
        header('Content-Type: application/json; charset=utf-8');

        $response = [
            'code' => $code,
            'message' => $message
        ];

        if ($data !== null) {
            $response['data'] = $data;
        }

        echo json_encode($response, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
        exit;
    }

    /**
     * 验证许可证格式
     * @param string $licenseCode 许可证代码
     * @return bool
     */
    public static function validateLicenseFormat($licenseCode) {
        // 许可证格式验证：应该是32位字符串，包含字母和数字
        if (empty($licenseCode)) {
            return false;
        }

        // 检查长度（通常是32位或更长）
        if (strlen($licenseCode) < 16) {
            return false;
        }

        // 检查是否只包含字母、数字和连字符
        if (!preg_match('/^[A-Za-z0-9\-]+$/', $licenseCode)) {
            return false;
        }

        return true;
    }

    /**
     * 记录日志
     * @param string $message 日志消息
     * @param string $level 日志级别 (INFO, WARNING, ERROR)
     */
    public static function log($message, $level = 'INFO') {
        $logFile = LOG_PATH . date('Y-m-d') . '.log';
        $timestamp = date('Y-m-d H:i:s');
        $logMessage = "[$timestamp] [$level] $message" . PHP_EOL;

        // 确保日志目录存在
        if (!file_exists(LOG_PATH)) {
            mkdir(LOG_PATH, 0755, true);
        }

        file_put_contents($logFile, $logMessage, FILE_APPEND);
    }

    /**
     * 获取客户端真实 IP 地址
     * @return string
     */
    public static function getClientIP() {
        $ip = '';

        if (isset($_SERVER['HTTP_X_FORWARDED_FOR']) && !empty($_SERVER['HTTP_X_FORWARDED_FOR'])) {
            $ips = explode(',', $_SERVER['HTTP_X_FORWARDED_FOR']);
            $ip = trim($ips[0]);
        } elseif (isset($_SERVER['HTTP_X_REAL_IP']) && !empty($_SERVER['HTTP_X_REAL_IP'])) {
            $ip = $_SERVER['HTTP_X_REAL_IP'];
        } elseif (isset($_SERVER['REMOTE_ADDR'])) {
            $ip = $_SERVER['REMOTE_ADDR'];
        }

        // 验证 IP 格式
        if (filter_var($ip, FILTER_VALIDATE_IP)) {
            return $ip;
        }

        return '0.0.0.0';
    }

    /**
     * 生成随机字符串
     * @param int $length 长度
     * @return string
     */
    public static function generateRandomString($length = 32) {
        $characters = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ';
        $charactersLength = strlen($characters);
        $randomString = '';

        for ($i = 0; $i < $length; $i++) {
            $randomString .= $characters[rand(0, $charactersLength - 1)];
        }

        return $randomString;
    }

    /**
     * 生成许可证代码
     * 格式: XXXX-XXXX-XXXX-XXXX-XXXX-XXXX (6组，每组4位)
     * @return string
     */
    public static function generateLicenseCode() {
        $segments = [];
        for ($i = 0; $i < 6; $i++) {
            $segments[] = self::generateRandomString(4);
        }
        return implode('-', $segments);
    }

    /**
     * 加密字符串
     * @param string $data 要加密的数据
     * @return string
     */
    public static function encrypt($data) {
        if (empty($data)) {
            return '';
        }

        $key = ENCRYPTION_KEY;
        $iv = openssl_random_pseudo_bytes(openssl_cipher_iv_length('aes-256-cbc'));
        $encrypted = openssl_encrypt($data, 'aes-256-cbc', $key, 0, $iv);

        return base64_encode($encrypted . '::' . $iv);
    }

    /**
     * 解密字符串
     * @param string $data 要解密的数据
     * @return string
     */
    public static function decrypt($data) {
        if (empty($data)) {
            return '';
        }

        $key = ENCRYPTION_KEY;
        list($encrypted, $iv) = explode('::', base64_decode($data), 2);

        return openssl_decrypt($encrypted, 'aes-256-cbc', $key, 0, $iv);
    }

    /**
     * 验证日期格式
     * @param string $date 日期字符串
     * @param string $format 日期格式
     * @return bool
     */
    public static function validateDate($date, $format = 'Y-m-d H:i:s') {
        $d = DateTime::createFromFormat($format, $date);
        return $d && $d->format($format) === $date;
    }

    /**
     * 清理输入数据
     * @param string $data 输入数据
     * @return string
     */
    public static function sanitize($data) {
        $data = trim($data);
        $data = stripslashes($data);
        $data = htmlspecialchars($data, ENT_QUOTES, 'UTF-8');
        return $data;
    }

    /**
     * 格式化文件大小
     * @param int $bytes 字节数
     * @return string
     */
    public static function formatBytes($bytes) {
        if ($bytes >= 1073741824) {
            $bytes = number_format($bytes / 1073741824, 2) . ' GB';
        } elseif ($bytes >= 1048576) {
            $bytes = number_format($bytes / 1048576, 2) . ' MB';
        } elseif ($bytes >= 1024) {
            $bytes = number_format($bytes / 1024, 2) . ' KB';
        } elseif ($bytes > 1) {
            $bytes = $bytes . ' bytes';
        } elseif ($bytes == 1) {
            $bytes = $bytes . ' byte';
        } else {
            $bytes = '0 bytes';
        }

        return $bytes;
    }
}
