-- =============================================
-- 邮箱配置表
-- =============================================

USE `license_system`;

-- 邮箱配置表
CREATE TABLE IF NOT EXISTS `email_settings` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `smtp_host` VARCHAR(200) NOT NULL COMMENT 'SMTP服务器地址',
  `smtp_port` INT NOT NULL DEFAULT 587 COMMENT 'SMTP端口',
  `smtp_username` VARCHAR(200) NOT NULL COMMENT 'SMTP用户名',
  `smtp_password` VARCHAR(500) NOT NULL COMMENT 'SMTP密码(加密存储)',
  `from_email` VARCHAR(200) NOT NULL COMMENT '发件人邮箱',
  `from_name` VARCHAR(200) DEFAULT 'LingCDN' COMMENT '发件人名称',
  `use_tls` TINYINT NOT NULL DEFAULT 1 COMMENT '是否使用TLS: 0-否, 1-是',
  `is_enabled` TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用: 0-否, 1-是',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮箱配置表';

-- 邮件发送日志表
CREATE TABLE IF NOT EXISTS `email_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `to_email` VARCHAR(200) NOT NULL COMMENT '收件人邮箱',
  `subject` VARCHAR(500) NOT NULL COMMENT '邮件主题',
  `body` TEXT COMMENT '邮件内容',
  `status` VARCHAR(20) NOT NULL COMMENT '状态: success, failed',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `sent_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX `idx_to_email` (`to_email`),
  INDEX `idx_status` (`status`),
  INDEX `idx_sent_at` (`sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件发送日志表';

-- 邮件模板表
CREATE TABLE IF NOT EXISTS `email_templates` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `code` VARCHAR(50) NOT NULL UNIQUE COMMENT '模板代码',
  `name` VARCHAR(200) NOT NULL COMMENT '模板名称',
  `subject` VARCHAR(500) NOT NULL COMMENT '邮件主题',
  `body` TEXT NOT NULL COMMENT '邮件内容(支持变量)',
  `variables` JSON DEFAULT NULL COMMENT '可用变量说明',
  `is_html` TINYINT NOT NULL DEFAULT 1 COMMENT '是否HTML: 0-否, 1-是',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件模板表';

-- 插入默认邮件模板
INSERT INTO `email_templates` (`code`, `name`, `subject`, `body`, `variables`, `is_html`) VALUES
('license_expire', '许可证过期提醒', 'LingCDN 许可证即将过期',
'<html><body><h2>许可证过期提醒</h2><p>尊敬的 {{customer_name}}：</p><p>您的许可证 <strong>{{license_code}}</strong> 将于 <strong>{{expire_date}}</strong> 过期。</p><p>请及时续费以避免服务中断。</p><p>LingCDN 团队</p></body></html>',
'{"customer_name": "客户名称", "license_code": "许可证密钥", "expire_date": "过期日期"}', 1),

('system_alert', '系统告警通知', 'LingCDN 系统告警',
'<html><body><h2>系统告警</h2><p>告警时间：{{alert_time}}</p><p>告警级别：{{alert_level}}</p><p>告警内容：{{alert_message}}</p><p>请及时处理。</p></body></html>',
'{"alert_time": "告警时间", "alert_level": "告警级别", "alert_message": "告警内容"}', 1),

('welcome', '欢迎邮件', '欢迎使用 LingCDN',
'<html><body><h2>欢迎使用 LingCDN</h2><p>尊敬的 {{username}}：</p><p>您的账户已创建成功！</p><p>登录地址：{{login_url}}</p><p>用户名：{{username}}</p><p>请及时修改初始密码。</p><p>LingCDN 团队</p></body></html>',
'{"username": "用户名", "login_url": "登录地址"}', 1);
