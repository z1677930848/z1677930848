<?php
/**
 * 版本检查 API (客户端引导)
 * 用于客户端检查更新和获取最新版本信息
 *
 * 参数:
 * - component: 组件代码 (admin, api, node)
 * - os: 操作系统 (linux, windows, darwin, all)
 * - arch: 架构 (amd64, arm64, 386, all)
 * - current_version: 当前版本号 (可选，用于版本比较)
 */

require_once __DIR__ . '/../../config.php';
require_once __DIR__ . '/../../database.php';
require_once __DIR__ . '/../../utils.php';

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
    $os = $_GET['os'] ?? 'linux';
    $arch = $_GET['arch'] ?? 'amd64';
    $currentVersion = $_GET['current_version'] ?? null;

    $db = Database::getInstance()->getConnection();

    // 如果没有指定 component，返回所有组件的最新版本列表（后台更新检查使用）
    if (empty($component)) {
        $validComponents = ['admin', 'api', 'node'];
        $versions = [];

        foreach ($validComponents as $comp) {
            $sql = "SELECT * FROM versions
                    WHERE component_code = :component
                    AND status = 1
                    AND (os = :os OR os = 'all')
                    AND (arch = :arch OR arch = 'all')
                    ORDER BY release_time DESC, id DESC
                    LIMIT 1";

            $stmt = $db->prepare($sql);
            $stmt->execute([
                'component' => $comp,
                'os' => $os,
                'arch' => $arch
            ]);

            $version = $stmt->fetch();
            if ($version) {
                $downloadUrl = sprintf(
                    '/updates/%s/%s/%s/%s',
                    $version['component_code'],
                    $version['os'],
                    $version['arch'],
                    $version['filename']
                );

                $versions[] = [
                    'code' => $version['component_code'],
                    'name' => $version['component_name'],
                    'version' => $version['version'],
                    'url' => $downloadUrl,
                    'size' => (int)$version['file_size'],
                    'md5' => $version['file_md5'],
                    'sha256' => $version['file_sha256'],
                    'releaseTime' => $version['release_time'],
                    'description' => $version['description'],
                    'changelog' => $version['changelog'],
                    'changelogUrl' => 'http://dl.lingcdn.cloud/changelog/v' . $version['version'] . '.php',
                    'isRequired' => (bool)$version['is_required_update']
                ];
            }
        }

        // 返回所有组件版本列表
        $response = [
            'code' => 200,
            'message' => 'success',
            'data' => [
                'host' => 'http://dl.lingcdn.cloud',
                'versions' => $versions
            ]
        ];

        Utils::log("版本列表查询: os={$os}, arch={$arch}, 找到 " . count($versions) . " 个组件", 'INFO');
        echo json_encode($response, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
        exit;
    }

    // 验证 component 参数
    $validComponents = ['admin', 'api', 'node'];
    if (!in_array($component, $validComponents)) {
        Utils::jsonResponse(400, '无效的组件类型');
    }

    // 查询最新启用的版本
    // 匹配逻辑：精确匹配 OS/ARCH 或者匹配 'all'
    $sql = "SELECT * FROM versions
            WHERE component_code = :component
            AND status = 1
            AND (os = :os OR os = 'all')
            AND (arch = :arch OR arch = 'all')
            ORDER BY release_time DESC, id DESC
            LIMIT 1";

    $stmt = $db->prepare($sql);
    $stmt->execute([
        'component' => $component,
        'os' => $os,
        'arch' => $arch
    ]);

    $version = $stmt->fetch();

    if (!$version) {
        Utils::jsonResponse(404, '未找到可用版本', [
            'component' => $component,
            'os' => $os,
            'arch' => $arch
        ]);
    }

    // 构建下载URL
    $downloadUrl = sprintf(
        '/updates/%s/%s/%s/%s',
        $version['component_code'],
        $version['os'],
        $version['arch'],
        $version['filename']
    );

    // 判断是否需要更新
    $needUpdate = false;
    if ($currentVersion) {
        $needUpdate = version_compare($version['version'], $currentVersion, '>');
    }

    // 构建响应数据
    $versionData = [
        'code' => $version['component_code'],
        'name' => $version['component_name'],
        'version' => $version['version'],
        'url' => $downloadUrl,
        'size' => (int)$version['file_size'],
        'md5' => $version['file_md5'],
        'sha256' => $version['file_sha256'],
        'releaseTime' => $version['release_time'],
        'description' => $version['description'],
        'changelog' => $version['changelog'],
        'changelogUrl' => 'http://dl.lingcdn.cloud/changelog/v' . $version['version'] . '.php',
        'isRequired' => (bool)$version['is_required_update'],
        'needUpdate' => $needUpdate,
        'downloadCount' => (int)$version['download_count']
    ];

    // 返回响应
    $response = [
        'code' => 200,
        'message' => 'success',
        'data' => [
            'host' => 'http://dl.lingcdn.cloud',
            'currentVersion' => $currentVersion,
            'latestVersion' => $version['version'],
            'needUpdate' => $needUpdate,
            'version' => $versionData
        ]
    ];

    // 记录访问日志（可选）
    Utils::log("版本检查: {$component} v{$currentVersion} -> v{$version['version']}, needUpdate: " . ($needUpdate ? 'YES' : 'NO'), 'INFO');

    echo json_encode($response, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);

} catch (Exception $e) {
    Utils::log("版本检查API异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误');
}
