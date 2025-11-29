<?php
/**
 * 鎺堟潈楠岃瘉 API
 * 绔偣: /api/authorization.php
 *
 * 鎺ユ敹鍙傛暟锛圥OST JSON锛?
 * - license_code: 璁稿彲璇佸瘑閽?
 * - system_token: 绯荤粺浠ょ墝锛堟満鍣ㄧ爜锛?
 * - domain: 鍩熷悕
 */

require_once '../config.php';
require_once '../database.php';
require_once '../utils.php';

// 璁剧疆鍝嶅簲澶?
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

// 澶勭悊 OPTIONS 璇锋眰
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// 鍙帴鍙?POST 璇锋眰
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    Utils::jsonResponse(405, 'Method Not Allowed');
}

// 鑾峰彇璇锋眰鏁版嵁
$input = file_get_contents('php://input');
$data = json_decode($input, true);

if (!$data) {
    Utils::jsonResponse(400, '鏃犳晥鐨勮姹傛暟鎹?);
}

// 楠岃瘉蹇呭～瀛楁
$licenseCode = $data['license_code'] ?? '';
$systemToken = $data['system_token'] ?? '';
$domain = $data['domain'] ?? '';

if (empty($licenseCode) || empty($systemToken) || empty($domain)) {
    Utils::jsonResponse(400, '缂哄皯蹇呭～鍙傛暟');
}

// 楠岃瘉璁稿彲璇佹牸寮?
if (!Utils::validateLicenseFormat($licenseCode)) {
    Utils::jsonResponse(400, '璁稿彲璇佹牸寮忔棤鏁?);
}

try {
    $db = Database::getInstance()->getConnection();

    // 鏌ヨ璁稿彲璇?
    $stmt = $db->prepare("
        SELECT * FROM licenses
        WHERE license_code = :license_code
        AND status = 1
    ");
    $stmt->execute(['license_code' => $licenseCode]);
    $license = $stmt->fetch();

    if (!$license) {
        Utils::log("鎺堟潈澶辫触: 璁稿彲璇佷笉瀛樺湪鎴栧凡绂佺敤 - $licenseCode - IP: " . Utils::getClientIP(), 'WARNING');
        Utils::jsonResponse(401, '璁稿彲璇佹棤鏁堟垨宸茶绂佺敤');
    }

    // 妫€鏌ヨ繃鏈熸椂闂?
    if ($license['expire_time'] && strtotime($license['expire_time']) < time()) {
        Utils::log("鎺堟潈澶辫触: 璁稿彲璇佸凡杩囨湡 - $licenseCode", 'WARNING');
        Utils::jsonResponse(401, '璁稿彲璇佸凡杩囨湡');
    }

    // 妫€鏌ュ煙鍚嶇粦瀹?
    if (!empty($license['allowed_domains'])) {
        $allowedDomains = json_decode($license['allowed_domains'], true);
        if (!in_array($domain, $allowedDomains) && !in_array('*', $allowedDomains)) {
            Utils::log("鎺堟潈澶辫触: 鍩熷悕鏈巿鏉?- $domain - 璁稿彲璇? $licenseCode", 'WARNING');
            Utils::jsonResponse(401, '璇ュ煙鍚嶆湭鑾峰緱鎺堟潈');
        }
    }

    // 妫€鏌ユ満鍣ㄧ爜缁戝畾
    $boundTokens = json_decode($license['bound_system_tokens'], true);
    if (!is_array($boundTokens)) {
        $boundTokens = [];
    }

    $alreadyBound = in_array($systemToken, $boundTokens, true);
    $maxDevices = isset($license['max_devices']) ? (int)$license['max_devices'] : 1;

    if (!$alreadyBound) {
        if ($maxDevices > 0 && count($boundTokens) >= $maxDevices) {
            Utils::log("鎺堟潈澶辫触: 宸茶揪鍒版渶澶ц澶囨暟 - 璁稿彲璇? $licenseCode", 'WARNING');
            Utils::jsonResponse(401, '宸茶揪鍒版渶澶ф巿鏉冭澶囨暟');
        }

        $boundTokens[] = $systemToken;
        $stmt = $db->prepare("
            UPDATE licenses
            SET bound_system_tokens = :tokens,
                last_check_time = NOW(),
                check_count = check_count + 1
            WHERE id = :id
        ");
        $stmt->execute([
            'tokens' => json_encode(array_values(array_unique($boundTokens)), JSON_UNESCAPED_UNICODE),
            'id' => $license['id']
        ]);

        Utils::log("鏂拌澶囩粦瀹氭垚鍔?- 璁稿彲璇? $licenseCode - Token: $systemToken", 'INFO');
    } else {
        $stmt = $db->prepare("
            UPDATE licenses
            SET last_check_time = NOW(),
                check_count = check_count + 1
            WHERE id = :id
        ");
        $stmt->execute(['id' => $license['id']]);
    }

    // 璁板綍楠岃瘉鏃ュ織
    $stmt = $db->prepare("
        INSERT INTO authorization_logs
        (license_id, system_token, domain, ip_address, status, created_at)
        VALUES (:license_id, :system_token, :domain, :ip, 'success', NOW())
    ");
    $stmt->execute([
        'license_id' => $license['id'],
        'system_token' => $systemToken,
        'domain' => $domain,
        'ip' => Utils::getClientIP()
    ]);

    // 杩斿洖鎴愬姛鍝嶅簲
    Utils::jsonResponse(200, '鎺堟潈楠岃瘉鎴愬姛', [
        'license_type' => $license['license_type'],
        'expire_time' => $license['expire_time'],
        'customer_name' => $license['customer_name'],
        'features' => json_decode($license['features'], true)
    ]);

} catch (Exception $e) {
    Utils::log("鎺堟潈楠岃瘉寮傚父: " . $e->getMessage(), 'ERROR');
    Utils::jsonResponse(500, '鏈嶅姟鍣ㄥ唴閮ㄩ敊璇?);
}
