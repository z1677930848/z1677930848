<?php
/**
 * 日志查询 API
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

$action = $_GET['action'] ?? 'auth_logs';

AdminAuth::requireAuth();

try {
    $db = Database::getInstance()->getConnection();

    switch ($action) {
        case 'auth_logs':
            // 授权日志
            $page = intval($_GET['page'] ?? 1);
            $pageSize = intval($_GET['pageSize'] ?? 20);
            $offset = ($page - 1) * $pageSize;

            $startDate = $_GET['start_date'] ?? '';
            $endDate = $_GET['end_date'] ?? '';
            $status = $_GET['status'] ?? '';
            $licenseId = intval($_GET['license_id'] ?? 0);
            $search = $_GET['search'] ?? '';

            $where = [];
            $params = [];

            if (!empty($startDate)) {
                $where[] = "l.created_at >= :start_date";
                $params['start_date'] = $startDate . ' 00:00:00';
            }

            if (!empty($endDate)) {
                $where[] = "l.created_at <= :end_date";
                $params['end_date'] = $endDate . ' 23:59:59';
            }

            if (!empty($status)) {
                $where[] = "l.status = :status";
                $params['status'] = $status;
            }

            if ($licenseId > 0) {
                $where[] = "l.license_id = :license_id";
                $params['license_id'] = $licenseId;
            }

            if (!empty($search)) {
                $where[] = "(l.domain LIKE :search OR l.ip_address LIKE :search OR lic.license_code LIKE :search OR lic.customer_name LIKE :search)";
                $params['search'] = "%$search%";
            }

            $whereSQL = !empty($where) ? 'WHERE ' . implode(' AND ', $where) : '';

            // 获取总数
            $stmt = $db->prepare("
                SELECT COUNT(*) as total
                FROM authorization_logs l
                LEFT JOIN licenses lic ON l.license_id = lic.id
                $whereSQL
            ");
            $stmt->execute($params);
            $total = $stmt->fetch()['total'];

            // 获取列表
            $stmt = $db->prepare("
                SELECT
                    l.*,
                    lic.license_code,
                    lic.customer_name,
                    lic.license_type
                FROM authorization_logs l
                LEFT JOIN licenses lic ON l.license_id = lic.id
                $whereSQL
                ORDER BY l.created_at DESC
                LIMIT $pageSize OFFSET $offset
            ");
            $stmt->execute($params);
            $logs = $stmt->fetchAll();

            Utils::jsonResponse(200, 'success', [
                'list' => $logs,
                'total' => $total,
                'page' => $page,
                'pageSize' => $pageSize
            ]);
            break;

        case 'download_logs':
            // 下载日志
            $page = intval($_GET['page'] ?? 1);
            $pageSize = intval($_GET['pageSize'] ?? 20);
            $offset = ($page - 1) * $pageSize;

            $startDate = $_GET['start_date'] ?? '';
            $endDate = $_GET['end_date'] ?? '';
            $versionId = intval($_GET['version_id'] ?? 0);

            $where = [];
            $params = [];

            if (!empty($startDate)) {
                $where[] = "l.download_time >= :start_date";
                $params['start_date'] = $startDate . ' 00:00:00';
            }

            if (!empty($endDate)) {
                $where[] = "l.download_time <= :end_date";
                $params['end_date'] = $endDate . ' 23:59:59';
            }

            if ($versionId > 0) {
                $where[] = "l.version_id = :version_id";
                $params['version_id'] = $versionId;
            }

            $whereSQL = !empty($where) ? 'WHERE ' . implode(' AND ', $where) : '';

            // 获取总数
            $stmt = $db->prepare("
                SELECT COUNT(*) as total
                FROM download_logs l
                $whereSQL
            ");
            $stmt->execute($params);
            $total = $stmt->fetch()['total'];

            // 获取列表
            $stmt = $db->prepare("
                SELECT
                    l.*,
                    v.component_name,
                    v.version,
                    v.os,
                    v.arch,
                    v.filename
                FROM download_logs l
                LEFT JOIN versions v ON l.version_id = v.id
                $whereSQL
                ORDER BY l.download_time DESC
                LIMIT $pageSize OFFSET $offset
            ");
            $stmt->execute($params);
            $logs = $stmt->fetchAll();

            Utils::jsonResponse(200, 'success', [
                'list' => $logs,
                'total' => $total,
                'page' => $page,
                'pageSize' => $pageSize
            ]);
            break;

        case 'stats':
            // 统计数据
            $startDate = $_GET['start_date'] ?? date('Y-m-d', strtotime('-30 days'));
            $endDate = $_GET['end_date'] ?? date('Y-m-d');

            // 授权统计
            $stmt = $db->prepare("
                SELECT
                    DATE(created_at) as date,
                    COUNT(*) as total,
                    SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success,
                    SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed
                FROM authorization_logs
                WHERE created_at >= :start_date AND created_at <= :end_date
                GROUP BY DATE(created_at)
                ORDER BY date ASC
            ");
            $stmt->execute([
                'start_date' => $startDate . ' 00:00:00',
                'end_date' => $endDate . ' 23:59:59'
            ]);
            $authStats = $stmt->fetchAll();

            // 下载统计
            $stmt = $db->prepare("
                SELECT
                    DATE(download_time) as date,
                    COUNT(*) as total
                FROM download_logs
                WHERE download_time >= :start_date AND download_time <= :end_date
                GROUP BY DATE(download_time)
                ORDER BY date ASC
            ");
            $stmt->execute([
                'start_date' => $startDate . ' 00:00:00',
                'end_date' => $endDate . ' 23:59:59'
            ]);
            $downloadStats = $stmt->fetchAll();

            // 热门版本
            $stmt = $db->prepare("
                SELECT
                    v.component_name,
                    v.version,
                    COUNT(l.id) as download_count
                FROM download_logs l
                LEFT JOIN versions v ON l.version_id = v.id
                WHERE l.download_time >= :start_date AND l.download_time <= :end_date
                GROUP BY l.version_id
                ORDER BY download_count DESC
                LIMIT 10
            ");
            $stmt->execute([
                'start_date' => $startDate . ' 00:00:00',
                'end_date' => $endDate . ' 23:59:59'
            ]);
            $topVersions = $stmt->fetchAll();

            // 今日统计
            $stmt = $db->query("
                SELECT
                    (SELECT COUNT(*) FROM authorization_logs WHERE DATE(created_at) = CURDATE()) as today_auth,
                    (SELECT COUNT(*) FROM authorization_logs WHERE DATE(created_at) = CURDATE() AND status = 'success') as today_auth_success,
                    (SELECT COUNT(*) FROM download_logs WHERE DATE(download_time) = CURDATE()) as today_download,
                    (SELECT COUNT(*) FROM licenses WHERE DATE(created_at) = CURDATE()) as today_issue,
                    (SELECT COUNT(*) FROM download_logs) as total_downloads,
                    (SELECT MIN(created_at) FROM authorization_logs) as first_auth_time,
                    (SELECT MAX(created_at) FROM authorization_logs) as last_auth_time,
                    (SELECT COUNT(*) FROM licenses WHERE status = 1) as active_licenses,
                    (SELECT COUNT(*) FROM licenses) as total_licenses
            ");
            $todayStats = $stmt->fetch();

            Utils::jsonResponse(200, 'success', [
                'auth_stats' => $authStats,
                'download_stats' => $downloadStats,
                'top_versions' => $topVersions,
                'today' => $todayStats
            ]);
            break;

        case 'system_logs':
            // 系统日志文件
            $date = $_GET['date'] ?? date('Y-m-d');
            $logFile = LOG_PATH . $date . '.log';

            if (!file_exists($logFile)) {
                Utils::jsonResponse(404, '日志文件不存在');
            }

            $lines = file($logFile, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
            $lines = array_reverse($lines); // 最新的在前

            $page = intval($_GET['page'] ?? 1);
            $pageSize = intval($_GET['pageSize'] ?? 50);
            $total = count($lines);

            $start = ($page - 1) * $pageSize;
            $logs = array_slice($lines, $start, $pageSize);

            Utils::jsonResponse(200, 'success', [
                'list' => $logs,
                'total' => $total,
                'page' => $page,
                'pageSize' => $pageSize
            ]);
            break;

        default:
            Utils::jsonResponse(400, '无效的操作');
    }

} catch (Exception $e) {
    Utils::log("日志查询异常: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误');
}
