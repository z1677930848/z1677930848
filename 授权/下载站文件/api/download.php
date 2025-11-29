<?php
/**
 * 文件下载 API
 * 提供版本文件下载服务
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';
require_once __DIR__ . '/license_guard.php';

// 获取请求的文件路径
$requestPath = trim($_GET['file'] ?? '');

if (empty($requestPath)) {
    Utils::jsonResponse(400, 'missing file parameter');
}

// 安全检查：防止目录遍历攻击
$requestPath = str_replace(['../', '..\\'], '', $requestPath);

// 构建实际文件路径
$filePath = UPDATE_FILES_PATH . $requestPath;

// 检查文件是否存在
if (!file_exists($filePath) || !is_file($filePath)) {
    Utils::jsonResponse(404, 'file not found');
}

try {
    $db = Database::getInstance()->getConnection();

    // 授权校验
    $credentials = LicenseGuard::getCredentialsFromRequest();
    $licenseInfo = LicenseGuard::verify(
        $db,
        $credentials['license_code'],
        $credentials['system_token'],
        $credentials['domain'],
        'download'
    );

    // 从路径中提取信息
    $pathParts = explode('/', $requestPath);
    $filename = end($pathParts);

    // 查询数据库验证文件
    $stmt = $db->prepare("
        SELECT * FROM versions
        WHERE filename = :filename
        AND status = 1
        LIMIT 1
    ");
    $stmt->execute(['filename' => $filename]);
    $version = $stmt->fetch();

    if (!$version) {
        http_response_code(403);
        die('Access denied');
    }

    // 记录下载日志
    $stmt = $db->prepare("
        INSERT INTO download_logs
        (version_id, ip_address, user_agent, download_time)
        VALUES (:version_id, :ip, :user_agent, NOW())
    ");
    $stmt->execute([
        'version_id' => $version['id'],
        'ip' => Utils::getClientIP(),
        'user_agent' => $_SERVER['HTTP_USER_AGENT'] ?? ''
    ]);

    // 更新下载次数
    $stmt = $db->prepare("UPDATE versions SET download_count = download_count + 1 WHERE id = :id");
    $stmt->execute(['id' => $version['id']]);

    // 设置下载响应头
    header('Content-Type: application/octet-stream');
    header('Content-Disposition: attachment; filename="' . basename($filePath) . '"');
    header('Content-Length: ' . filesize($filePath));
    header('Cache-Control: must-revalidate');
    header('Pragma: public');

    // 输出文件内容
    readfile($filePath);

    Utils::log("文件下载: $filename - License: {$licenseInfo['license_code']} - IP: " . Utils::getClientIP(), 'INFO');

} catch (Exception $e) {
    Utils::log("文件下载异常: " . $e->getMessage(), 'ERROR');
    http_response_code(500);
    die('Server error');
}
