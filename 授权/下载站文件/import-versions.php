<?php
/**
 * 版本数据导入脚本
 * 扫描 updates 目录，将现有版本文件导入到数据库
 */

require_once 'config.php';
require_once 'database.php';
require_once 'utils.php';

echo "==========================================\n";
echo "开始导入版本数据...\n";
echo "==========================================\n\n";

try {
    $db = Database::getInstance()->getConnection();

    // 扫描 updates 目录
    $basePath = UPDATE_FILES_PATH;
    $components = ['admin', 'api', 'node'];
    $imported = 0;
    $skipped = 0;

    foreach ($components as $component) {
        $componentPath = $basePath . $component . '/';

        if (!is_dir($componentPath)) {
            echo "⚠️  组件目录不存在: {$component}\n";
            continue;
        }

        // 扫描 linux/amd64 目录
        $versionPath = $componentPath . 'linux/amd64/';

        if (!is_dir($versionPath)) {
            echo "⚠️  版本目录不存在: {$versionPath}\n";
            continue;
        }

        echo "📁 扫描目录: {$versionPath}\n";

        $files = glob($versionPath . '*.{zip,tar.gz,tgz}', GLOB_BRACE);

        foreach ($files as $filePath) {
            $filename = basename($filePath);

            // 解析文件名提取版本号
            // 格式: sk-admin-v1.0.5-linux-amd64.zip 或 edge-admin-v1.0.0-linux-amd64.zip
            if (preg_match('/(sk|edge)-(admin|api|node)-v([\d.]+)-(linux|windows|darwin)-(amd64|arm64|386)/i', $filename, $matches)) {
                $componentCode = $matches[2];
                $version = $matches[3];
                $os = $matches[4];
                $arch = $matches[5];

                // 检查是否已存在
                $checkSql = "SELECT id FROM versions
                            WHERE component_code = :component
                            AND version = :version
                            AND os = :os
                            AND arch = :arch";
                $checkStmt = $db->prepare($checkSql);
                $checkStmt->execute([
                    'component' => $componentCode,
                    'version' => $version,
                    'os' => $os,
                    'arch' => $arch
                ]);

                if ($checkStmt->fetch()) {
                    echo "  ⏭️  跳过（已存在）: {$filename}\n";
                    $skipped++;
                    continue;
                }

                // 获取文件信息
                $fileSize = filesize($filePath);
                $md5 = md5_file($filePath);
                $sha256 = hash_file('sha256', $filePath);

                // 组件名称映射
                $componentNames = [
                    'admin' => 'LingCDN Admin',
                    'api' => 'LingCDN API',
                    'node' => 'LingCDN Node'
                ];

                $componentName = $componentNames[$componentCode] ?? ucfirst($componentCode);

                // 插入数据库
                $sql = "INSERT INTO versions (
                    component_code, component_name, version, os, arch,
                    filename, file_size, file_md5, file_sha256,
                    description, release_time, status
                ) VALUES (
                    :component_code, :component_name, :version, :os, :arch,
                    :filename, :file_size, :file_md5, :file_sha256,
                    :description, :release_time, :status
                )";

                $stmt = $db->prepare($sql);
                $result = $stmt->execute([
                    'component_code' => $componentCode,
                    'component_name' => $componentName,
                    'version' => $version,
                    'os' => $os,
                    'arch' => $arch,
                    'filename' => $filename,
                    'file_size' => $fileSize,
                    'file_md5' => $md5,
                    'file_sha256' => $sha256,
                    'description' => "从现有文件导入 - {$componentName} v{$version}",
                    'release_time' => date('Y-m-d H:i:s', filemtime($filePath)),
                    'status' => 1
                ]);

                if ($result) {
                    echo "  ✅ 导入成功: {$filename} (v{$version}, " . Utils::formatBytes($fileSize) . ")\n";
                    $imported++;
                } else {
                    echo "  ❌ 导入失败: {$filename}\n";
                }

            } else {
                echo "  ⚠️  无法解析文件名: {$filename}\n";
            }
        }

        echo "\n";
    }

    echo "==========================================\n";
    echo "导入完成！\n";
    echo "成功导入: {$imported} 个版本\n";
    echo "跳过: {$skipped} 个版本\n";
    echo "==========================================\n\n";

    // 显示导入的版本列表
    $sql = "SELECT component_code, component_name, version, os, arch, file_size
            FROM versions
            ORDER BY component_code, version DESC";
    $stmt = $db->query($sql);
    $versions = $stmt->fetchAll();

    echo "📋 当前数据库中的版本列表：\n\n";
    $currentComponent = '';
    foreach ($versions as $v) {
        if ($currentComponent !== $v['component_code']) {
            $currentComponent = $v['component_code'];
            echo "\n🔹 {$v['component_name']}:\n";
        }
        echo "   v{$v['version']} ({$v['os']}/{$v['arch']}) - " . Utils::formatBytes($v['file_size']) . "\n";
    }

    echo "\n";
    echo "✅ 导入完成！现在可以通过 API 获取版本信息了。\n\n";

} catch (Exception $e) {
    echo "❌ 错误: " . $e->getMessage() . "\n";
    echo "详细信息: " . $e->getTraceAsString() . "\n";
    exit(1);
}
