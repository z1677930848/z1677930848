#!/bin/bash
# LingCDN 授权安装脚本

set -e

LICENSE_SERVER="https://dl.lingcdn.cloud"
INSTALL_DIR="/home/lingcdn"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
echo_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# 检测系统
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo_error "不支持的架构: $ARCH" ;;
esac

echo_info "系统: $OS/$ARCH"

# 验证授权
read -p "请输入许可证密钥: " LICENSE_KEY
read -p "请输入域名: " DOMAIN
MACHINE_ID=$(cat /etc/machine-id 2>/dev/null || hostname)

echo_info "验证授权..."
RESPONSE=$(curl -s -X POST "$LICENSE_SERVER/api/authorization.php" \
    -H "Content-Type: application/json" \
    -d "{\"license_code\":\"$LICENSE_KEY\",\"system_token\":\"$MACHINE_ID\",\"domain\":\"$DOMAIN\"}")

CODE=$(echo $RESPONSE | grep -o '"code":[0-9]*' | cut -d: -f2)
[ "$CODE" != "200" ] && echo_error "授权验证失败"

echo_info "授权验证成功"

# 选择组件
echo "请选择组件: 1)Admin 2)Node 3)API"
read -p "选择 [1-3]: " CHOICE
case $CHOICE in
    1) COMPONENT="admin" ;;
    2) COMPONENT="node" ;;
    3) COMPONENT="api" ;;
    *) echo_error "无效选项" ;;
esac

# 获取版本
echo_info "获取最新版本..."
RESPONSE=$(curl -s "$LICENSE_SERVER/api/versions.php?action=latest&component=$COMPONENT&os=$OS&arch=$ARCH")
VERSION=$(echo $RESPONSE | grep -o '"version":"[^"]*"' | head -1 | cut -d'"' -f4)
DOWNLOAD_URL=$(echo $RESPONSE | grep -o '"download_url":"[^"]*"' | cut -d'"' -f4)
FILE_MD5=$(echo $RESPONSE | grep -o '"file_md5":"[^"]*"' | cut -d'"' -f4)

[ -z "$VERSION" ] && echo_error "获取版本失败"
echo_info "版本: $VERSION"

# 下载
TEMP_FILE="/tmp/lingcdn-$COMPONENT.zip"
echo_info "下载中..."
curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL" || echo_error "下载失败"

# 校验
echo_info "校验文件..."
ACTUAL_MD5=$(md5sum "$TEMP_FILE" | cut -d' ' -f1)
[ "$ACTUAL_MD5" != "$FILE_MD5" ] && echo_error "文件校验失败"

# 安装
echo_info "安装中..."
mkdir -p "$INSTALL_DIR/$COMPONENT"
unzip -q "$TEMP_FILE" -d "$INSTALL_DIR/$COMPONENT"
chmod +x "$INSTALL_DIR/$COMPONENT"/*
rm -f "$TEMP_FILE"

echo_info "安装完成！目录: $INSTALL_DIR/$COMPONENT"
