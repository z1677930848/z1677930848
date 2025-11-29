<?php
/**
 * 文件上传 API
 * 端点: /api/admin/upload.php
 *
 * 用于上传版本文件到 updates 目录，并同步 versions 表为最新版。
 *
 * 参数:
 * - file: 文件 (multipart/form-data)
 * - component: 组件代码 (admin, api, node)
 * - os: 操作系统 (linux, windows, darwin)
 * - arch: 架构 (amd64, arm64, 386)
 * - overwrite: 是否覆盖同名不同内容文件（可选，值为 "true" 时覆盖）
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';
require_once __DIR__ . '/auth_guard.php';

/**
 * 从文件名中解析语义化版本号，形如 v1.2.3 或 1.2.3
 */
function parseVersionFromFilename(string $filename): ?string
{
    if (preg_match('/v?(\\d+\\.\\d+\\.\\d+)/', $filename, $matches)) {
        return $matches[1];
    }
    return null;
}

/**
 * 获取组件名称，便于写入 versions 表
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

/**
 * 生成 changelog HTML，便于前端直接查看
 */
function generateChangelogHtml(string $component, string $version, string $description, string $changelog): void
{
    $changelogDir = __DIR__ . '/../../changelog';
    if (!is_dir($changelogDir)) {
        @mkdir($changelogDir, 0755, true);
    }
    $content = $changelog ?: $description;
    $html = <<<HTML
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{$component} v{$version} 更新说明</title>
    <style>
        body { font-family: -apple-system,BlinkMacSystemFont,'Segoe UI',Arial,sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; line-height: 1.6; }
        h1 { color: #1f2937; }
        .meta { color: #6b7280; margin-bottom: 16px; }
        pre { white-space: pre-wrap; background: #f9fafb; padding: 12px; border-radius: 8px; }
    </style>
    </head>
<body>
    <h1>{$component} v{$version}</h1>
    <div class="meta">自动生成的更新说明</div>
    <pre>{$content}</pre>
</body>
</html>
HTML;
    file_put_contents($changelogDir . '/v' . $version . '.html', $html);
}

// CORS
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization, X-Admin-Token');

// 处理 OPTIONS 请求
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// 只接收 POST
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    Utils::jsonResponse(405, '只支持 POST 请求');
}

AdminAuth::requireAuth();

try {
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
    $os = $_POST['os'] ?? 'linux';
    $arch = $_POST['arch'] ?? 'amd64';

    // 校验
    $validComponents = ['admin', 'api', 'node'];
    if (!in_array($component, $validComponents, true)) {
        Utils::jsonResponse(400, '无效的组件类型');
    }

    // 仅支持 linux 平台
    $validOs = ['linux'];
    if (!in_array($os, $validOs, true)) {
        Utils::jsonResponse(400, '当前仅支持 Linux 平台');
    }

    $validArch = ['amd64', 'arm64', '386'];
    if (!in_array($arch, $validArch, true)) {
        Utils::jsonResponse(400, '无效的架构');
    }

    $uploadedFile = $_FILES['file'];
    $tmpPath = $uploadedFile['tmp_name'];
    $originalFilename = basename($uploadedFile['name']);
    $fileSize = $uploadedFile['size'];

    // 扩展名校验
    $allowedExtensions = ['zip', 'tar.gz', 'tgz', 'tar', 'gz'];
    $fileExt = strtolower(pathinfo($originalFilename, PATHINFO_EXTENSION));
    if ($fileExt === 'gz' && substr($originalFilename, -7) === '.tar.gz') {
        $fileExt = 'tar.gz';
    }
    if (!in_array($fileExt, $allowedExtensions, true)) {
        Utils::jsonResponse(400, '不支持的文件格式，只支持: ' . implode(', ', $allowedExtensions));
    }

    // 哈希
    $md5 = md5_file($tmpPath);
    $sha256 = hash_file('sha256', $tmpPath);

    // 目录与路径
    $targetDir = UPDATE_FILES_PATH . "{$component}/{$os}/{$arch}/";
    if (!is_dir($targetDir) && !mkdir($targetDir, 0755, true)) {
        Utils::jsonResponse(500, '创建目标目录失败');
    }
    $targetPath = $targetDir . $originalFilename;

    // 已存在处理
    $needMove = true;
    if (file_exists($targetPath)) {
        $existingMd5 = md5_file($targetPath);
        if ($existingMd5 === $md5) {
            // 同名同内容，跳过移动但继续同步数据库
            $needMove = false;
        } else {
            // 同名不同内容，需要明确覆盖
            if (!isset($_POST['overwrite']) || $_POST['overwrite'] !== 'true') {
                Utils::jsonResponse(409, '文件已存在但内容不同，请确认是否覆盖', [
                    'filename' => $originalFilename,
                    'existingMd5' => $existingMd5,
                    'newMd5' => $md5
                ]);
            }
        }
    }

    if ($needMove && !move_uploaded_file($tmpPath, $targetPath)) {
        Utils::jsonResponse(500, '移动文件失败');
    }

    // 权限
    @chmod($targetPath, 0644);

    // 文件名规范校验（强制 linux 平台）
    $pattern = '/^(ling-)?' . preg_quote($component, '/') . '-v\\d+\\.\\d+\\.\\d+-linux-(amd64|arm64|386)\\.(zip|tar\\.gz|tgz|tar|gz)$/i';
    if (!preg_match($pattern, $originalFilename)) {
        Utils::jsonResponse(400, '文件名不符合规范，示例：ling-admin-v1.2.3-linux-amd64.zip');
    }

    // 版本号
    $version = parseVersionFromFilename($originalFilename);
    if (!$version) {
        Utils::jsonResponse(400, '无法从文件名解析版本号，示例：ling-admin-v1.2.3-linux-amd64.zip');
    }

    // 同步数据库：同组件+OS+架构只保留最新版启用
    $db = Database::getInstance()->getConnection();
    $releaseTime = date('Y-m-d H:i:s');

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
        $updateStmt = $db->prepare("
            UPDATE versions SET
                filename = :filename,
                file_size = :file_size,
                file_md5 = :file_md5,
                file_sha256 = :file_sha256,
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
            'release_time' => $releaseTime,
            'id' => $existing['id']
        ]);
        $versionId = $existing['id'];
    } else {
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
            'component_name' => getComponentName($component),
            'version' => $version,
            'os' => $os,
            'arch' => $arch,
            'filename' => $originalFilename,
            'file_size' => $fileSize,
            'file_md5' => $md5,
            'file_sha256' => $sha256,
            'description' => '通过上传接口自动入库',
            'changelog' => '',
            'release_time' => $releaseTime
        ]);
        $versionId = $db->lastInsertId();
    }

    // 禁用其他版本，确保当前为最新版
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

    // 自动补充描述/更新说明，便于生成 changelog
    $desc = "自动生成：{$originalFilename}，大小 " . Utils::formatBytes($fileSize) . "，MD5 {$md5}";
    $stmtSelect = $db->prepare("SELECT description, changelog FROM versions WHERE id = ?");
    $stmtSelect->execute([$versionId]);
    $meta = $stmtSelect->fetch();
    $description = $meta['description'] ?: $desc;
    $changelog = $meta['changelog'] ?: $desc;
    $stmtUpdateDesc = $db->prepare("UPDATE versions SET description = :description, changelog = :changelog WHERE id = :id");
    $stmtUpdateDesc->execute([
        'description' => $description,
        'changelog' => $changelog,
        'id' => $versionId
    ]);
    generateChangelogHtml($componentName, $version, $description, $changelog);

    Utils::log(
        "文件上传成功并入库: {$originalFilename}, 版本: {$version}, 平台: {$os}/{$arch}, MD5: {$md5}, 大小: " . Utils::formatBytes($fileSize),
        'INFO'
    );

    Utils::jsonResponse(200, '上传成功', [
        'filename' => $originalFilename,
        'path' => "/updates/{$component}/{$os}/{$arch}/{$originalFilename}",
        'fullPath' => $targetPath,
        'size' => $fileSize,
        'sizeFormatted' => Utils::formatBytes($fileSize),
        'md5' => $md5,
        'sha256' => $sha256,
        'component' => $component,
        'os' => $os,
        'arch' => $arch,
        'version' => $version,
        'versionId' => $versionId,
        'existed' => isset($existingMd5) && $existingMd5 === $md5
    ]);

} catch (Exception $e) {
    Utils::log("文件上传异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误: ' . $e->getMessage());
}
