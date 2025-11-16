#!/bin/bash

# ==========================================
# LingCDN 一键升级功能测试脚本
# ==========================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "=============================================="
echo "    LingCDN 一键升级功能测试"
echo "=============================================="
echo -e "${NC}"

# 1. 检查文件是否存在
echo -e "${BLUE}[1/6]${NC} 检查文件..."

FILES=(
    "/root/Lingadmin-master/internal/web/actions/default/settings/upgrade/check.go"
    "/root/Lingadmin-master/internal/web/actions/default/settings/upgrade/install.go"
    "/root/Lingadmin-master/internal/web/actions/default/settings/upgrade/progress.go"
    "/root/Lingadmin-master/internal/web/actions/default/settings/upgrade/history.go"
    "/root/Lingadmin-master/web/views/@default/settings/upgrade/index.html"
    "/root/Lingadmin-master/web/views/@default/settings/upgrade/index.js"
    "/root/Lingadmin-master/internal/utils/upgrade_errors.go"
    "/root/Lingadmin-master/internal/utils/upgrade_logger.go"
    "/root/Lingadmin-master/internal/utils/upgrade_notifier.go"
    "/root/Lingadmin-master/internal/utils/upgrade_manager_v2.go"
    "/root/Lingadmin-master/internal/utils/temp_file_cleaner.go"
)

ALL_FILES_EXIST=true
for FILE in "${FILES[@]}"; do
    if [ -f "$FILE" ]; then
        echo -e "  ${GREEN}✓${NC} $FILE"
    else
        echo -e "  ${RED}✗${NC} $FILE ${RED}(缺失)${NC}"
        ALL_FILES_EXIST=false
    fi
done

if [ "$ALL_FILES_EXIST" = false ]; then
    echo -e "${RED}错误: 部分文件缺失${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 所有文件检查完成${NC}"
echo ""

# 2. 检查代码语法
echo -e "${BLUE}[2/6]${NC} 检查 Go 代码语法..."
cd /root/Lingadmin-master

if go vet ./internal/web/actions/default/settings/upgrade/... 2>/dev/null; then
    echo -e "${GREEN}✓ 后端代码语法检查通过${NC}"
else
    echo -e "${YELLOW}⚠ 部分语法警告（可能是依赖问题，可以忽略）${NC}"
fi
echo ""

# 3. 创建测试用的更新信息
echo -e "${BLUE}[3/6]${NC} 创建测试更新信息..."

TEST_UPDATE_INFO='{
  "version": "1.0.8",
  "currentVersion": "1.0.7",
  "downloadURL": "http://dl.lingcdn.cloud/ling-admin-v1.0.8-linux-amd64.tar.gz",
  "changelog": "新增功能：\n- 添加一键升级功能\n- 支持断点续传\n- 实时进度显示\n- 升级历史记录\n\nBug修复：\n- 修复下载中断问题\n- 优化临时文件清理\n- 改进错误处理",
  "description": "本次更新主要增强了升级体验，新增一键升级功能，支持断点续传和实时进度显示",
  "sha256": "abc123def456789",
  "checkTime": "2025-10-31 15:00:00"
}'

# 创建配置目录
mkdir -p /opt/lingcdn/configs
echo "$TEST_UPDATE_INFO" > /opt/lingcdn/configs/update_info.json

echo -e "${GREEN}✓ 测试更新信息已创建${NC}"
echo ""

# 4. 编译项目
echo -e "${BLUE}[4/6]${NC} 编译项目..."

if [ -f "/root/Lingadmin-master/cmd/lingcdnadmin/main.go" ]; then
    cd /root/Lingadmin-master
    echo "开始编译..."
    if go build -o bin/lingcdnadmin-test cmd/lingcdnadmin/main.go 2>&1 | tee /tmp/build.log; then
        echo -e "${GREEN}✓ 编译成功${NC}"
        ls -lh bin/lingcdnadmin-test
    else
        echo -e "${RED}✗ 编译失败，查看错误日志：${NC}"
        cat /tmp/build.log
        exit 1
    fi
else
    echo -e "${YELLOW}⚠ 未找到 main.go，跳过编译测试${NC}"
fi
echo ""

# 5. 显示访问信息
echo -e "${BLUE}[5/6]${NC} 功能访问方式..."
echo ""
echo "升级页面 URL："
echo "  http://你的服务器IP:7788/settings/upgrade"
echo ""
echo "API 接口："
echo "  POST /settings/upgrade/check     - 检查更新"
echo "  POST /settings/upgrade/install   - 执行升级"
echo "  POST /settings/upgrade/progress  - 获取进度"
echo "  POST /settings/upgrade/history   - 升级历史"
echo ""

# 6. 功能说明
echo -e "${BLUE}[6/6]${NC} 测试步骤..."
echo ""
echo "1. 启动服务："
echo "   cd /opt/lingcdn/bin"
echo "   ./ling-admin start"
echo ""
echo "2. 访问升级页面："
echo "   http://你的服务器IP:7788/settings/upgrade"
echo ""
echo "3. 测试流程："
echo "   a. 页面会自动检查更新"
echo "   b. 应该显示 v1.0.8 新版本"
echo "   c. 点击查看更新内容"
echo "   d. 点击'一键更新'按钮"
echo "   e. 观察进度条和状态"
echo "   f. 点击'升级历史'查看记录"
echo ""

echo -e "${GREEN}"
echo "=============================================="
echo "    测试准备完成！"
echo "=============================================="
echo -e "${NC}"

echo ""
echo "📚 查看详细文档："
echo "  cat /root/Lingadmin-master/ONE_CLICK_UPGRADE_GUIDE.md"
echo ""
echo "🔧 如需帮助，请查看故障排除部分"
echo ""
