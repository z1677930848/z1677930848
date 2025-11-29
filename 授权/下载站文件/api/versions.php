<?php
/**
 * 版本管理 API
 */

require_once __DIR__ . '/../config.php';
require_once __DIR__ . '/../database.php';
require_once __DIR__ . '/../utils.php';
require_once __DIR__ . '/response.php';
require_once __DIR__ . '/license_guard.php';
require_once __DIR__ . '/admin/auth_guard.php';

header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization, X-Admin-Token');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

class VersionAPI {
    private $pdo;

    public function __construct($pdo = null) {
        if ($pdo === null) {
            $pdo = Database::getInstance()->getConnection();
        }
        $this->pdo = $pdo;
    }

    /**
     * 平台限制：仅支持 Linux
     */
    private function ensureLinux($os) {
        if ($os !== 'linux') {
            Response::error('当前仅支持 Linux 平台', 400);
        }
    }

    /**
     * 使用 version_compare 选出最高版本
     */
    private function pickLatestByVersion(array $rows) {
        if (empty($rows)) {
            return null;
        }
        $latest = $rows[0];
        foreach ($rows as $row) {
            $cmp = version_compare($row['version'], $latest['version']);
            if ($cmp > 0) {
                $latest = $row;
            } elseif ($cmp === 0) {
                if ($row['release_time'] > $latest['release_time'] || ($row['release_time'] === $latest['release_time'] && $row['id'] > $latest['id'])) {
                    $latest = $row;
                }
            }
        }
        return $latest;
    }

    /**
     * 读取升级策略（system_config.update_policy_{component}，JSON）
     * 示例：
     * {
     *   "target_version": "1.1.0",
     *   "min_compatible_version": "1.0.0",
     *   "rollout_percent": 100,
     *   "canary_licenses": ["ABC-123"],
     *   "canary_domains": ["example.com"]
     * }
     */
    private function loadPolicy(string $component): array {
        $default = [
            'target_version' => null,
            'min_compatible_version' => null,
            'rollout_percent' => 100,
            'canary_licenses' => [],
            'canary_domains' => [],
        ];

        $key = 'update_policy_' . $component;
        $stmt = $this->pdo->prepare("SELECT config_value FROM system_config WHERE config_key = ?");
        $stmt->execute([$key]);
        $row = $stmt->fetch();
        if (!$row || empty($row['config_value'])) {
            return $default;
        }

        $json = json_decode($row['config_value'], true);
        if (!is_array($json)) {
            return $default;
        }

        return array_merge($default, $json);
    }

    /**
     * 灰度/白名单判断：优先白名单，其次按 rollout_percent 分桶
     */
    private function inRollout(array $policy, string $licenseCode = '', string $domain = ''): bool {
        $rollout = isset($policy['rollout_percent']) ? intval($policy['rollout_percent']) : 100;
        if ($rollout >= 100) {
            return true;
        }
        if (!empty($policy['canary_licenses']) && in_array($licenseCode, (array)$policy['canary_licenses'], true)) {
            return true;
        }
        if (!empty($policy['canary_domains']) && in_array($domain, (array)$policy['canary_domains'], true)) {
            return true;
        }

        $seed = $licenseCode ?: $domain;
        if ($seed === '') {
            $seed = uniqid('rollout', true);
        }
        $bucket = crc32($seed) % 100;
        return $bucket < $rollout;
    }

    /**
     * 获取最新版
     */
    public function getLatestVersion($component, $os = 'linux', $arch = 'amd64') {
        try {
            $this->ensureLinux($os);

            $stmt = $this->pdo->prepare("
                SELECT * FROM versions
                WHERE component_code = ? AND os = ? AND arch = ? AND status = 1
                ORDER BY release_time DESC, id DESC
            ");
            $stmt->execute([$component, $os, $arch]);
            $rows = $stmt->fetchAll();
            $version = $this->pickLatestByVersion($rows);

            if (!$version) {
                Response::error('未找到版本信息', 404);
            }

            Response::success([
                'id' => $version['id'],
                'component' => $version['component_name'],
                'version' => $version['version'],
                'os' => $version['os'],
                'arch' => $version['arch'],
                'filename' => $version['filename'],
                'file_size' => $version['file_size'],
                'file_md5' => $version['file_md5'],
                'download_url' => $this->getDownloadUrl($version),
                'release_time' => $version['release_time'],
                'description' => $version['description'],
                'changelog' => $version['changelog']
            ]);
        } catch (PDOException $e) {
            Response::error('数据库错误: ' . $e->getMessage());
        }
    }

    /**
     * 获取版本列表
     */
    public function getVersionList($component = null, $os = 'linux', $arch = null, $page = 1, $limit = 20) {
        try {
            $this->ensureLinux($os);
            $where = ['status = 1', 'os = ?'];
            $params = [$os];

            if ($component) {
                $where[] = 'component_code = ?';
                $params[] = $component;
            }
            if ($arch) {
                $where[] = 'arch = ?';
                $params[] = $arch;
            }

            $whereClause = implode(' AND ', $where);

            $countStmt = $this->pdo->prepare("SELECT COUNT(*) as total FROM versions WHERE {$whereClause}");
            $countStmt->execute($params);
            $total = $countStmt->fetch()['total'];

            $offset = ($page - 1) * $limit;
            $stmt = $this->pdo->prepare("
                SELECT * FROM versions
                WHERE {$whereClause}
                ORDER BY release_time DESC, id DESC
                LIMIT ? OFFSET ?
            ");
            $paramsWithLimit = array_merge($params, [$limit, $offset]);
            $stmt->execute($paramsWithLimit);
            $versions = $stmt->fetchAll();

            $list = array_map(function($v) {
                return [
                    'id' => $v['id'],
                    'component' => $v['component_name'],
                    'component_code' => $v['component_code'],
                    'version' => $v['version'],
                    'os' => $v['os'],
                    'arch' => $v['arch'],
                    'filename' => $v['filename'],
                    'file_size' => $v['file_size'],
                    'file_md5' => $v['file_md5'],
                    'download_url' => $this->getDownloadUrl($v),
                    'release_time' => $v['release_time'],
                    'download_count' => $v['download_count']
                ];
            }, $versions);

            Response::success([
                'list' => $list,
                'total' => $total,
                'page' => $page,
                'limit' => $limit,
                'pages' => ceil($total / $limit)
            ]);
        } catch (PDOException $e) {
            Response::error('数据库错误: ' . $e->getMessage());
        }
    }

    /**
     * 检查更新
     */
    public function checkUpdate($component, $currentVersion, $os = 'linux', $arch = 'amd64') {
        try {
            $this->ensureLinux($os);

            $stmt = $this->pdo->prepare("
                SELECT * FROM versions
                WHERE component_code = ? AND os = ? AND arch = ? AND status = 1
                ORDER BY release_time DESC, id DESC
            ");
            $stmt->execute([$component, $os, $arch]);
            $rows = $stmt->fetchAll();
            $latest = $this->pickLatestByVersion($rows);

            if (!$latest) {
                Response::error('未找到版本信息', 404);
            }

            // 读取升级策略（可选配置）
            $policy = $this->loadPolicy($component);
            $credentials = LicenseGuard::getCredentialsFromRequest();
            $targetVersion = $policy['target_version'] ?: $latest['version'];
            $minCompatible = $policy['min_compatible_version'];

            $isRequired = false;
            if (!empty($minCompatible) && version_compare($currentVersion, $minCompatible, '<')) {
                $isRequired = true;
            }

            $isRecommended = version_compare($currentVersion, $targetVersion, '<') && $this->inRollout(
                $policy,
                $credentials['license_code'] ?? '',
                $credentials['domain'] ?? ''
            );

            $hasUpdate = $isRequired || $isRecommended;

            Response::success([
                'has_update' => $hasUpdate,
                'current_version' => $currentVersion,
                'latest_version' => $latest['version'],
                'target_version' => $targetVersion,
                'min_compatible_version' => $minCompatible,
                'is_required' => $isRequired || (bool)$latest['is_required_update'],
                'is_recommended' => $isRecommended,
                'download_url' => $hasUpdate ? $this->getDownloadUrl($latest) : null,
                'file_sha256' => $latest['file_sha256'] ?? null,
                'file_md5' => $latest['file_md5'] ?? null,
                'file_size' => $latest['file_size'] ?? null,
                'changelog' => $hasUpdate ? $latest['changelog'] : null,
                'policy' => $policy
            ]);
        } catch (PDOException $e) {
            Response::error('数据库错误: ' . $e->getMessage());
        }
    }

    /**
     * 获取下载URL
     */
    private function getDownloadUrl($version) {
        $baseUrl = rtrim($this->getDownloadHost(), '/');
        $path = "updates/{$version['component_code']}/{$version['os']}/{$version['arch']}/{$version['filename']}";
        return $baseUrl . '/' . $path;
    }

    /**
     * 获取下载服务器地址
     */
    private function getDownloadHost() {
        static $host = null;
        if ($host === null) {
            $stmt = $this->pdo->prepare("SELECT config_value FROM system_config WHERE config_key = 'download_host'");
            $stmt->execute();
            $config = $stmt->fetch();
            // 若未在 system_config.download_host 配置，则回退到正式域名
            $host = $config ? $config['config_value'] : 'https://dl.lingcdn.cloud';
        }
        return $host;
    }

    /**
     * 获取版本详情
     */
    public function getVersionDetail($id) {
        try {
            $stmt = $this->pdo->prepare("SELECT * FROM versions WHERE id = ?");
            $stmt->execute([$id]);
            $version = $stmt->fetch();

            if (!$version) {
                Response::error('版本不存在', 404);
            }
            $this->ensureLinux($version['os']);

            Response::success([
                'id' => $version['id'],
                'component' => $version['component_name'],
                'component_code' => $version['component_code'],
                'version' => $version['version'],
                'os' => $version['os'],
                'arch' => $version['arch'],
                'filename' => $version['filename'],
                'file_size' => $version['file_size'],
                'file_md5' => $version['file_md5'],
                'download_url' => $this->getDownloadUrl($version),
                'release_time' => $version['release_time'],
                'description' => $version['description'],
                'changelog' => $version['changelog'],
                'download_count' => $version['download_count'],
                'is_required_update' => (bool)$version['is_required_update']
            ]);
        } catch (PDOException $e) {
            Response::error('数据库错误: ' . $e->getMessage());
        }
    }

    /**
     * 删除版本
     */
    public function deleteVersion($id) {
        try {
            $stmt = $this->pdo->prepare("SELECT * FROM versions WHERE id = ?");
            $stmt->execute([$id]);
            $version = $stmt->fetch();

            if (!$version) {
                Response::error('版本不存在', 404);
            }
            $this->ensureLinux($version['os']);

            $filePath = UPDATE_FILES_PATH .
                        $version['component_code'] . '/' .
                        $version['os'] . '/' .
                        $version['arch'] . '/' .
                        $version['filename'];

            if (file_exists($filePath) && !unlink($filePath)) {
                Response::error('删除文件失败');
            }

            $deleteStmt = $this->pdo->prepare("DELETE FROM versions WHERE id = ?");
            if (!$deleteStmt->execute([$id])) {
                Response::error('删除数据库记录失败');
            }

            Response::success([
                'id' => $id,
                'filename' => $version['filename'],
                'message' => '删除成功'
            ]);
        } catch (PDOException $e) {
            Response::error('数据库错误: ' . $e->getMessage());
        } catch (Exception $e) {
            Response::error('删除失败: ' . $e->getMessage());
        }
    }
}

// 处理请求
$action = $_GET['action'] ?? 'list';
$pdo = Database::getInstance()->getConnection();

$licenseProtected = ['latest', 'list', 'check', 'detail'];
$adminProtected = ['delete'];

if (in_array($action, $licenseProtected, true)) {
    $credentials = LicenseGuard::getCredentialsFromRequest();
    LicenseGuard::verify(
        $pdo,
        $credentials['license_code'],
        $credentials['system_token'],
        $credentials['domain'],
        'versions_' . $action
    );
}

if (in_array($action, $adminProtected, true)) {
    AdminAuth::requireAuth();
}
$api = new VersionAPI($pdo);

switch ($action) {
    case 'latest':
        $component = $_GET['component'] ?? '';
        $os = $_GET['os'] ?? 'linux';
        $arch = $_GET['arch'] ?? 'amd64';
        $api->getLatestVersion($component, $os, $arch);
        break;

    case 'list':
        $component = $_GET['component'] ?? null;
        $os = $_GET['os'] ?? 'linux';
        $arch = $_GET['arch'] ?? null;
        $page = intval($_GET['page'] ?? 1);
        $limit = intval($_GET['limit'] ?? 20);
        $api->getVersionList($component, $os, $arch, $page, $limit);
        break;

    case 'check':
        $component = $_GET['component'] ?? '';
        $currentVersion = $_GET['version'] ?? '0.0.0';
        $os = $_GET['os'] ?? 'linux';
        $arch = $_GET['arch'] ?? 'amd64';
        $api->checkUpdate($component, $currentVersion, $os, $arch);
        break;

    case 'detail':
        $id = intval($_GET['id'] ?? 0);
        if ($id <= 0) {
            Response::error('无效的版本ID', 400);
        }
        $api->getVersionDetail($id);
        break;

    case 'delete':
        $id = intval($_GET['id'] ?? 0);
        if ($id <= 0) {
            Response::error('无效的版本ID', 400);
        }
        $api->deleteVersion($id);
        break;

    default:
        Response::error('无效的操作', 400);
}
