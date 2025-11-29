<?php
/**
 * 鎵弿鐗堟湰鏂囦欢骞惰嚜鍔ㄥ悓姝ュ埌鏁版嵁搴? */

require_once __DIR__ . '/config.php';
require_once __DIR__ . '/database.php';

header('Content-Type: application/json; charset=utf-8');

/**
 * 纭繚鐩綍瀛樺湪涓斿彲鍐欙紝涓嶅彲鍐欑洿鎺ユ姏閿欎究浜庡畾浣嶆潈闄愰棶棰? */
function ensure_dir_writable($dir, $desc) {
    if (!is_dir($dir)) {
        if (!@mkdir($dir, 0755, true)) {
            throw new RuntimeException($desc . ' 鍒涘缓澶辫触锛屾鏌ユ潈闄? ' . $dir);
        }
    }
    if (!is_writable($dir)) {
        throw new RuntimeException($desc . ' 涓嶅彲鍐欙紝妫€鏌ユ潈闄? ' . $dir);
    }
}

class VersionScanner {
    private $pdo;
    private $stats = [
        'scanned' => 0,
        'added' => 0,
        'updated' => 0,
        'skipped' => 0,
        'errors' => []
    ];

    public function __construct() {
        $this->pdo = Database::getInstance()->getConnection();
        // 核验更新目录权限，避免后续频繁报 Permission denied
        ensure_dir_writable(UPDATE_FILES_PATH, '更新目录');
        ensure_dir_writable(LOG_PATH, '日志目录');
    }

    /**
     * 扫描所有版本文件
     */
    public function scan() {
        $basePath = UPDATE_FILES_PATH;
        $components = ['admin', 'api', 'node'];

        foreach ($components as $component) {
            $componentPath = $basePath . $component;
            if (!is_dir($componentPath)) {
                continue;
            }

            $this->scanComponent($component, $componentPath);
        }

        // 扫描完成后，只启用每个组件/平台组合的最新版本
        $this->enableLatestVersionsOnly();

        return $this->stats;
    }

    /**
     * 只启用每个组件+OS+架构的最新版本，禁用其他版本
     */
    private function enableLatestVersionsOnly() {
        $components = ['admin', 'api', 'node'];

        foreach ($components as $component) {
            $comboStmt = $this->pdo->prepare("\n                SELECT DISTINCT os, arch\n                FROM versions\n                WHERE component_code = ?\n            ");
            $comboStmt->execute([$component]);
            $combos = $comboStmt->fetchAll();

            foreach ($combos as $combo) {
                $os = $combo['os'];
                $arch = $combo['arch'];

                $disableStmt = $this->pdo->prepare("\n                    UPDATE versions SET status = 0\n                    WHERE component_code = ? AND os = ? AND arch = ?\n                ");
                $disableStmt->execute([$component, $os, $arch]);

                $latestStmt = $this->pdo->prepare("\n                    SELECT id FROM versions\n                    WHERE component_code = ? AND os = ? AND arch = ?\n                    ORDER BY release_time DESC, id DESC\n                    LIMIT 1\n                ");
                $latestStmt->execute([$component, $os, $arch]);
                $latest = $latestStmt->fetch();

                if ($latest) {
                    $enableStmt = $this->pdo->prepare("UPDATE versions SET status = 1 WHERE id = ?");
                    $enableStmt->execute([$latest['id']]);
                }
            }
        }
    }
    private function scanComponent($component, $path) {
        $systems = ['linux', 'windows', 'darwin'];

        foreach ($systems as $os) {
            $osPath = $path . '/' . $os;
            if (!is_dir($osPath)) continue;

            $archs = ['amd64', 'arm64', '386'];
            foreach ($archs as $arch) {
                $archPath = $osPath . '/' . $arch;
                if (!is_dir($archPath)) continue;

                $this->scanFiles($component, $os, $arch, $archPath);
            }
        }
    }

    /**
     * 鎵弿鏂囦欢
     */
    private function scanFiles($component, $os, $arch, $path) {
        $files = glob($path . '/*');

        foreach ($files as $file) {
            if (!is_file($file)) continue;

            $filename = basename($file);

            // 璺宠繃鏍￠獙鏂囦欢
            if (preg_match('/\.(md5|sha256)$/', $filename)) continue;

            $this->stats['scanned']++;

            // 瑙ｆ瀽鐗堟湰鍙?            $version = $this->parseVersion($filename);
            if (!$version) {
                $this->stats['skipped']++;
                continue;
            }

            // 鑾峰彇鏂囦欢淇℃伅
            $fileSize = filesize($file);
            $md5File = $file . '.md5';
            $sha256File = $file . '.sha256';

            // 璇诲彇MD5鍜孲HA256锛屽彧鍙栧搱甯屽€奸儴鍒?            if (file_exists($md5File)) {
                $md5Content = trim(file_get_contents($md5File));
                $md5 = preg_match('/^([a-f0-9]{32})/i', $md5Content, $m) ? $m[1] : md5_file($file);
            } else {
                $md5 = md5_file($file);
            }

            if (file_exists($sha256File)) {
                $sha256Content = trim(file_get_contents($sha256File));
                $sha256 = preg_match('/^([a-f0-9]{64})/i', $sha256Content, $m) ? $m[1] : null;
            } else {
                $sha256 = null;
            }

            // 妫€鏌ユ暟鎹簱涓槸鍚﹀瓨鍦?            $existing = $this->findVersion($component, $version, $os, $arch);

            if ($existing) {
                // 鏇存柊鏂囦欢淇℃伅
                if ($existing['file_md5'] !== $md5 || $existing['file_size'] != $fileSize) {
                    $this->updateVersion($existing['id'], $filename, $fileSize, $md5, $sha256);
                    $this->stats['updated']++;
                } else {
                    $this->stats['skipped']++;
                }
                // 姣忔閮介噸鏂扮敓鎴恈hangelog
                $this->createChangelog($component, $version);
            } else {
                // 鏂板鐗堟湰
                $this->addVersion($component, $version, $os, $arch, $filename, $fileSize, $md5, $sha256);
                $this->createChangelog($component, $version);
                $this->stats['added']++;
            }
        }
    }

    /**
     * 瑙ｆ瀽鐗堟湰鍙?     */
    private function parseVersion($filename) {
        // 鍖归厤鏍煎紡: ling-admin-v1.0.17-linux-amd64
        if (preg_match('/v?(\d+\.\d+\.\d+)/', $filename, $matches)) {
            return $matches[1];
        }
        return null;
    }

    /**
     * 鏌ユ壘鐗堟湰
     */
    private function findVersion($component, $version, $os, $arch) {
        $stmt = $this->pdo->prepare("
            SELECT * FROM versions
            WHERE component_code = ? AND version = ? AND os = ? AND arch = ?
        ");
        $stmt->execute([$component, $version, $os, $arch]);
        return $stmt->fetch();
    }

    /**
     * 娣诲姞鐗堟湰
     */
    private function addVersion($component, $version, $os, $arch, $filename, $fileSize, $md5, $sha256) {
        $componentNames = [
            'admin' => 'Ling Admin',
            'api' => 'Ling API',
            'node' => 'Ling Node'
        ];

        $stmt = $this->pdo->prepare("
            INSERT INTO versions
            (component_code, component_name, version, os, arch, filename, file_size, file_md5, file_sha256, release_time, status)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), 0)
        ");

        $stmt->execute([
            $component,
            $componentNames[$component] ?? ucfirst($component),
            $version,
            $os,
            $arch,
            $filename,
            $fileSize,
            $md5,
            $sha256
        ]);
    }

    /**
     * 鏇存柊鐗堟湰
     */
    private function updateVersion($id, $filename, $fileSize, $md5, $sha256) {
        $stmt = $this->pdo->prepare("
            UPDATE versions
            SET filename = ?, file_size = ?, file_md5 = ?, file_sha256 = ?, updated_at = NOW()
            WHERE id = ?
        ");

        $stmt->execute([$filename, $fileSize, $md5, $sha256, $id]);
    }

    /**
     * 鍒涘缓鏇存柊鏃ュ織鏂囦欢
     */
    private function createChangelog($component, $version) {
        $changelogDir = __DIR__ . '/changelog';
        ensure_dir_writable($changelogDir, '鏇存柊鏃ュ織鐩綍');

        $changelogFile = $changelogDir . '/v' . $version . '.php';
        // 绉婚櫎妫€鏌ワ紝姣忔閮介噸鏂扮敓鎴?        // if (file_exists($changelogFile)) {
        //     return;
        // }

        $php = <<<'PHP'
<?php
require_once __DIR__ . '/../config.php';
require_once __DIR__ . '/../database.php';

$version = 'VERSION_PLACEHOLDER';
$db = Database::getInstance()->getConnection();
$stmt = $db->prepare("SELECT * FROM versions WHERE version = ? AND component_code = 'admin' LIMIT 1");
$stmt->execute([$version]);
$data = $stmt->fetch();

$changelog = $data['changelog'] ?? '鏆傛棤鏇存柊璇存槑';
$description = $data['description'] ?? '';
$releaseTime = $data['release_time'] ?? date('Y-m-d');
?>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LingCDN v<?= htmlspecialchars($version) ?> 鏇存柊鏃ュ織</title>
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
    <p class="version">鍙戝竷鏃ユ湡: <?= htmlspecialchars(substr($releaseTime, 0, 10)) ?></p>

    <?php if ($description): ?>
    <div class="description"><?= nl2br(htmlspecialchars($description)) ?></div>
    <?php endif; ?>

    <div class="changelog"><?= nl2br(htmlspecialchars($changelog)) ?></div>
</body>
</html>
PHP;

        $php = str_replace('VERSION_PLACEHOLDER', $version, $php);
        file_put_contents($changelogFile, $php);
        chmod($changelogFile, 0644);

        // 鐢熸垚闈欐€?HTML锛屽吋瀹规棫鐗堟湰锛堟瘡娆￠兘閲嶆柊鐢熸垚浠ュ悓姝ユ渶鏂板唴瀹癸級
        $htmlFile = $changelogDir . '/v' . $version . '.html';
        $htmlContent = '';
        $renderError = null;
        ob_start();
        try {
            include $changelogFile;
            $htmlContent = ob_get_clean();
        } catch (Throwable $e) {
            ob_end_clean();
            $renderError = $e->getMessage();
        }
        if ($renderError !== null) {
            throw new RuntimeException('鐢熸垚鏇存柊鏃ュ織HTML澶辫触: ' . $renderError);
        }

        if (file_put_contents($htmlFile, $htmlContent) === false) {
            throw new RuntimeException('鏇存柊鏃ュ織鏂囦欢鍐欏叆澶辫触锛屽彲鑳芥棤鏉冮檺, path: ' . $htmlFile);
        }
        @chmod($htmlFile, 0644);
    }
}

// 澶勭悊璇锋眰
if (isset($_GET['scan'])) {
    try {
        $scanner = new VersionScanner();
        $stats = $scanner->scan();

        echo json_encode([
            'code' => 200,
            'message' => 'success',
            'data' => $stats
        ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
    } catch (Exception $e) {
        echo json_encode([
            'code' => 500,
            'message' => '鎵弿澶辫触: ' . $e->getMessage(),
            'data' => null
        ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
    }
} else {
    echo json_encode([
        'code' => 400,
        'message' => '璇蜂娇鐢??scan=1 鍙傛暟',
        'data' => null
    ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
}







