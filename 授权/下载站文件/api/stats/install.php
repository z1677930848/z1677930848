<?php
/**
 * 安装统计 API
 * 端点: /api/stats/install
 *
 * 功能：记录安装请求
 *
 * 请求参数（POST）:
 * - os: 操作系统 (linux, windows, darwin)
 * - arch: 架构 (amd64, arm64, 386)
 * - version: 安装版本
 * - install_type: 安装类型 (script, manual)
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';

// 设置响应头
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

// 处理 OPTIONS 请求
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// 只接受 POST 请求
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    Utils::jsonResponse(405, 'Method Not Allowed');
}

try {
    $db = Database::getInstance()->getConnection();

    // 获取请求参数
    $os = $_POST['os'] ?? $_GET['os'] ?? 'unknown';
    $arch = $_POST['arch'] ?? $_GET['arch'] ?? 'unknown';
    $version = $_POST['version'] ?? $_GET['version'] ?? 'unknown';
    $installType = $_POST['install_type'] ?? $_GET['install_type'] ?? 'script';

    // 获取客户端信息
    $ipAddress = Utils::getClientIP();
    $userAgent = $_SERVER['HTTP_USER_AGENT'] ?? '';

    // 简单的地理位置识别（基于IP，这里做简单处理）
    $country = null;
    $region = null;
    $city = null;

    // 可以集成第三方IP地理位置API，这里暂时留空
    // 如果需要可以使用 ip-api.com 或其他服务

    // 插入统计记录
    $sql = "INSERT INTO install_stats
            (install_time, ip_address, country, region, city, os, arch, version, install_type, user_agent)
            VALUES (NOW(), :ip_address, :country, :region, :city, :os, :arch, :version, :install_type, :user_agent)";

    $stmt = $db->prepare($sql);
    $result = $stmt->execute([
        'ip_address' => $ipAddress,
        'country' => $country,
        'region' => $region,
        'city' => $city,
        'os' => $os,
        'arch' => $arch,
        'version' => $version,
        'install_type' => $installType,
        'user_agent' => $userAgent
    ]);

    if ($result) {
        // 记录日志
        Utils::log("安装统计: IP=$ipAddress, OS=$os, ARCH=$arch, VERSION=$version", 'INFO');

        Utils::jsonResponse(200, 'success', [
            'message' => '统计记录成功',
            'install_id' => $db->lastInsertId()
        ]);
    } else {
        Utils::jsonResponse(500, '记录失败');
    }

} catch (Exception $e) {
    Utils::log("安装统计异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误');
}
