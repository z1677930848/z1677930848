<?php
/**
 * 查询安装统计 API
 * 端点: /api/stats/query
 *
 * 功能：查询安装统计数据
 *
 * 请求参数（GET）:
 * - type: 查询类型 (total, today, recent, trend)
 * - days: 天数 (用于 recent 和 trend)
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';
require_once __DIR__ . '/../admin/auth_guard.php';

// 设置响应头
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization, X-Admin-Token');

// 处理 OPTIONS 请求
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// 只接受 GET 请求
if ($_SERVER['REQUEST_METHOD'] !== 'GET') {
    Utils::jsonResponse(405, 'Method Not Allowed');
}

try {
    AdminAuth::requireAuth();
    $db = Database::getInstance()->getConnection();

    $type = $_GET['type'] ?? 'total';
    $days = (int)($_GET['days'] ?? 7);

    $result = [];

    switch ($type) {
        case 'total':
            // 总安装次数
            $stmt = $db->query("SELECT COUNT(*) as total FROM install_stats");
            $result['total'] = $stmt->fetch()['total'];
            break;

        case 'today':
            // 今日安装次数
            $stmt = $db->query("SELECT COUNT(*) as total FROM install_stats WHERE DATE(install_time) = CURDATE()");
            $result['today'] = $stmt->fetch()['total'];
            break;

        case 'recent':
            // 最近N天
            $stmt = $db->prepare("SELECT COUNT(*) as total FROM install_stats WHERE install_time >= DATE_SUB(NOW(), INTERVAL :days DAY)");
            $stmt->execute(['days' => $days]);
            $result['recent'] = $stmt->fetch()['total'];
            $result['days'] = $days;
            break;

        case 'trend':
            // 趋势数据（每天的安装量）
            $stmt = $db->prepare("
                SELECT DATE(install_time) as date, COUNT(*) as count
                FROM install_stats
                WHERE install_time >= DATE_SUB(NOW(), INTERVAL :days DAY)
                GROUP BY DATE(install_time)
                ORDER BY date ASC
            ");
            $stmt->execute(['days' => $days]);
            $result['trend'] = $stmt->fetchAll();
            $result['days'] = $days;
            break;

        case 'stats':
            // 综合统计
            // 总数
            $stmt = $db->query("SELECT COUNT(*) as total FROM install_stats");
            $result['total'] = $stmt->fetch()['total'];

            // 今日
            $stmt = $db->query("SELECT COUNT(*) as total FROM install_stats WHERE DATE(install_time) = CURDATE()");
            $result['today'] = $stmt->fetch()['total'];

            // 本周
            $stmt = $db->query("SELECT COUNT(*) as total FROM install_stats WHERE YEARWEEK(install_time) = YEARWEEK(NOW())");
            $result['week'] = $stmt->fetch()['total'];

            // 本月
            $stmt = $db->query("SELECT COUNT(*) as total FROM install_stats WHERE YEAR(install_time) = YEAR(NOW()) AND MONTH(install_time) = MONTH(NOW())");
            $result['month'] = $stmt->fetch()['total'];

            // 按操作系统
            $stmt = $db->query("SELECT os, COUNT(*) as count FROM install_stats GROUP BY os ORDER BY count DESC");
            $result['by_os'] = $stmt->fetchAll();

            // 按架构
            $stmt = $db->query("SELECT arch, COUNT(*) as count FROM install_stats GROUP BY arch ORDER BY count DESC");
            $result['by_arch'] = $stmt->fetchAll();

            // 最近7天趋势
            $stmt = $db->query("
                SELECT DATE(install_time) as date, COUNT(*) as count
                FROM install_stats
                WHERE install_time >= DATE_SUB(NOW(), INTERVAL 7 DAY)
                GROUP BY DATE(install_time)
                ORDER BY date ASC
            ");
            $result['trend_7days'] = $stmt->fetchAll();

            break;

        default:
            Utils::jsonResponse(400, '无效的查询类型');
    }

    Utils::jsonResponse(200, 'success', $result);

} catch (Exception $e) {
    Utils::log("查询统计异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误');
}
