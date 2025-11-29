-- =============================================
-- 授权和更新系统数据库结构
-- =============================================

-- 创建数据�?
CREATE DATABASE IF NOT EXISTS `license_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `license_system`;

-- =============================================
-- 1. 许可证表
-- =============================================
CREATE TABLE IF NOT EXISTS `licenses` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `license_code` VARCHAR(100) NOT NULL UNIQUE COMMENT '许可证密�?,
  `license_type` VARCHAR(50) NOT NULL DEFAULT 'standard' COMMENT '许可证类�? trial(试用), standard(标准), professional(专业), enterprise(企业)',
  `customer_name` VARCHAR(200) NOT NULL COMMENT '客户名称',
  `customer_email` VARCHAR(200) DEFAULT NULL COMMENT '客户邮箱',
  `customer_phone` VARCHAR(50) DEFAULT NULL COMMENT '客户电话',
  `company_name` VARCHAR(200) DEFAULT NULL COMMENT '公司名称',

  -- 授权控制
  `allowed_domains` JSON DEFAULT NULL COMMENT '允许的域名列表，JSON数组',
  `bound_system_tokens` JSON DEFAULT NULL COMMENT '已绑定的系统令牌（机器码），JSON数组',
  `max_devices` INT DEFAULT 1 COMMENT '最大授权设备数',

  -- 功能控制
  `features` JSON DEFAULT NULL COMMENT '授权功能列表，JSON对象',

  -- 时间控制
  `issue_time` DATETIME NOT NULL COMMENT '发放时间',
  `expire_time` DATETIME DEFAULT NULL COMMENT '过期时间，NULL表示永久',
  `last_check_time` DATETIME DEFAULT NULL COMMENT '最后验证时�?,

  -- 状态与统计
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状�? 0-禁用, 1-启用',
  `check_count` INT UNSIGNED DEFAULT 0 COMMENT '验证次数',

  -- 备注
  `remark` TEXT DEFAULT NULL COMMENT '备注',

  -- 时间�?
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX `idx_license_code` (`license_code`),
  INDEX `idx_customer_name` (`customer_name`),
  INDEX `idx_status` (`status`),
  INDEX `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='许可证表';

-- =============================================
-- 2. 授权验证日志�?
-- =============================================
CREATE TABLE IF NOT EXISTS `authorization_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `license_id` INT UNSIGNED NOT NULL COMMENT '许可证ID',
  `system_token` VARCHAR(100) NOT NULL COMMENT '系统令牌',
  `domain` VARCHAR(200) NOT NULL COMMENT '请求域名',
  `ip_address` VARCHAR(50) NOT NULL COMMENT 'IP地址',
  `status` VARCHAR(20) NOT NULL COMMENT '状�? success, failed',
  `error_message` VARCHAR(500) DEFAULT NULL COMMENT '错误信息',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX `idx_license_id` (`license_id`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='授权验证日志�?;

-- =============================================
-- 3. 版本�?
-- =============================================
CREATE TABLE IF NOT EXISTS `versions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `component_code` VARCHAR(50) NOT NULL COMMENT '组件代码: admin, node, api',
  `component_name` VARCHAR(100) NOT NULL COMMENT '组件名称',
  `version` VARCHAR(50) NOT NULL COMMENT '版本�?,
  `os` VARCHAR(20) NOT NULL COMMENT '操作系统: linux, windows, darwin, all',
  `arch` VARCHAR(20) NOT NULL COMMENT '架构: amd64, arm64, 386, all',

  -- 文件信息
  `filename` VARCHAR(200) NOT NULL COMMENT '文件�?,
  `file_size` BIGINT UNSIGNED NOT NULL COMMENT '文件大小(字节)',
  `file_md5` VARCHAR(32) NOT NULL COMMENT 'MD5校验�?,
  `file_sha256` VARCHAR(64) DEFAULT NULL COMMENT 'SHA256校验�?,

  -- 版本信息
  `description` TEXT DEFAULT NULL COMMENT '版本描述',
  `changelog` TEXT DEFAULT NULL COMMENT '更新日志',
  `release_time` DATETIME NOT NULL COMMENT '发布时间',
  `is_required_update` TINYINT DEFAULT 0 COMMENT '是否强制更新: 0-�? 1-�?,

  -- 状态与统计
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状�? 0-禁用, 1-启用',
  `download_count` INT UNSIGNED DEFAULT 0 COMMENT '下载次数',

  -- 时间�?
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY `uk_component_version_os_arch` (`component_code`, `version`, `os`, `arch`),
  INDEX `idx_component_code` (`component_code`),
  INDEX `idx_version` (`version`),
  INDEX `idx_status` (`status`),
  INDEX `idx_release_time` (`release_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='版本�?;

-- =============================================
-- 4. 下载日志�?
-- =============================================
CREATE TABLE IF NOT EXISTS `download_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `version_id` INT UNSIGNED NOT NULL COMMENT '版本ID',
  `ip_address` VARCHAR(50) NOT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
  `download_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下载时间',

  INDEX `idx_version_id` (`version_id`),
  INDEX `idx_download_time` (`download_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='下载日志�?;

-- =============================================
-- 5. 系统配置�?
-- =============================================
CREATE TABLE IF NOT EXISTS `system_config` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `config_key` VARCHAR(100) NOT NULL UNIQUE COMMENT '配置�?,
  `config_value` TEXT NOT NULL COMMENT '配置�?,
  `config_type` VARCHAR(20) DEFAULT 'string' COMMENT '配置类型: string, number, json',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '描述',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX `idx_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置�?;

-- =============================================
-- 6. 管理员表
-- =============================================
CREATE TABLE IF NOT EXISTS `admins` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户�?,
  `password` VARCHAR(255) NOT NULL COMMENT '密码(加密)',
  `real_name` VARCHAR(100) DEFAULT NULL COMMENT '真实姓名',
  `email` VARCHAR(200) DEFAULT NULL COMMENT '邮箱',
  `role` VARCHAR(20) NOT NULL DEFAULT 'admin' COMMENT '角色: superadmin, admin',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状�? 0-禁用, 1-启用',
  `last_login_time` DATETIME DEFAULT NULL COMMENT '最后登录时�?,
  `last_login_ip` VARCHAR(50) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- =============================================
-- 7. 安装统计�?
-- =============================================
CREATE TABLE IF NOT EXISTS `install_stats` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `install_time` DATETIME NOT NULL COMMENT '安装时间',
  `ip_address` VARCHAR(50) NOT NULL COMMENT 'IP地址',
  `country` VARCHAR(100) DEFAULT NULL COMMENT '国家',
  `region` VARCHAR(100) DEFAULT NULL COMMENT '地区',
  `city` VARCHAR(100) DEFAULT NULL COMMENT '城市',
  `os` VARCHAR(20) NOT NULL COMMENT '操作系统: linux, windows, darwin',
  `arch` VARCHAR(20) NOT NULL COMMENT '架构: amd64, arm64, 386',
  `version` VARCHAR(50) NOT NULL COMMENT '安装版本',
  `install_type` VARCHAR(20) NOT NULL DEFAULT 'script' COMMENT '安装类型: script, manual',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX `idx_install_time` (`install_time`),
  INDEX `idx_os` (`os`),
  INDEX `idx_arch` (`arch`),
  INDEX `idx_version` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='安装统计�?;

-- =============================================
-- 插入初始数据
-- =============================================

-- 插入默认配置
INSERT INTO `system_config` (`config_key`, `config_value`, `config_type`, `description`) VALUES
('download_host', 'http://localhost', 'string', '下载服务器地址'),
('site_name', '授权管理系统', 'string', '站点名称'),
('allow_auto_bind', '1', 'number', '是否允许自动绑定新设�?);

-- 插入默认管理员账�?(用户�? admin, 密码: admin123)
-- 密码使用 password_hash 加密
INSERT INTO `admins` (`username`, `password`, `real_name`, `role`) VALUES
('admin', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '系统管理�?, 'superadmin');

-- =============================================
-- 示例数据
-- =============================================

-- 插入示例许可�?
INSERT INTO `licenses`
(`license_code`, `license_type`, `customer_name`, `customer_email`, `allowed_domains`, `max_devices`, `features`, `issue_time`, `expire_time`, `remark`)
VALUES
('ABCD-EFGH-IJKL-MNOP-QRST-UVWX', 'enterprise', '测试客户', 'test@example.com',
 '["*"]', 5,
 '{"cdn": true, "waf": true, "ssl": true, "max_bandwidth": "10Gbps"}',
 NOW(), DATE_ADD(NOW(), INTERVAL 1 YEAR), '企业版测试许可证');

-- 插入示例版本
INSERT INTO `versions`
(`component_code`, `component_name`, `version`, `os`, `arch`, `filename`, `file_size`, `file_md5`, `description`, `release_time`)
VALUES
('admin', 'Edge Admin', '1.0.0', 'linux', 'amd64', 'edge-admin-v1.0.0-linux-amd64.zip', 10485760, 'd41d8cd98f00b204e9800998ecf8427e', '初始版本发布', NOW()),
('node', 'Ling Node', '1.1.0', 'linux', 'amd64', 'ling-node-v1.1.0-linux-amd64.zip', 20971520, 'd41d8cd98f00b204e9800998ecf8427e', '边缘节点初始版本', NOW()),
('api', 'Edge API', '1.0.0', 'linux', 'amd64', 'edge-api-v1.0.0-linux-amd64.zip', 15728640, 'd41d8cd98f00b204e9800998ecf8427e', 'API服务初始版本', NOW());

