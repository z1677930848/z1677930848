<?php
/**
 * CI/CD 专用上传接口
 * 端点: /api/ci/upload.php
 *
 * 用于 GitHub Actions 等 CI/CD 系统自动上传版本文件
 * 使用静态 API Key 认证，无需登录获取 token
 *
 * 认证方式:
 * - Header: X-CI-Key: your-ci-api-key
 *
 * 参数:
 * - file: 文件 (multipart/form-data)
 * - component: 组件代码 (admin, api, node)
 * - arch: 架构 (amd64, arm64)
 * - changelog: 更新日志 (可选)
 * - description: 版本描述 (可选)
 */

require_once dirname(__DIR__, 2) . '/config.php';
require_once dirname(__DIR__, 2) . '/database.php';
require_once dirname(__DIR__, 2) . '/utils.php';

// CI/CD API Key - 从环境变量读取，请在服务器上配置
define('CI_API_KEY', env_or_default('LINGCDN_CI_API_KEY', ''));

// CORS
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, X-CI-Key');

// 处理 OPTIONS 请求
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// 只接收 POST
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    Utils::jsonResponse(405, '只支持 POST 请求');
}

// 验证 CI API Key
function verifyCIKey(): void
{
    $key = $_SERVER['HTTP_X_CI_KEY'] ?? '';

    if (empty(CI_API_KEY)) {
        Utils::log("CI上传失败: 服务器未配置 LINGCDN_CI_API_KEY 环境变量", 'ERROR');
        Utils::jsonResponse(500, '服务器未配置 CI API Key');
    }

    if (empty($key)) {
        Utils::log("CI上传失败: 请求未提供 X-CI-Key", 'ERROR');
        Utils::jsonResponse(401, '缺少 X-CI-Key 认证头');
    }

    if (!hash_equals(CI_API_KEY, $key)) {
        Utils::log("CI上传失败: API Key 无效", 'ERROR');
        Utils::jsonResponse(401, 'API Key 无效');
    }
}

/**
 * 从文件名中解析版本号
 */
function parseVersionFromFilename(string $filename): ?string
{
    if (preg_match('/v?(\d+\.\d+\.\d+)/', $filename, $matches)) {
        return $matches[1];
    }
    return null;
}

/**
 * 获取组件名称
 */
function getComponentName(string $component): string
{
    $map = [
        'admin' => 'Ling Admin',
        'api'   => 'Ling API',
        'node'  => 'Ling Node',
    ];
    return $map[$component] ?? ucfirst($component);
}

try {
    // 验证 API Key
    verifyCIKey();

    // 检查文件
    if (!isset($_FILES['file']) || $_FILES['file']['error'] !== UPLOAD_ERR_OK) {
        $errorMsg = '文件上传失败';
        if (isset($_FILES['file']['error'])) {
            switch ($_FILES['file']['error']) {
                case UPLOAD_ERR_INI_SIZE:
                case UPLOAD_ERR_FORM_SIZE:
                    $errorMsg = '文件大小超过限制';
                    break;
                case UPLOAD_ERR_PARTIAL:
                    $errorMsg = '文件只上传了一部分';
                    break;
                case UPLOAD_ERR_NO_FILE:
                    $errorMsg = '没有选择文件';
                    break;
                default:
                    $errorMsg = '文件上传失败，错误代码：' . $_FILES['file']['error'];
            }
        }
        Utils::jsonResponse(400, $errorMsg);
    }

    // 参数
    $component = $_POST['component'] ?? '';
    $arch = $_POST['arch'] ?? 'amd64';
    $changelog = $_POST['changelog'] ?? '';
    $description = $_POST['description'] ?? '';
    $os = 'linux'; // 固定为 linux

    // 校验组件
    $validComponents = ['admin', 'api', 'node'];
    if (!in_array($component, $validComponents, true)) {
        Utils::jsonResponse(400, '无效的组件类型，支持: ' . implode(', ', $validComponents));
    }

    // 校验架构
    $validArch = ['amd64', 'arm64'];
    if (!in_array($arch, $validArch, true)) {
        Utils::jsonResponse(400, '无效的架构，支持: ' . implode(', ', $validArch));
    }

    $uploadedFile = $_FILES['file'];
    $tmpPath = $uploadedFile['tmp_name'];
    $originalFilename = basename($uploadedFile['name']);
    $fileSize = $uploadedFile['size'];

    // 扩展名校验
    $allowedExtensions = ['zip', 'tar.gz', 'tgz'];
    $fileExt = strtolower(pathinfo($originalFilename, PATHINFO_EXTENSION));
    if ($fileExt === 'gz' && substr($originalFilename, -7) === '.tar.gz') {
        $fileExt = 'tar.gz';
    }
    if (!in_array($fileExt, $allowedExtensions, true)) {
        Utils::jsonResponse(400, '不支持的文件格式，只支持: ' . implode(', ', $allowedExtensions));
    }

    // 文件名规范校验
    $pattern = '/^(ling-)?' . preg_quote($component, '/') . '-v\d+\.\d+\.\d+-linux-(amd64|arm64)\.(zip|tar\.gz|tgz)$/i';
    if (!preg_match($pattern, $originalFilename)) {
        Utils::jsonResponse(400, '文件名不符合规范，示例：ling-admin-v1.2.3-linux-amd64.zip');
    }

    // 解析版本号
    $version = parseVersionFromFilename($originalFilename);
    if (!$version) {
        Utils::jsonResponse(400, '无法从文件名解析版本号');
    }

    // 计算哈希
    $md5 = md5_file($tmpPath);
    $sha256 = hash_file('sha256', $tmpPath);

    // 目标目录
    $targetDir = UPDATE_FILES_PATH . "{$component}/{$os}/{$arch}/";
    if (!is_dir($targetDir) && !mkdir($targetDir, 0755, true)) {
        Utils::jsonResponse(500, '创建目标目录失败');
    }
    $targetPath = $targetDir . $originalFilename;

    // 检查是否已存在
    $needMove = true;
    if (file_exists($targetPath)) {
        $existingMd5 = md5_file($targetPath);
        if ($existingMd5 === $md5) {
            // 同名同内容，跳过移动
            $needMove = false;
            Utils::log("CI上传: 文件已存在且内容相同，跳过: {$originalFilename}", 'INFO');
        }
    }

    if ($needMove && !move_uploaded_file($tmpPath, $targetPath)) {
        Utils::jsonResponse(500, '移动文件失败');
    }

    @chmod($targetPath, 0644);

    // 同步数据库
    $db = Database::getInstance()->getConnection();
    $releaseTime = date('Y-m-d H:i:s');
    $componentName = getComponentName($component);

    // 检查是否已存在该版本
    $checkStmt = $db->prepare("
        SELECT id FROM versions
        WHERE component_code = :component AND version = :version AND os = :os AND arch = :arch
        LIMIT 1
    ");
    $checkStmt->execute([
        'component' => $component,
        'version' => $version,
        'os' => $os,
        'arch' => $arch
    ]);
    $existing = $checkStmt->fetch();

    if ($existing) {
        // 更新现有记录
        $updateStmt = $db->prepare("
            UPDATE versions SET
                filename = :filename,
                file_size = :file_size,
                file_md5 = :file_md5,
                file_sha256 = :file_sha256,
                description = CASE WHEN :description = '' THEN description ELSE :description END,
                changelog = CASE WHEN :changelog = '' THEN changelog ELSE :changelog END,
                release_time = :release_time,
                status = 1,
                updated_at = NOW()
            WHERE id = :id
        ");
        $updateStmt->execute([
            'filename' => $originalFilename,
            'file_size' => $fileSize,
            'file_md5' => $md5,
            'file_sha256' => $sha256,
            'description' => $description,
            'changelog' => $changelog,
            'release_time' => $releaseTime,
            'id' => $existing['id']
        ]);
        $versionId = $existing['id'];
    } else {
        // 插入新记录
        $insertStmt = $db->prepare("
            INSERT INTO versions (
                component_code, component_name, version, os, arch,
                filename, file_size, file_md5, file_sha256,
                description, changelog, release_time,
                is_required_update, status, download_count
            ) VALUES (
                :component_code, :component_name, :version, :os, :arch,
                :filename, :file_size, :file_md5, :file_sha256,
                :description, :changelog, :release_time,
                0, 1, 0
            )
        ");
        $insertStmt->execute([
            'component_code' => $component,
            'component_name' => $componentName,
            'version' => $version,
            'os' => $os,
            'arch' => $arch,
            'filename' => $originalFilename,
            'file_size' => $fileSize,
            'file_md5' => $md5,
            'file_sha256' => $sha256,
            'description' => $description ?: "CI/CD 自动发布 v{$version}",
            'changelog' => $changelog ?: "版本 {$version} 自动发布",
            'release_time' => $releaseTime
        ]);
        $versionId = $db->lastInsertId();
    }

    // 禁用同组件同架构的其他版本，确保当前为最新
    $disableStmt = $db->prepare("
        UPDATE versions
        SET status = 0
        WHERE component_code = :component AND os = :os AND arch = :arch AND id <> :id
    ");
    $disableStmt->execute([
        'component' => $component,
        'os' => $os,
        'arch' => $arch,
        'id' => $versionId
    ]);

    Utils::log(
        "CI上传成功: {$originalFilename}, 版本: {$version}, 平台: {$os}/{$arch}, SHA256: {$sha256}",
        'INFO'
    );

    Utils::jsonResponse(200, '上传成功', [
        'filename' => $originalFilename,
        'path' => "/updates/{$component}/{$os}/{$arch}/{$originalFilename}",
        'version' => $version,
        'component' => $component,
        'os' => $os,
        'arch' => $arch,
        'size' => $fileSize,
        'sizeFormatted' => Utils::formatBytes($fileSize),
        'md5' => $md5,
        'sha256' => $sha256,
        'versionId' => $versionId
    ]);

} catch (Exception $e) {
    Utils::log("CI上传异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误: ' . $e->getMessage());
}
