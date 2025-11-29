#!/bin/bash

# SKCDN 全新安装脚本 - 强制使用 SK 表前缀
# 使用方法: curl -fsSL http://dl.skcdn.cn/install-fresh.sh | bash

set -e

echo "==========================================="
echo "SKCDN 全新安装 - SK 表前缀版本"
echo "==========================================="
echo ""

# 配置
DOWNLOAD_HOST="http://dl.skcdn.cn"
INSTALL_DIR="/opt/edge"
TMP_DIR="/tmp/sk-install-$$"

# 下载最新的 SK API v1.0.1
echo "[1/5] 下载 SK API v1.0.1..."
mkdir -p "$TMP_DIR"
cd "$TMP_DIR"

# 强制下载最新版本，添加时间戳避免缓存
wget -O sk-api.zip "http://dl.skcdn.cn/updates/api/linux/amd64/sk-api-v1.0.1-linux-amd64.zip?t=$(date +%s)"

echo "[2/5] 解压文件..."
unzip -q sk-api.zip

echo "[3/5] 验证版本..."
VERSION=$(./sk-api-package/bin/sk-api version | grep -oP 'v\K[0-9.]+')
if [ "$VERSION" != "1.0.1" ]; then
    echo "错误: 版本不匹配，期望 1.0.1，实际 $VERSION"
    exit 1
fi
echo "版本验证通过: v$VERSION"

echo "[4/5] 安装程序..."
mkdir -p "$INSTALL_DIR/edge-api/bin"
cp -f sk-api-package/bin/sk-api "$INSTALL_DIR/edge-api/bin/"
chmod +x "$INSTALL_DIR/edge-api/bin/sk-api"
ln -sf sk-api "$INSTALL_DIR/edge-api/bin/edge-api"

echo "[5/5] 清理临时文件..."
cd ~
rm -rf "$TMP_DIR"

echo ""
echo "==========================================="
echo "安装完成！"
echo "==========================================="
echo ""
echo "程序路径: $INSTALL_DIR/edge-api/bin/sk-api"
echo "版本: v$VERSION"
echo ""
echo "重要提示:"
echo "1. 此版本使用 SK 表前缀（如 SKAdmins, SKServers）"
echo "2. 如果数据库中有旧的 edge 表，请先清空数据库"
echo "3. 数据库配置文件: $INSTALL_DIR/edge-api/configs/db.yaml"
echo ""
echo "启动命令:"
echo "  $INSTALL_DIR/edge-api/bin/sk-api"
echo ""
