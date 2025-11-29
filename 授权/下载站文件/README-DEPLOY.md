# LingCDN 自动化部署指南

## 🚀 一键自动化部署

### 方式1: 完全自动化（推荐）

```bash
curl -fsSL https://dl.lingcdn.cloud/install-complete.sh | sudo bash
```

**特点**:
- ✅ 自动安装所有依赖
- ✅ 自动配置数据库
- ✅ 自动生成随机密码
- ✅ 3-5分钟完成部署
- ✅ 适合快速部署和测试

---

### 方式2: 自定义管理员信息

```bash
curl -fsSL https://dl.lingcdn.cloud/install-complete.sh | \
  sudo ADMIN_USERNAME=admin \
       ADMIN_EMAIL=admin@example.com \
       ADMIN_PASSWORD=YourSecurePassword123 \
       bash
```

**环境变量说明**:
- `ADMIN_USERNAME`: 管理员用户名（默认: admin）
- `ADMIN_EMAIL`: 管理员邮箱（默认: admin@lingcdn.cloud）
- `ADMIN_PASSWORD`: 管理员密码（默认: 自动生成）

---

### 方式3: 交互式安装

```bash
# 下载脚本
wget https://dl.lingcdn.cloud/install-complete.sh

# 执行安装（会提示输入管理员信息）
sudo bash install-complete.sh
```

---

## 📋 使用示例

### 示例1: 生产环境快速部署

```bash
# 一键部署，自动生成密码
curl -fsSL https://dl.lingcdn.cloud/install-complete.sh | sudo bash

# 安装完成后会显示:
# 管理员账户:
#   用户名: admin
#   密码: xK9mP2nQ7vR4
#
# 数据库信息:
#   数据库: lingcdn
#   用户: lingcdn
#   密码: aB3cD5eF7gH9
```

### 示例2: 自定义管理员账户

```bash
curl -fsSL https://dl.lingcdn.cloud/install-complete.sh | \
  sudo ADMIN_USERNAME=myuser \
       ADMIN_EMAIL=myuser@company.com \
       ADMIN_PASSWORD=MySecurePass123 \
       bash
```

### 示例3: 批量部署脚本

```bash
#!/bin/bash
# deploy-lingcdn.sh

# 配置变量
export ADMIN_USERNAME="admin"
export ADMIN_EMAIL="admin@company.com"
export ADMIN_PASSWORD="SecurePassword123"

# 执行安装
curl -fsSL https://dl.lingcdn.cloud/install-complete.sh | sudo -E bash

# 安装完成后的操作
echo "LingCDN 部署完成"
echo "访问地址: http://$(hostname -I | awk '{print $1}'):7788"
```

---

## 🔧 高级配置

### 自定义安装目录

```bash
curl -fsSL https://dl.lingcdn.cloud/install-complete.sh | \
  sudo INSTALL_DIR=/usr/local/lingcdn \
       bash
```

### 跳过MySQL安装（使用外部数据库）

```bash
# 修改脚本中的 install_mysql 函数
# 或手动配置数据库后再运行脚本
```

---

## 📊 部署时间

| 环境 | 时间 | 说明 |
|------|------|------|
| 全新服务器 | 3-5分钟 | 包含MySQL安装 |
| 已有MySQL | 2-3分钟 | 跳过MySQL安装 |
| 网络较慢 | 5-10分钟 | 下载时间较长 |

---

## ✅ 部署后检查

### 1. 检查服务状态

```bash
systemctl status ling-api ling-admin
```

### 2. 检查端口监听

```bash
netstat -tlnp | grep -E "7788|8001"
```

### 3. 查看日志

```bash
journalctl -u ling-api -u ling-admin -f
```

### 4. 访问管理面板

```bash
# 获取服务器IP
hostname -I

# 浏览器访问
# http://服务器IP:7788
```

---

## 🔒 安全建议

### 1. 修改默认密码

```bash
# 登录后立即修改管理员密码
# 系统设置 -> 管理员 -> 修改密码
```

### 2. 配置防火墙

```bash
# 仅允许特定IP访问
firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="YOUR_IP" port port="7788" protocol="tcp" accept'
firewall-cmd --reload
```

### 3. 配置HTTPS

```bash
# 在管理面板中配置SSL证书
# 系统设置 -> HTTPS -> 上传证书
```

### 4. 定期备份

```bash
# 备份数据库
mysqldump -u lingcdn -p lingcdn > lingcdn_backup_$(date +%Y%m%d).sql

# 备份配置文件
tar -czf lingcdn_config_$(date +%Y%m%d).tar.gz /opt/lingcdn/configs
```

---

## 🐛 故障排查

### 问题1: 服务无法启动

```bash
# 查看详细日志
journalctl -u ling-api -n 100 --no-pager
journalctl -u ling-admin -n 100 --no-pager

# 检查配置文件
cat /opt/lingcdn/ling-api/configs/db.yaml
cat /opt/lingcdn/configs/api.yaml
```

### 问题2: 无法访问管理面板

```bash
# 检查端口
ss -tlnp | grep 7788

# 检查防火墙
firewall-cmd --list-all
ufw status

# 检查进程
ps aux | grep ling-admin
```

### 问题3: 数据库连接失败

```bash
# 测试数据库连接
mysql -u lingcdn -p lingcdn

# 检查MySQL状态
systemctl status mysql
```

---

## 📞 技术支持

- 官网: https://lingcdn.cloud
- 文档: https://docs.lingcdn.cloud
- 下载站: https://dl.lingcdn.cloud
- 问题反馈: support@lingcdn.cloud

---

## 🎯 最佳实践

1. **生产环境**: 使用自定义密码部署
2. **测试环境**: 使用自动生成密码快速部署
3. **批量部署**: 使用脚本统一配置
4. **安全加固**: 部署后立即修改密码、配置防火墙
5. **定期维护**: 定期备份数据库和配置文件
