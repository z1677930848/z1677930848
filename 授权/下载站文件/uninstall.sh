#!/bin/bash

# ==========================================
# LingCDN 卸载脚本
# 域名: dl.lingcdn.cloud
# ==========================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认安装目录（可通过环境变量覆盖）
INSTALL_DIR="${LINGCDN_INSTALL_DIR:-/home/lingcdn}"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示 Logo
show_logo() {
    cat << "EOF"
====================================================
           LingCDN 卸载脚本
====================================================
EOF
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "请使用 root 权限运行此脚本"
        exit 1
    fi
}

# 确认卸载
confirm_uninstall() {
    echo ""
    print_warning "此操作将完全卸载 LingCDN 及其所有组件"
    print_warning "安装目录: $INSTALL_DIR"
    echo ""
    read -p "$(echo -e ${YELLOW}是否继续卸载？[y/N]: ${NC})" -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "取消卸载"
        exit 0
    fi
}

# 停止所有服务
stop_services() {
    print_info "停止 LingCDN 服务..."

    # 停止 systemd 服务
    local services=("lingcdnadmin" "ling-admin" "edge-api" "ling-api" "edge-node" "ling-node")
    for service in "${services[@]}"; do
        if systemctl is-active --quiet "$service" 2>/dev/null; then
            print_info "停止服务: $service"
            systemctl stop "$service" || true
        fi

        if systemctl is-enabled --quiet "$service" 2>/dev/null; then
            print_info "禁用服务: $service"
            systemctl disable "$service" || true
        fi

        # 删除 systemd 服务文件
        if [ -f "/etc/systemd/system/$service.service" ]; then
            print_info "删除服务文件: /etc/systemd/system/$service.service"
            rm -f "/etc/systemd/system/$service.service"
        fi
    done

    systemctl daemon-reload || true
    print_success "服务已停止"
}

# 终止进程
kill_processes() {
    print_info "终止 LingCDN 相关进程..."

    local processes=("lingcdnadmin" "ling-admin" "edge-api" "ling-api" "edge-node" "ling-node")
    for proc in "${processes[@]}"; do
        if pgrep -x "$proc" > /dev/null; then
            print_info "终止进程: $proc"
            pkill -9 "$proc" || true
        fi
    done

    print_success "进程已终止"
}

# 删除安装目录
remove_installation() {
    print_info "删除安装目录..."

    if [ -d "$INSTALL_DIR" ]; then
        # 创建备份（可选）
        read -p "$(echo -e ${YELLOW}是否创建备份？[y/N]: ${NC})" -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            BACKUP_FILE="/tmp/lingcdn-backup-$(date +%Y%m%d_%H%M%S).tar.gz"
            print_info "创建备份: $BACKUP_FILE"
            tar -czf "$BACKUP_FILE" -C "$(dirname $INSTALL_DIR)" "$(basename $INSTALL_DIR)" 2>/dev/null || true
            print_success "备份已创建: $BACKUP_FILE"
        fi

        print_info "删除目录: $INSTALL_DIR"
        rm -rf "$INSTALL_DIR"
        print_success "安装目录已删除"
    else
        print_warning "安装目录不存在: $INSTALL_DIR"
    fi
}

# 删除配置文件
remove_configs() {
    print_info "删除配置文件..."

    # 删除可能的配置文件位置
    local config_dirs=(
        "/etc/lingcdn"
        "/etc/edge-admin"
        "/etc/edge-api"
        "/etc/edge-node"
    )

    for dir in "${config_dirs[@]}"; do
        if [ -d "$dir" ]; then
            print_info "删除配置目录: $dir"
            rm -rf "$dir"
        fi
    done

    print_success "配置文件已删除"
}

# 删除日志文件
remove_logs() {
    print_info "删除日志文件..."

    local log_dirs=(
        "/var/log/lingcdn"
        "/var/log/edge-admin"
        "/var/log/edge-api"
        "/var/log/edge-node"
    )

    for dir in "${log_dirs[@]}"; do
        if [ -d "$dir" ]; then
            print_info "删除日志目录: $dir"
            rm -rf "$dir"
        fi
    done

    print_success "日志文件已删除"
}

# 删除数据库（可选）
remove_database() {
    echo ""
    read -p "$(echo -e ${YELLOW}是否删除数据库数据？[y/N]: ${NC})" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_warning "请手动删除数据库，数据库名称通常为: edgeadmin, edgeapi"
        print_info "MySQL 删除命令示例:"
        echo "  DROP DATABASE IF EXISTS edgeadmin;"
        echo "  DROP DATABASE IF EXISTS edgeapi;"
    fi
}

# 主函数
main() {
    show_logo
    check_root

    # 检测安装目录
    if [ ! -d "$INSTALL_DIR" ]; then
        # 尝试检测其他可能的安装位置
        if [ -d "/opt/lingcdn" ]; then
            INSTALL_DIR="/opt/lingcdn"
            print_info "检测到安装目录: $INSTALL_DIR"
        else
            print_error "未找到 LingCDN 安装目录"
            print_info "如果安装在自定义目录，请使用: LINGCDN_INSTALL_DIR=/your/path bash uninstall.sh"
            exit 1
        fi
    fi

    confirm_uninstall

    echo ""
    print_info "开始卸载 LingCDN..."
    echo ""

    stop_services
    kill_processes
    remove_installation
    remove_configs
    remove_logs
    remove_database

    echo ""
    print_success "======================================================"
    print_success "LingCDN 卸载完成！"
    print_success "======================================================"
    echo ""
}

main "$@"
