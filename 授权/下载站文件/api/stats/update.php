<?php
// LingCDN 更新统计接口

header('Content-Type: application/json; charset=utf-8');
header('Access-Control-Allow-Origin: *');

// 日志文件
$logFile = __DIR__ . '/../logs/update-stats.log';
$logDir = dirname($logFile);

// 确保日志目录存在
if (!is_dir($logDir)) {
    @mkdir($logDir, 0755, true);
}

// 获取请求数据
$os = $_POST['os'] ?? $_GET['os'] ?? 'unknown';
$arch = $_POST['arch'] ?? $_GET['arch'] ?? 'unknown';
$components = $_POST['components'] ?? $_GET['components'] ?? 'unknown';
$updateType = $_POST['update_type'] ?? $_GET['update_type'] ?? 'unknown';

// 获取客户端信息
$ip = $_SERVER['HTTP_X_FORWARDED_FOR'] ?? $_SERVER['REMOTE_ADDR'] ?? 'unknown';
$userAgent = $_SERVER['HTTP_USER_AGENT'] ?? 'unknown';
$timestamp = date('Y-m-d H:i:s');

// 记录日志
$logData = [
    'timestamp' => $timestamp,
    'ip' => $ip,
    'os' => $os,
    'arch' => $arch,
    'components' => $components,
    'update_type' => $updateType,
    'user_agent' => $userAgent
];

$logLine = json_encode($logData, JSON_UNESCAPED_UNICODE) . "\n";

// 写入日志
@file_put_contents($logFile, $logLine, FILE_APPEND | LOCK_EX);

// 返回成功响应
echo json_encode([
    'code' => 0,
    'message' => 'success',
    'data' => [
        'received' => true,
        'timestamp' => $timestamp
    ]
], JSON_UNESCAPED_UNICODE);
