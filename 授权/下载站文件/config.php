<?php
/**
 * 授权与更新系统 - 配置文件
 */

// 环境变量读取工具（没有则使用默认值）
function env_or_default(string $key, $default) {
    $value = getenv($key);
    return ($value === false || $value === '') ? $default : $value;
}

// 数据库配置（支持环境变量覆盖，便于安全部署）
define('DB_HOST', env_or_default('LINGCDN_DB_HOST', 'localhost'));
define('DB_NAME', env_or_default('LINGCDN_DB_NAME', 'lingcdn'));
define('DB_USER', env_or_default('LINGCDN_DB_USER', 'lingcdn'));
define('DB_PASS', env_or_default('LINGCDN_DB_PASS', '123456'));
define('DB_CHARSET', env_or_default('LINGCDN_DB_CHARSET', 'utf8mb4'));
define('DB_TIMEZONE', env_or_default('LINGCDN_DB_TIMEZONE', '+8:00'));

// 系统配置（务必在生产环境用环境变量覆盖默认值）
define('API_SECRET_KEY', env_or_default('LINGCDN_API_SECRET_KEY', 'your-secret-key-change-this')); // 请修改此密钥
define('ENCRYPTION_KEY', env_or_default('LINGCDN_ENCRYPTION_KEY', 'your-encryption-key-32-chars')); // 32位加密密钥

// 管理员校验登录数据
define('ADMIN_USERNAME', env_or_default('LINGCDN_ADMIN_USERNAME', 'admin')); // 默认管理员名称
define('ADMIN_PASSWORD', env_or_default('LINGCDN_ADMIN_PASSWORD', 'admin123')); // 默认密码，请在环境处必须修改
define('ADMIN_PASSWORD_HASH', env_or_default('LINGCDN_ADMIN_PASSWORD_HASH', '')); // 如果使用hash数据，确保填写在此
define('ADMIN_TOKEN_TTL', intval(env_or_default('LINGCDN_ADMIN_TOKEN_TTL', 7200))); // 登录token有效期为2小时

// 更新文件存储路径
define('UPDATE_FILES_PATH', __DIR__ . '/updates/');

// 日志路径
define('LOG_PATH', __DIR__ . '/logs/');

// 时区设置
date_default_timezone_set('Asia/Shanghai');

// 错误报告
error_reporting(E_ALL);
ini_set('display_errors', 0);
ini_set('log_errors', 1);

// 确保日志目录、更新目录可用，不打断请求，仅记录可读性更好的错误
if (!is_dir(LOG_PATH)) {
    @mkdir(LOG_PATH, 0755, true);
}
if (is_dir(LOG_PATH) && is_writable(LOG_PATH)) {
    ini_set('error_log', LOG_PATH . 'error.log');
} else {
    // 回退到系统临时目录，避免因权限导致日志写入失败
    ini_set('error_log', sys_get_temp_dir() . DIRECTORY_SEPARATOR . 'lingcdn-download.log');
}

if (!is_dir(UPDATE_FILES_PATH)) {
    @mkdir(UPDATE_FILES_PATH, 0755, true);
}
