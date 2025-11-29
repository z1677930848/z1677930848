<?php
/**
 * 管理员认证入口
 */

$rootPath = dirname(__DIR__, 2);
require_once $rootPath . '/config.php';
require_once $rootPath . '/utils.php';
require_once __DIR__ . '/auth_guard.php';

header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, GET, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization, X-Admin-Token');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

$method = $_SERVER['REQUEST_METHOD'];
$action = $_GET['action'] ?? ($method === 'POST' ? ($_POST['action'] ?? 'login') : 'verify');

try {
    switch ($action) {
        case 'login':
            if ($method !== 'POST') {
                Utils::jsonResponse(405, 'Method Not Allowed');
            }

            $payload = json_decode(file_get_contents('php://input'), true);
            if (!is_array($payload)) {
                $payload = $_POST;
            }

            $username = trim($payload['username'] ?? '');
            $password = $payload['password'] ?? '';

            if ($username === '' || $password === '') {
                Utils::jsonResponse(400, '请输入账号和密码');
            }

            $result = AdminAuth::login($username, $password);
            Utils::jsonResponse(200, '登录成功', $result);
            break;

        case 'verify':
            $payload = AdminAuth::requireAuth();
            Utils::jsonResponse(200, 'authorized', [
                'username' => $payload['sub'] ?? ADMIN_USERNAME,
                'expire_at' => isset($payload['exp']) ? date('c', $payload['exp']) : null,
                'issued_at' => isset($payload['iat']) ? date('c', $payload['iat']) : null,
            ]);
            break;

        case 'logout':
            // token 只存在客户端，本接口主要用于前端保持一致
            Utils::jsonResponse(200, 'logout success');
            break;

        default:
            Utils::jsonResponse(400, '无效的操作');
    }
} catch (Exception $e) {
    Utils::log('Admin auth exception: ' . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '服务器内部错误');
}
