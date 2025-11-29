<?php
/**
 * 版本列表 API (客户端)
 * 用于获取所有可用版本列表，支持版本比较
 *
 * 参数:
 * - component: 组件代码 (admin, api, node) - 可选，不传则返回所有组件
 * - os: 操作系统过滤 - 可选
 * - arch: 架构过滤 - 可选
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';

// 设置响应头
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');
header('Content-Type: application/json; charset=utf-8');

// 处理 OPTIONS 请求
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

try {
    // 获取请求参数
    $component = $_GET['component'] ?? null;
    $os = $_GET['os'] ?? null;
    $arch = $_GET['arch'] ?? null;

    $db = Database::getInstance()->getConnection();

    // 构建查询条件
    $conditions = ['status = 1'];
    $params = [];

    if ($component) {
        $conditions[] = 'component_code = :component';
        $params['component'] = $component;
    }

    if ($os) {
        $conditions[] = '(os = :os OR os = \'all\')';
        $params['os'] = $os;
    }

    if ($arch) {
        $conditions[] = '(arch = :arch OR arch = \'all\')';
        $params['arch'] = $arch;
    }

    $whereClause = implode(' AND ', $conditions);

    // 查询所有符合条件的版本
    $sql = "SELECT * FROM versions
            WHERE {$whereClause}
            ORDER BY component_code ASC, release_time DESC, id DESC";

    $stmt = $db->prepare($sql);
    $stmt->execute($params);
    $versions = $stmt->fetchAll();

    // 按组件分组
    $groupedVersions = [];
    $latestVersions = [];

    foreach ($versions as $version) {
        $componentCode = $version['component_code'];

        // 构建下载URL
        $downloadUrl = sprintf(
            '/updates/%s/%s/%s/%s',
            $version['component_code'],
            $version['os'],
            $version['arch'],
            $version['filename']
        );

        $versionData = [
            'id' => (int)$version['id'],
            'code' => $version['component_code'],
            'name' => $version['component_name'],
            'version' => $version['version'],
            'os' => $version['os'],
            'arch' => $version['arch'],
            'url' => $downloadUrl,
            'filename' => $version['filename'],
            'size' => (int)$version['file_size'],
            'sizeFormatted' => Utils::formatBytes($version['file_size']),
            'md5' => $version['file_md5'],
            'sha256' => $version['file_sha256'],
            'releaseTime' => $version['release_time'],
            'description' => $version['description'],
            'changelog' => $version['changelog'],
            'isRequired' => (bool)$version['is_required_update'],
            'downloadCount' => (int)$version['download_count']
        ];

        // 添加到分组
        if (!isset($groupedVersions[$componentCode])) {
            $groupedVersions[$componentCode] = [];
        }
        $groupedVersions[$componentCode][] = $versionData;

        // 记录最新版本（每个组件的第一个版本就是最新的，因为已经按时间倒序排列）
        if (!isset($latestVersions[$componentCode])) {
            $latestVersions[$componentCode] = $versionData;
        }
    }

    // 构建响应
    $response = [
        'code' => 200,
        'message' => 'success',
        'data' => [
            'host' => 'https://dl.lingcdn.cloud',
            'total' => count($versions),
            'components' => array_keys($groupedVersions),
            'latestVersions' => $latestVersions,
            'allVersions' => $groupedVersions,
            'serverTime' => date('Y-m-d H:i:s')
        ]
    ];

    echo json_encode($response, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);

} catch (Exception $e) {
    Utils::log("版本列表API异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误');
}
