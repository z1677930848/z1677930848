<?php
require_once __DIR__ . '/../config.php';
require_once __DIR__ . '/../database.php';

$version = '1.0.17';
$db = Database::getInstance()->getConnection();
$stmt = $db->prepare("SELECT * FROM versions WHERE version = ? AND component_code = 'admin' LIMIT 1");
$stmt->execute([$version]);
$data = $stmt->fetch();

$changelog = $data['changelog'] ?? '暂无更新说明';
$description = $data['description'] ?? '';
$releaseTime = $data['release_time'] ?? date('Y-m-d');
?>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LingCDN v<?= htmlspecialchars($version) ?> 更新日志</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; line-height: 1.6; }
        h1 { color: #667eea; border-bottom: 2px solid #667eea; padding-bottom: 10px; }
        .version { color: #999; font-size: 14px; }
        .description { background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .changelog { white-space: pre-wrap; }
    </style>
</head>
<body>
    <h1>LingCDN v<?= htmlspecialchars($version) ?></h1>
    <p class="version">发布日期: <?= htmlspecialchars(substr($releaseTime, 0, 10)) ?></p>

    <?php if ($description): ?>
    <div class="description"><?= nl2br(htmlspecialchars($description)) ?></div>
    <?php endif; ?>

    <div class="changelog"><?= nl2br(htmlspecialchars($changelog)) ?></div>
</body>
</html>