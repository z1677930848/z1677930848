<?php
// 数据库诊断和修复工具
// 访问: http://dl.skcdn.cn/fix-database.php

header('Content-Type: text/html; charset=utf-8');

// 从表单获取数据库配置，或使用默认值
$dbUser = $_POST['db_user'] ?? 'edge';
$dbPass = $_POST['db_pass'] ?? '123456';
?>
<!DOCTYPE html>
<html>
<head>
    <title>数据库诊断修复工具</title>
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; background: #f5f5f5; }
        .container { max-width: 900px; margin: 0 auto; background: white; padding: 20px; border-radius: 5px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; }
        .log { background: #000; color: #0f0; padding: 15px; border-radius: 3px; margin: 20px 0; max-height: 600px; overflow-y: auto; font-family: monospace; font-size: 14px; }
        .success { color: #0f0; }
        .error { color: #f00; }
        .warning { color: #ff0; }
        .info { color: #0ff; }
        button { background: #007bff; color: white; padding: 12px 24px; border: none; border-radius: 3px; cursor: pointer; font-size: 16px; margin: 5px; }
        button:hover { background: #0056b3; }
        .btn-warning { background: #ffc107; color: #000; }
        .btn-warning:hover { background: #e0a800; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
    </style>
</head>
<body>
<div class="container">
    <h1>🔧 数据库诊断修复工具</h1>
    <p>自动检测并修复数据库配置问题</p>

    <form method="post">
        <div style="margin: 15px 0; padding: 15px; background: #f8f9fa; border-radius: 3px;">
            <label style="display: block; margin-bottom: 10px;">
                <strong>数据库用户名:</strong><br>
                <input type="text" name="db_user" value="<?php echo htmlspecialchars($dbUser); ?>" style="padding: 8px; width: 200px; border: 1px solid #ccc; border-radius: 3px;">
            </label>
            <label style="display: block; margin-bottom: 10px;">
                <strong>数据库密码:</strong><br>
                <input type="password" name="db_pass" value="<?php echo htmlspecialchars($dbPass); ?>" style="padding: 8px; width: 200px; border: 1px solid #ccc; border-radius: 3px;">
            </label>
        </div>
        <button type="submit" name="action" value="diagnose">诊断问题</button>
        <button type="submit" name="action" value="fix" class="btn-warning" onclick="return confirm('确定要自动修复吗？')">自动修复</button>
    </form>

    <div class="log">
<?php
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $action = $_POST['action'] ?? '';

    echo "<span class='info'>========================================</span>\n";
    echo "<span class='info'>开始诊断...</span>\n";
    echo "<span class='info'>========================================</span>\n\n";

    // 连接 MySQL
    $mysqli = new mysqli('127.0.0.1', $dbUser, $dbPass);

    if ($mysqli->connect_error) {
        echo "<span class='error'>✗ MySQL 连接失败: " . htmlspecialchars($mysqli->connect_error) . "</span>\n";
        exit;
    }

    echo "<span class='success'>✓ MySQL 连接成功</span>\n\n";

    // 1. 检查数据库
    echo "<span class='info'>[1] 检查数据库...</span>\n";
    $result = $mysqli->query("SHOW DATABASES LIKE 'edge%'");
    $databases = [];

    if ($result) {
        while ($row = $result->fetch_array()) {
            $databases[] = $row[0];
            echo "  找到数据库: <span class='warning'>{$row[0]}</span>\n";
        }
    }
    echo "\n";

    // 2. 检查每个数据库中的表
    $dbInfo = [];
    foreach ($databases as $dbName) {
        echo "<span class='info'>[2] 检查数据库 '{$dbName}' 的表...</span>\n";

        $mysqli->select_db($dbName);

        // 检查 edge 表
        $result = $mysqli->query("SHOW TABLES LIKE 'edge%'");
        $edgeTables = $result ? $result->num_rows : 0;

        // 检查 SK 表
        $result = $mysqli->query("SHOW TABLES LIKE 'SK%'");
        $skTables = $result ? $result->num_rows : 0;

        // 总表数
        $result = $mysqli->query("SHOW TABLES");
        $totalTables = $result ? $result->num_rows : 0;

        $dbInfo[$dbName] = [
            'edge' => $edgeTables,
            'SK' => $skTables,
            'total' => $totalTables
        ];

        echo "  edge 前缀表: <span class='warning'>{$edgeTables}</span>\n";
        echo "  SK 前缀表: <span class='success'>{$skTables}</span>\n";
        echo "  总表数: {$totalTables}\n\n";
    }

    // 3. 诊断问题
    echo "<span class='info'>[3] 问题诊断...</span>\n";

    $problems = [];
    $correctDb = null;

    // 找到有最多表的数据库
    foreach ($dbInfo as $dbName => $info) {
        if ($info['total'] > 0) {
            if ($correctDb === null || $info['total'] > $dbInfo[$correctDb]['total']) {
                $correctDb = $dbName;
            }
        }
    }

    if ($correctDb) {
        echo "  <span class='success'>✓ 主数据库应该是: {$correctDb}</span>\n";

        if ($dbInfo[$correctDb]['edge'] > 0) {
            $problems[] = "数据库 '{$correctDb}' 中有 {$dbInfo[$correctDb]['edge']} 个 edge 前缀的表需要重命名为 SK";
            echo "  <span class='error'>✗ 发现 edge 前缀的表</span>\n";
        }

        if ($dbInfo[$correctDb]['SK'] === 0) {
            $problems[] = "数据库 '{$correctDb}' 中没有 SK 前缀的表";
            echo "  <span class='error'>✗ 没有 SK 前缀的表</span>\n";
        }
    } else {
        $problems[] = "没有找到包含表的数据库";
        echo "  <span class='error'>✗ 所有数据库都是空的</span>\n";
    }

    echo "\n";

    // 4. 自动修复
    if ($action === 'fix' && !empty($problems)) {
        echo "<span class='info'>[4] 开始自动修复...</span>\n\n";

        if ($correctDb) {
            $mysqli->select_db($correctDb);

            // 重命名 edge 表为 SK 表
            if ($dbInfo[$correctDb]['edge'] > 0) {
                echo "<span class='info'>重命名 edge 表为 SK 表...</span>\n";

                $result = $mysqli->query("SHOW TABLES LIKE 'edge%'");
                $tables = [];

                if ($result) {
                    while ($row = $result->fetch_array()) {
                        $tables[] = $row[0];
                    }
                }

                $success = 0;
                $failed = 0;
                $total = count($tables);

                foreach ($tables as $i => $table) {
                    $newTable = preg_replace('/^edge/', 'SK', $table);
                    $num = $i + 1;

                    echo "  [{$num}/{$total}] {$table} → {$newTable} ... ";

                    $sql = "RENAME TABLE `{$table}` TO `{$newTable}`";
                    if ($mysqli->query($sql)) {
                        echo "<span class='success'>✓</span>\n";
                        $success++;
                    } else {
                        echo "<span class='error'>✗ {$mysqli->error}</span>\n";
                        $failed++;
                    }
                }

                echo "\n";
                echo "<span class='success'>✓ 重命名完成: {$success} 成功, {$failed} 失败</span>\n";
            }
        }

        echo "\n<span class='info'>========================================</span>\n";
        echo "<span class='success'>✓ 修复完成！</span>\n";
        echo "<span class='info'>========================================</span>\n";
    } elseif ($action === 'fix') {
        echo "<span class='success'>✓ 没有发现需要修复的问题</span>\n";
    }

    // 5. 显示建议
    echo "\n<span class='info'>建议:</span>\n";
    if ($correctDb === 'edge') {
        echo "  <span class='success'>✓ 数据库名称正确 (edge)</span>\n";
    } elseif ($correctDb === 'edges') {
        echo "  <span class='warning'>! 数据库名称是 'edges'，建议在安装向导中填写 'edges'</span>\n";
    }

    if (!empty($problems) && $action !== 'fix') {
        echo "  <span class='warning'>! 点击 '自动修复' 按钮解决以上问题</span>\n";
    }

    $mysqli->close();
} else {
    echo "<span class='info'>点击 '诊断问题' 开始检查数据库...</span>\n";
}
?>
    </div>
</div>
</body>
</html>
