<?php
/**
 * 版本管理 API（后台）
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';
require_once __DIR__ . '/auth_guard.php';

header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization, X-Admin-Token');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// 平台限制：仅支持 Linux
function ensureLinux(string $os) {
    if ($os !== 'linux') {
        Utils::jsonResponse(400, '当前仅支持 Linux 平台');
    }
}

// 重新计算文件哈希并更新数据库
function verifyFileHash(array $version) {
    $filePath = UPDATE_FILES_PATH .
        $version['component_code'] . '/' .
        $version['os'] . '/' .
        $version['arch'] . '/' .
        $version['filename'];

    if (!file_exists($filePath)) {
        Utils::jsonResponse(404, '文件不存在: ' . $filePath);
    }

    $size = filesize($filePath);
    $md5 = md5_file($filePath);
    $sha256 = hash_file('sha256', $filePath);

    return [$size, $md5, $sha256, $filePath];
}

try {
    AdminAuth::requireAuth();
    $db = Database::getInstance()->getConnection();

    if ($_SERVER['REQUEST_METHOD'] === 'GET') {
        $action = $_GET['action'] ?? '';
        if ($action === 'list') {
            $showAll = isset($_GET['show_all']) && $_GET['show_all'] == '1';
            $sql = $showAll
                ? "SELECT * FROM versions WHERE os = 'linux' ORDER BY release_time DESC, id DESC"
                : "SELECT * FROM versions WHERE os = 'linux' AND status = 1 ORDER BY release_time DESC, id DESC";
            $stmt = $db->prepare($sql);
            $stmt->execute();
            $versions = $stmt->fetchAll();
            Utils::jsonResponse(200, 'success', $versions);
        } else {
            Utils::jsonResponse(400, '无效的操作');
        }
    } elseif ($_SERVER['REQUEST_METHOD'] === 'POST') {
        $input = file_get_contents('php://input');
        $data = json_decode($input, true);
        if (!$data) {
            $data = $_POST;
        }
        $action = $data['action'] ?? '';

        if ($action === 'add') {
            $required = ['component_code', 'component_name', 'version', 'os', 'arch', 'filename', 'file_size', 'file_md5', 'release_time'];
            foreach ($required as $field) {
                if (empty($data[$field])) {
                    Utils::jsonResponse(400, "缺少必填字段: {$field}");
                }
            }
            ensureLinux($data['os']);

            $checkSql = "SELECT id FROM versions WHERE component_code = :component_code AND version = :version AND os = :os AND arch = :arch";
            $checkStmt = $db->prepare($checkSql);
            $checkStmt->execute([
                'component_code' => $data['component_code'],
                'version' => $data['version'],
                'os' => $data['os'],
                'arch' => $data['arch']
            ]);
            if ($checkStmt->fetch()) {
                Utils::jsonResponse(400, '该版本已存在（组件、版本号、系统、架构组合重复）');
            }

            $sql = "INSERT INTO versions (
                component_code, component_name, version, os, arch,
                filename, file_size, file_md5, file_sha256,
                description, changelog, release_time,
                is_required_update, status
            ) VALUES (
                :component_code, :component_name, :version, :os, :arch,
                :filename, :file_size, :file_md5, :file_sha256,
                :description, :changelog, :release_time,
                :is_required_update, :status
            )";
            $stmt = $db->prepare($sql);
            $result = $stmt->execute([
                'component_code' => $data['component_code'],
                'component_name' => $data['component_name'],
                'version' => $data['version'],
                'os' => $data['os'],
                'arch' => $data['arch'],
                'filename' => $data['filename'],
                'file_size' => $data['file_size'],
                'file_md5' => $data['file_md5'],
                'file_sha256' => $data['file_sha256'] ?? null,
                'description' => $data['description'] ?? null,
                'changelog' => $data['changelog'] ?? null,
                'release_time' => $data['release_time'],
                'is_required_update' => isset($data['is_required_update']) && $data['is_required_update'] ? 1 : 0,
                'status' => isset($data['status']) && $data['status'] ? 1 : 0
            ]);

            if ($result) {
                $newId = $db->lastInsertId();
                Utils::log("添加版本: {$data['component_code']} v{$data['version']}", 'INFO');
                Utils::jsonResponse(200, '添加成功', ['id' => $newId]);
            } else {
                Utils::jsonResponse(500, '添加失败');
            }
        } elseif ($action === 'update') {
            if (empty($data['id'])) {
                Utils::jsonResponse(400, '缺少版本ID');
            }
            if (!empty($data['os'])) {
                ensureLinux($data['os']);
            }

            $checkSql = "SELECT id FROM versions WHERE id = :id";
            $checkStmt = $db->prepare($checkSql);
            $checkStmt->execute(['id' => $data['id']]);
            if (!$checkStmt->fetch()) {
                Utils::jsonResponse(404, '版本不存在');
            }

            $checkSql = "SELECT id FROM versions WHERE component_code = :component_code AND version = :version AND os = :os AND arch = :arch AND id != :id";
            $checkStmt = $db->prepare($checkSql);
            $checkStmt->execute([
                'component_code' => $data['component_code'],
                'version' => $data['version'],
                'os' => $data['os'],
                'arch' => $data['arch'],
                'id' => $data['id']
            ]);
            if ($checkStmt->fetch()) {
                Utils::jsonResponse(400, '该版本已存在（组件、版本号、系统、架构组合重复）');
            }

            $sql = "UPDATE versions SET
                component_code = :component_code,
                component_name = :component_name,
                version = :version,
                os = :os,
                arch = :arch,
                filename = :filename,
                file_size = :file_size,
                file_md5 = :file_md5,
                file_sha256 = :file_sha256,
                description = :description,
                changelog = :changelog,
                release_time = :release_time,
                is_required_update = :is_required_update,
                status = :status
            WHERE id = :id";
            $stmt = $db->prepare($sql);
            $result = $stmt->execute([
                'id' => $data['id'],
                'component_code' => $data['component_code'],
                'component_name' => $data['component_name'],
                'version' => $data['version'],
                'os' => $data['os'],
                'arch' => $data['arch'],
                'filename' => $data['filename'],
                'file_size' => $data['file_size'],
                'file_md5' => $data['file_md5'],
                'file_sha256' => $data['file_sha256'] ?? null,
                'description' => $data['description'] ?? null,
                'changelog' => $data['changelog'] ?? null,
                'release_time' => $data['release_time'],
                'is_required_update' => isset($data['is_required_update']) && $data['is_required_update'] ? 1 : 0,
                'status' => isset($data['status']) && $data['status'] ? 1 : 0
            ]);

            if ($result) {
                Utils::log("更新版本: ID={$data['id']}, {$data['component_code']} v{$data['version']}", 'INFO');
                Utils::jsonResponse(200, '更新成功');
            } else {
                Utils::jsonResponse(500, '更新失败');
            }
        } elseif ($action === 'verify') {
            if (empty($data['id'])) {
                Utils::jsonResponse(400, '缺少版本ID');
            }

            $stmt = $db->prepare("SELECT * FROM versions WHERE id = ?");
            $stmt->execute([$data['id']]);
            $version = $stmt->fetch();
            if (!$version) {
                Utils::jsonResponse(404, '版本不存在');
            }
            ensureLinux($version['os']);

            list($size, $md5, $sha256, $filePath) = verifyFileHash($version);

            $updateStmt = $db->prepare("
                UPDATE versions SET file_size = :file_size, file_md5 = :file_md5, file_sha256 = :file_sha256, updated_at = NOW()
                WHERE id = :id
            ");
            $updateStmt->execute([
                'file_size' => $size,
                'file_md5' => $md5,
                'file_sha256' => $sha256,
                'id' => $data['id']
            ]);

            Utils::jsonResponse(200, '校验完成', [
                'id' => $data['id'],
                'file' => $filePath,
                'size' => $size,
                'md5' => $md5,
                'sha256' => $sha256,
                'sizeFormatted' => Utils::formatBytes($size)
            ]);
        } elseif ($action === 'delete') {
            if (empty($data['id'])) {
                Utils::jsonResponse(400, '缺少版本ID');
            }

            $checkSql = "SELECT component_code, version FROM versions WHERE id = :id";
            $checkStmt = $db->prepare($checkSql);
            $checkStmt->execute(['id' => $data['id']]);
            $version = $checkStmt->fetch();

            if (!$version) {
                Utils::jsonResponse(404, '版本不存在');
            }

            $sql = "DELETE FROM versions WHERE id = :id";
            $stmt = $db->prepare($sql);
            $result = $stmt->execute(['id' => $data['id']]);

            if ($result) {
                Utils::log("删除版本: ID={$data['id']}, {$version['component_code']} v{$version['version']}", 'INFO');
                Utils::jsonResponse(200, '删除成功');
            } else {
                Utils::jsonResponse(500, '删除失败');
            }
        } else {
            Utils::jsonResponse(400, '无效的操作');
        }
    } else {
        Utils::jsonResponse(405, '不支持的请求方法');
    }
} catch (Exception $e) {
    Utils::log("版本管理API异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误: ' . $e->getMessage());
}
