<?php
/**
 * 许可证管理 API
 */

require_once '../../config.php';
require_once '../../database.php';
require_once '../../utils.php';
require_once __DIR__ . '/auth_guard.php';

header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization, X-Admin-Token');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

$method = $_SERVER['REQUEST_METHOD'];
$action = $_GET['action'] ?? '';

// 绠＄悊鍛樺瘑鎶や护鏍￠獙
AdminAuth::requireAuth();

try {
    $db = Database::getInstance()->getConnection();

    switch ($action) {
        case 'list':
            // 获取许可证列表
            $page = intval($_GET['page'] ?? 1);
            $pageSize = intval($_GET['pageSize'] ?? 20);
            $offset = ($page - 1) * $pageSize;

            $search = $_GET['search'] ?? '';
            $status = $_GET['status'] ?? '';

            $where = [];
            $params = [];

            if (!empty($search)) {
                $where[] = "(license_code LIKE :search OR customer_name LIKE :search OR customer_email LIKE :search)";
                $params['search'] = "%$search%";
            }

            if ($status !== '') {
                $where[] = "status = :status";
                $params['status'] = $status;
            }

            $whereSQL = !empty($where) ? 'WHERE ' . implode(' AND ', $where) : '';

            // 获取总数
            $stmt = $db->prepare("SELECT COUNT(*) as total FROM licenses $whereSQL");
            $stmt->execute($params);
            $total = $stmt->fetch()['total'];

            // 获取列表
            $stmt = $db->prepare("
                SELECT *
                FROM licenses
                $whereSQL
                ORDER BY created_at DESC
                LIMIT $pageSize OFFSET $offset
            ");
            $stmt->execute($params);
            $licenses = $stmt->fetchAll();

            foreach ($licenses as &$license) {
                $license['bound_system_tokens'] = json_decode($license['bound_system_tokens'], true);
            }

            Utils::jsonResponse(200, 'success', [
                'list' => $licenses,
                'total' => $total,
                'page' => $page,
                'pageSize' => $pageSize
            ]);
            break;

        case 'add':
            // 添加许可证
            $input = file_get_contents('php://input');
            $data = json_decode($input, true);

            // 如果license_code为空或不存在，自动生成
            $licenseCode = (!empty($data['license_code']) && trim($data['license_code']) !== '')
                ? trim($data['license_code'])
                : Utils::generateLicenseCode();

            $licenseType = $data['license_type'] ?? 'online';
            $customerName = $data['customer_name'] ?? '';
            $customerEmail = $data['customer_email'] ?? '';
            $customerPhone = $data['customer_phone'] ?? '';
            $companyName = $data['company_name'] ?? '';
            $maxDevices = 1; // 一机一码，固定为1
            $maxNodes = intval($data['max_nodes'] ?? 0);

            // 处理到期时间：如果为空或null，设置为null
            $expireTime = null;
            if (isset($data['expire_time']) && !empty($data['expire_time']) && trim($data['expire_time']) !== '') {
                $expireTime = $data['expire_time'];
            }

            $remark = $data['remark'] ?? '';

            if (empty($customerName)) {
                Utils::jsonResponse(400, '客户名称不能为空');
            }

            $stmt = $db->prepare("
                INSERT INTO licenses
                (license_code, license_type, customer_name, customer_email, customer_phone, company_name,
                 max_devices, max_nodes, issue_time, expire_time, remark, status)
                VALUES
                (:license_code, :license_type, :customer_name, :customer_email, :customer_phone, :company_name,
                 :max_devices, :max_nodes, NOW(), :expire_time, :remark, 1)
            ");

            $stmt->execute([
                'license_code' => $licenseCode,
                'license_type' => $licenseType,
                'customer_name' => $customerName,
                'customer_email' => $customerEmail,
                'customer_phone' => $customerPhone,
                'company_name' => $companyName,
                'max_devices' => $maxDevices,
                'max_nodes' => $maxNodes,
                'expire_time' => $expireTime,
                'remark' => $remark
            ]);

            Utils::log("添加许可证: $licenseCode - 客户: $customerName", 'INFO');

            Utils::jsonResponse(200, '添加成功', [
                'license_code' => $licenseCode,
                'id' => $db->lastInsertId()
            ]);
            break;

        case 'batch_add':
            // 批量添加许可证
            $input = file_get_contents('php://input');
            $data = json_decode($input, true);

            $count = intval($data['count'] ?? 1);
            $licenseType = $data['license_type'] ?? 'online';
            $customerNamePrefix = $data['customer_name_prefix'] ?? '客户';
            $companyName = $data['company_name'] ?? '';
            $maxNodes = intval($data['max_nodes'] ?? 0);

            // 处理到期时间
            $expireTime = null;
            if (isset($data['expire_time']) && !empty($data['expire_time']) && trim($data['expire_time']) !== '') {
                $expireTime = $data['expire_time'];
            }

            $remark = $data['remark'] ?? '';

            if ($count < 1 || $count > 1000) {
                Utils::jsonResponse(400, '批量数量必须在1-1000之间');
            }

            $stmt = $db->prepare("
                INSERT INTO licenses
                (license_code, license_type, customer_name, customer_email, customer_phone, company_name,
                 max_devices, max_nodes, issue_time, expire_time, remark, status)
                VALUES
                (:license_code, :license_type, :customer_name, :customer_email, :customer_phone, :company_name,
                 :max_devices, :max_nodes, NOW(), :expire_time, :remark, 1)
            ");

            $successCount = 0;
            $failedCount = 0;
            $generatedLicenses = [];

            for ($i = 1; $i <= $count; $i++) {
                try {
                    $licenseCode = Utils::generateLicenseCode();
                    $customerName = $customerNamePrefix . sprintf('%04d', $i);

                    $stmt->execute([
                        'license_code' => $licenseCode,
                        'license_type' => $licenseType,
                        'customer_name' => $customerName,
                        'customer_email' => '',
                        'customer_phone' => '',
                        'company_name' => $companyName,
                        'max_devices' => 1,
                        'max_nodes' => $maxNodes,
                        'expire_time' => $expireTime,
                        'remark' => $remark
                    ]);

                    $generatedLicenses[] = [
                        'license_code' => $licenseCode,
                        'customer_name' => $customerName,
                        'id' => $db->lastInsertId()
                    ];
                    $successCount++;
                } catch (Exception $e) {
                    $failedCount++;
                    Utils::log("批量添加许可证失败 ($i/$count): " . $e->getMessage(), 'ERROR');
                }
            }

            Utils::log("批量添加许可证: 成功 $successCount 个, 失败 $failedCount 个", 'INFO');

            Utils::jsonResponse(200, "批量添加完成: 成功 $successCount 个, 失败 $failedCount 个", [
                'success_count' => $successCount,
                'failed_count' => $failedCount,
                'licenses' => $generatedLicenses
            ]);
            break;

        case 'update':
            // 更新许可证
            $input = file_get_contents('php://input');
            $data = json_decode($input, true);

            $id = intval($data['id'] ?? 0);
            if ($id <= 0) {
                Utils::jsonResponse(400, '无效的许可证ID');
            }

            $updates = [];
            $params = ['id' => $id];

            $allowedFields = [
                'license_type', 'customer_name', 'customer_email', 'customer_phone',
                'company_name', 'max_nodes', 'expire_time', 'remark', 'status'
            ];

            foreach ($allowedFields as $field) {
                if (isset($data[$field])) {
                    $updates[] = "$field = :$field";
                    $params[$field] = $data[$field];
                }
            }

            if (empty($updates)) {
                Utils::jsonResponse(400, '没有要更新的数据');
            }

            $updateSQL = implode(', ', $updates);
            $stmt = $db->prepare("UPDATE licenses SET $updateSQL WHERE id = :id");
            $stmt->execute($params);

            Utils::log("更新许可证ID: $id", 'INFO');

            Utils::jsonResponse(200, '更新成功');
            break;

        case 'reset_binding':
            // 重置绑定状态
            $id = intval($_GET['id'] ?? 0);
            if ($id <= 0) {
                Utils::jsonResponse(400, '无效的许可证ID');
            }

            $stmt = $db->prepare("UPDATE licenses SET bound_system_tokens = NULL WHERE id = :id");
            $stmt->execute(['id' => $id]);

            Utils::log("重置许可证绑定ID: $id", 'INFO');
            Utils::jsonResponse(200, '重置成功');
            break;

        case 'delete':
            // 删除许可证
            $id = intval($_GET['id'] ?? 0);
            if ($id <= 0) {
                Utils::jsonResponse(400, '无效的许可证ID');
            }

            $stmt = $db->prepare("DELETE FROM licenses WHERE id = :id");
            $stmt->execute(['id' => $id]);

            Utils::log("删除许可证ID: $id", 'WARNING');

            Utils::jsonResponse(200, '删除成功');
            break;

        case 'toggle_status':
            // 切换许可证状态
            $id = intval($_GET['id'] ?? 0);
            if ($id <= 0) {
                Utils::jsonResponse(400, '无效的许可证ID');
            }

            $stmt = $db->prepare("UPDATE licenses SET status = 1 - status WHERE id = :id");
            $stmt->execute(['id' => $id]);

            Utils::log("切换许可证状态ID: $id", 'INFO');

            Utils::jsonResponse(200, '状态更新成功');
            break;

        case 'detail':
            // 获取详情
            $id = intval($_GET['id'] ?? 0);
            if ($id <= 0) {
                Utils::jsonResponse(400, '无效的许可证ID');
            }

            $stmt = $db->prepare("SELECT * FROM licenses WHERE id = :id");
            $stmt->execute(['id' => $id]);
            $license = $stmt->fetch();

            if (!$license) {
                Utils::jsonResponse(404, '许可证不存在');
            }

            $license['allowed_domains'] = json_decode($license['allowed_domains'], true);
            $license['bound_system_tokens'] = json_decode($license['bound_system_tokens'], true);
            $license['features'] = json_decode($license['features'], true);

            Utils::jsonResponse(200, 'success', $license);
            break;

        default:
            Utils::jsonResponse(400, '无效的操作');
    }

} catch (Exception $e) {
    Utils::log("许可证管理异常: " . $e->getMessage(), 'ERROR');
    // 开发模式下显示详细错误信息
    $errorMsg = '服务器内部错误';
    if (defined('DEBUG') && DEBUG) {
        $errorMsg .= ': ' . $e->getMessage();
    }
    Utils::jsonResponse(500, $errorMsg, [
        'error' => $e->getMessage(),
        'file' => $e->getFile(),
        'line' => $e->getLine()
    ]);
}
