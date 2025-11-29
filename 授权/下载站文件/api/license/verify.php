<?php
require_once dirname(__DIR__, 2) . '/config.php';

header('Content-Type: application/json');

$code = isset($_GET['code']) ? trim($_GET['code']) : '';
$systemToken = isset($_GET['token']) ? trim($_GET['token']) : '';

if (empty($code)) {
    echo json_encode([
        'code' => 400,
        'message' => '请提供授权码',
        'data' => null
    ]);
    exit;
}

// 连接数据库
try {
    $dsn = "mysql:host=" . DB_HOST . ";dbname=" . DB_NAME . ";charset=" . DB_CHARSET;
    $pdo = new PDO($dsn, DB_USER, DB_PASS);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    $pdo->exec("SET time_zone = '+8:00'");
} catch (PDOException $e) {
    echo json_encode([
        'code' => 500,
        'message' => '数据库连接失败',
        'data' => null
    ]);
    exit;
}

// 查询授权信息
try {
    $stmt = $pdo->prepare("
        SELECT
            id,
            license_code,
            license_type,
            customer_name,
            customer_email,
            company_name,
            max_devices,
            max_nodes,
            bound_system_tokens,
            issue_time,
            expire_time,
            first_activated_at,
            status
        FROM licenses
        WHERE license_code = :code
    ");

    $stmt->execute(['code' => $code]);
    $license = $stmt->fetch(PDO::FETCH_ASSOC);

    if (!$license) {
        echo json_encode([
            'code' => 404,
            'message' => '授权码不存在',
            'data' => null
        ]);
        exit;
    }

    // 检查授权状态
    if ($license['status'] != 1) {
        echo json_encode([
            'code' => 403,
            'message' => '授权已被禁用',
            'data' => null
        ]);
        exit;
    }

    // 检查是否过期
    if ($license['expire_time']) {
        $expireTime = strtotime($license['expire_time']);
        if (time() > $expireTime) {
            echo json_encode([
                'code' => 410,
                'message' => '授权已过期',
                'data' => null
            ]);
            exit;
        }
    }

    // 更新最后检查时间和检查次数
    $updateStmt = $pdo->prepare("
        UPDATE licenses
        SET last_check_time = NOW(),
            check_count = check_count + 1
        WHERE license_code = :code
    ");
    $updateStmt->execute(['code' => $code]);

    // 一机一码：获取客户端IP
    $clientIP = $_SERVER['REMOTE_ADDR'] ?? '';
    if (isset($_SERVER['HTTP_X_FORWARDED_FOR'])) {
        $clientIP = explode(',', $_SERVER['HTTP_X_FORWARDED_FOR'])[0];
    }

    // 获取已绑定的IP
    $boundIP = $license['bound_system_tokens'] ? json_decode($license['bound_system_tokens'], true)[0] ?? null : null;

    // 检查IP绑定
    if ($boundIP && $boundIP !== $clientIP) {
        echo json_encode([
            'code' => 403,
            'message' => '此授权码已绑定其他IP地址',
            'data' => null
        ]);
        exit;
    }

    // 首次绑定IP
    if (!$boundIP && !empty($clientIP)) {
        $bindStmt = $pdo->prepare("
            UPDATE licenses
            SET bound_system_tokens = :tokens,
                first_activated_at = NOW()
            WHERE license_code = :code
        ");
        $bindStmt->execute([
            'tokens' => json_encode([$clientIP]),
            'code' => $code
        ]);
    }

    // 记录授权验证日志
    $logStmt = $pdo->prepare("
        INSERT INTO authorization_logs
        (license_id, system_token, domain, ip_address, status, user_agent, created_at)
        VALUES (:license_id, :token, :domain, :ip, :status, :ua, NOW())
    ");
    $logStmt->execute([
        'license_id' => $license['id'],
        'token' => '',
        'domain' => $_SERVER['HTTP_HOST'] ?? '',
        'ip' => $_SERVER['REMOTE_ADDR'] ?? '',
        'status' => 'success',
        'ua' => $_SERVER['HTTP_USER_AGENT'] ?? ''
    ]);

    // 构建响应数据
    $responseData = [
        'type' => $license['license_type'],
        'company' => $license['company_name'] ?: $license['customer_name'],
        'email' => $license['customer_email'] ?: '',
        'maxNodes' => (int)$license['max_nodes'],
        'expireAt' => $license['expire_time'] ?: '2099-12-31 23:59:59',
        'createdAt' => $license['issue_time'],
        'activatedAt' => $license['first_activated_at'] ?: '',
        'boundIP' => $boundIP ?: ''
    ];

    echo json_encode([
        'code' => 200,
        'message' => '授权验证成功',
        'data' => $responseData
    ]);

} catch (PDOException $e) {
    error_log('License verify error: ' . $e->getMessage());
    echo json_encode([
        'code' => 500,
        'message' => '授权验证失败',
        'data' => null
    ]);
}
