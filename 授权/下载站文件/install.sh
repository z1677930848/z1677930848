#!/bin/bash

# ==========================================
# LingCDN 一键安装脚本
# 域名: dl.lingcdn.cloud
# 支持: LingCDN Admin 原生安装向导
# ==========================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量（允许通过环境变量覆盖，默认使用固定官方下载域名）
DOWNLOAD_HOST="${DOWNLOAD_HOST:-https://dl.lingcdn.cloud}"
API_VERSION_URL="${DOWNLOAD_HOST}/api/boot/versions"
INSTALL_DIR="/home/lingcdn"

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

 ╔═══════════════════════════════════════╗
 ║                                      ║
 ║      LingCDN 一键安装脚本            ║
 ║                                      ║
 ║      https://dl.lingcdn.cloud        ║
 ║                                      ║
 ╚═══════════════════════════════════════╝
EOF
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "请使用 root 权限运行此脚本"
        print_info "使用命令: sudo $0"
        exit 1
    fi
}

# 检测操作系统和架构
detect_system() {
    print_info "检测系统信息..."

    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            print_error "不支持的架构: $ARCH"
            exit 1
            ;;
    esac

    print_success "系统: $OS, 架构: $ARCH"
}

# 检查必要的命令
check_commands() {
    print_info "检查必要的工具..."

    local commands=("curl" "tar" "unzip" "systemctl")
    local missing=()

    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            missing+=("$cmd")
        fi
    done

    if [ ${#missing[@]} -ne 0 ]; then
        print_warning "缺少以下工具: ${missing[*]}"
        print_info "正在安装..."

        if command -v apt-get &> /dev/null; then
            apt-get update -qq
            apt-get install -y curl tar unzip systemd
        elif command -v yum &> /dev/null; then
            yum install -y curl tar unzip systemd
        else
            print_error "无法自动安装依赖，请手动安装: ${missing[*]}"
            exit 1
        fi
    fi

    print_success "所有必要工具已就绪"
}

# 获取最新版本信息
get_latest_version() {
    local component=$1
    print_info "正在获取 ${component} 最新版本信息..."

    local url="${API_VERSION_URL}?os=${OS}&arch=${ARCH}"
    local response=$(curl -s "$url")

    if [ $? -ne 0 ]; then
        print_error "无法连接到服务器: $DOWNLOAD_HOST"
        exit 1
    fi

    # 保存响应到临时文件以便多次读取
    local tmp_json=$(mktemp)
    echo "$response" > "$tmp_json"

    # 解析 JSON - 优先使用jq，其次使用awk
    if command -v jq &> /dev/null; then
        # 使用jq解析JSON
        VERSION=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .version" "$tmp_json")
        local path=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .url" "$tmp_json")
        DOWNLOAD_URL="${DOWNLOAD_HOST}${path}"
        FILE_MD5=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .md5" "$tmp_json")
    else
        # 使用awk解析JSON（简单但可靠）
        # 将JSON转换为每行一个字段，然后提取对应组件的信息
        local in_target=0
        while IFS= read -r line; do
            # 检查是否找到目标组件
            if echo "$line" | grep -q "\"code\".*:.*\"${component}\""; then
                in_target=1
            fi

            # 如果在目标组件内，提取字段
            if [ $in_target -eq 1 ]; then
                if echo "$line" | grep -q "\"version\""; then
                    VERSION=$(echo "$line" | awk -F'"' '{print $4}')
                elif echo "$line" | grep -q "\"url\""; then
                    local path=$(echo "$line" | awk -F'"' '{print $4}' | sed 's|\\/|/|g')
                    DOWNLOAD_URL="${DOWNLOAD_HOST}${path}"
                elif echo "$line" | grep -q "\"md5\""; then
                    FILE_MD5=$(echo "$line" | awk -F'"' '{print $4}')
                fi

                # 遇到对象结束符号，停止
                if echo "$line" | grep -q "^[[:space:]]*}"; then
                    break
                fi
            fi
        done < "$tmp_json"
    fi

    # 清理临时文件
    rm -f "$tmp_json"

    if [ -z "$VERSION" ]; then
        print_error "未找到适用的版本"
        exit 1
    fi

    print_success "找到最新版本: $VERSION"
}

# 下载并安装组件
download_and_install() {
    local component=$1
    local component_name=$2

    get_latest_version "$component"

    print_info "开始下载 ${component_name} v${VERSION}..."

    local tmp_dir=$(mktemp -d)

    # 根据文件类型确定文件名和扩展名
    local file_ext=""
    if [[ "$DOWNLOAD_URL" == *.tar.gz ]]; then
        file_ext=".tar.gz"
    elif [[ "$DOWNLOAD_URL" == *.zip ]]; then
        file_ext=".zip"
    else
        file_ext=".zip"  # 默认使用 zip
    fi

    local filename="ling-${component}-v${VERSION}-${OS}-${ARCH}${file_ext}"
    local download_file="${tmp_dir}/${filename}"

    print_info "下载地址: $DOWNLOAD_URL"

    # 下载文件
    if ! curl -L -o "$download_file" "$DOWNLOAD_URL" --progress-bar; then
        print_error "下载失败: $DOWNLOAD_URL"
        rm -rf "$tmp_dir"
        exit 1
    fi

    print_success "下载完成"

    # 验证 MD5（如果有的话）
    if [ -n "$FILE_MD5" ] && [ "$FILE_MD5" != "null" ]; then
        local actual_md5=$(md5sum "$download_file" | awk '{print $1}')
        if [ "$actual_md5" != "$FILE_MD5" ]; then
            print_warning "MD5 校验失败！预期: $FILE_MD5, 实际: $actual_md5"
            print_warning "继续安装可能存在风险"
        else
            print_success "MD5 校验通过"
        fi
    else
        print_info "跳过 MD5 校验"
    fi

    # 解压
    print_info "正在解压..."
    if [[ "$file_ext" == ".tar.gz" ]]; then
        if ! tar -xzf "$download_file" -C "$tmp_dir"; then
            print_error "解压失败！文件可能已损坏"
            rm -rf "$tmp_dir"
            exit 1
        fi
    else
        if ! unzip -q "$download_file" -d "$tmp_dir"; then
            print_error "解压失败！请确保 unzip 工具已安装"
            print_info "尝试手动安装: apt-get install unzip 或 yum install unzip"
            rm -rf "$tmp_dir"
            exit 1
        fi
    fi

    print_success "解压完成"

    # 创建安装目录
    local install_path="${INSTALL_DIR}"
    if [ "$component" == "api" ]; then
        install_path="${INSTALL_DIR}/ling-api"
    fi

    mkdir -p "$install_path/bin"
    mkdir -p "$install_path/configs"

    # 复制二进制文件
    print_info "正在安装到 $install_path..."

    # 查找二进制文件，支持多种命名方式
    local binary_file=""
    local binary_names=()

    # 根据组件类型定义可能的二进制文件名
    if [ "$component" == "admin" ]; then
        binary_names=("lingcdnadmin" "ling-admin" "sk-admin" "edge-admin")
    elif [ "$component" == "api" ]; then
        binary_names=("ling-api" "sk-api" "edge-api")
    elif [ "$component" == "node" ]; then
        binary_names=("ling-node" "sk-node" "edge-node")
    else
        binary_names=("ling-${component}" "sk-${component}" "edge-${component}")
    fi

    # 查找二进制文件
    for name in "${binary_names[@]}"; do
        binary_file=$(find "$tmp_dir" -name "$name" -type f | head -1)
        if [ -n "$binary_file" ]; then
            print_info "找到二进制文件: $name"
            break
        fi
    done

    if [ -n "$binary_file" ]; then
        # 统一使用 ling- 命名
        local target_name="ling-${component}"
        [ "$component" == "admin" ] && target_name="lingcdnadmin"

        cp "$binary_file" "${install_path}/bin/${target_name}"
        chmod +x "${install_path}/bin/${target_name}"
        print_success "二进制文件已安装"
    else
        print_error "未找到二进制文件"
        print_info "查找的文件名: ${binary_names[*]}"
        print_info "临时目录内容:"
        find "$tmp_dir" -type f | head -20
        rm -rf "$tmp_dir"
        exit 1
    fi

    # 处理 web 目录 (LingCDN Admin 需要)
    if [ "$component" == "admin" ]; then
        print_info "正在安装 web 视图文件..."

        # 查找web目录
        local web_dir=$(find "$tmp_dir" -type d -name "web" -path "*/web" | head -1)

        if [ -n "$web_dir" ] && [ -d "$web_dir" ]; then
            # 检查是否包含public目录（前端项目标志）或index.html
            if [ -d "$web_dir/public" ] || [ -f "$web_dir/index.html" ]; then
                # 删除旧的web目录
                rm -rf "${install_path}/web"
                # 复制新的web目录
                cp -r "$web_dir" "${install_path}/"
                print_success "web 目录已安装"
            else
                print_warning "web 目录结构异常，但仍继续安装"
                rm -rf "${install_path}/web"
                cp -r "$web_dir" "${install_path}/"
            fi
        else
            print_error "安装包中未找到 web 目录"
            print_info "尝试的路径: $tmp_dir"
            find "$tmp_dir" -type d -name "web" || true
            rm -rf "$tmp_dir"
            exit 1
        fi

        # 验证 web 目录
        if [ -d "${install_path}/web" ]; then
            print_success "web 目录验证通过"
        else
            print_error "web 目录安装失败"
            rm -rf "$tmp_dir"
            exit 1
        fi
    fi

    # 清理临时文件
    rm -rf "$tmp_dir"

    print_success "${component_name} 安装完成"
}

# 配置 LingCDN Admin
configure_edge_admin() {
    print_info "配置 LingCDN Admin..."

    local admin_config_dir="${INSTALL_DIR}/configs"
    mkdir -p "$admin_config_dir"

    # 生成 server 配置文件
    cat > "${admin_config_dir}/server.yaml" << 'EOF'
# environment code
env: prod

# http
http:
  "on": true
  listen: [ "0.0.0.0:7788" ]

# https
https:
  "on": false
  listen: [ "0.0.0.0:443"]
  cert: ""
  key: ""
EOF

    print_success "LingCDN Admin 服务配置完成"
}

# 配置防火墙
configure_firewall() {
    print_info "配置防火墙..."

    local port="7788"

    if command -v firewall-cmd &> /dev/null && systemctl is-active --quiet firewalld; then
        firewall-cmd --permanent --add-port=${port}/tcp &>/dev/null || true
        firewall-cmd --reload &>/dev/null || true
        print_success "firewalld 规则已添加"
    elif command -v ufw &> /dev/null; then
        ufw allow ${port}/tcp &>/dev/null || true
        print_success "ufw 规则已添加"
    elif command -v iptables &> /dev/null; then
        iptables -I INPUT -p tcp --dport ${port} -j ACCEPT 2>/dev/null || true
        print_success "iptables 规则已添加"
    else
        print_warning "未检测到防火墙，请手动开放端口 ${port}"
    fi

    print_info "如果使用云服务器，请在安全组中开放端口 ${port}"
}

# 配置 systemd 服务
setup_systemd() {
    print_info "配置 LingCDN Admin 系统服务..."

    local service_file="/etc/systemd/system/lingcdnadmin.service"

    cat > "$service_file" << EOF
[Unit]
Description=LingCDN Admin
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/bin/lingcdnadmin
Restart=on-failure
RestartSec=5s
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable lingcdnadmin

    print_success "LingCDN Admin 系统服务配置完成"
}

# 检查并安装MySQL
install_mysql() {
    print_info "检查MySQL数据库..."

    # 检查MySQL是否已安装
    if command -v mysql &> /dev/null; then
        print_success "MySQL已安装"

        # 尝试连接MySQL检查是否正常运行
        if systemctl is-active --quiet mysql || systemctl is-active --quiet mysqld; then
            print_success "MySQL服务正在运行"
            return 0
        else
            print_info "MySQL已安装但未运行，正在启动..."
            systemctl start mysql 2>/dev/null || systemctl start mysqld 2>/dev/null || true
            return 0
        fi
    fi

    # 检查是否在非交互模式下运行
    print_info "MySQL未安装，正在自动安装MySQL数据库..."
    install_mysql_choice="y"

    print_info "开始安装MySQL..."

    # 根据不同的系统安装MySQL
    if command -v apt-get &> /dev/null; then
        # Debian/Ubuntu
        print_info "检测到 Debian/Ubuntu 系统，安装MySQL..."

        export DEBIAN_FRONTEND=noninteractive
        apt-get update -qq

        # 安装MySQL服务器
        apt-get install -y mysql-server mysql-client

        # 启动MySQL
        systemctl start mysql
        systemctl enable mysql

    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        print_info "检测到 CentOS/RHEL 系统，安装MySQL..."

        # 检测系统版本
        if [ -f /etc/redhat-release ]; then
            if grep -q "release 8" /etc/redhat-release; then
                # CentOS 8 / RHEL 8
                yum install -y mysql-server
            elif grep -q "release 7" /etc/redhat-release; then
                # CentOS 7 / RHEL 7
                if [ ! -f /etc/yum.repos.d/mysql-community.repo ]; then
                    yum install -y https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
                fi
                yum install -y mysql-community-server
            else
                yum install -y mysql-server
            fi
        fi

        # 启动MySQL
        systemctl start mysqld
        systemctl enable mysqld

    else
        print_error "不支持的系统，无法自动安装MySQL"
        print_info "请手动安装MySQL后重新运行此脚本"
        exit 1
    fi

    # 等待MySQL启动
    sleep 3

    # 检查MySQL是否成功启动
    if systemctl is-active --quiet mysql || systemctl is-active --quiet mysqld; then
        print_success "MySQL安装并启动成功"

        # 创建数据库和用户
        print_info "正在配置MySQL数据库..."

        # 生成随机密码
        DB_PASSWORD=$(openssl rand -base64 12 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)

        # 创建数据库和用户
        mysql -u root << EOF
CREATE DATABASE IF NOT EXISTS lingcdn CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'lingcdn'@'localhost' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON lingcdn.* TO 'lingcdn'@'localhost';
FLUSH PRIVILEGES;
EOF

        if [ $? -eq 0 ]; then
            print_success "MySQL数据库配置完成"
            print_info "数据库信息:"
            echo "  数据库名: lingcdn"
            echo "  用户名: lingcdn"
            echo "  密码: ${DB_PASSWORD}"
            echo ""
            print_warning "请记录以上数据库信息，安装向导中需要使用"

            # 保存数据库信息到临时文件
            cat > /tmp/lingcdn-db-info.txt << DBINFO
LingCDN 数据库信息
================
数据库主机: 127.0.0.1
数据库端口: 3306
数据库名称: lingcdn
用户名: lingcdn
密码: ${DB_PASSWORD}

请在安装向导中使用以上信息配置数据库连接
DBINFO
            print_info "数据库信息已保存到 /tmp/lingcdn-db-info.txt"

            # 只在交互模式下等待用户确认
            if [ -t 0 ]; then
                echo ""
                echo -n "按回车键继续..."
                read
            else
                echo ""
                sleep 2
            fi
        else
            print_error "MySQL数据库配置失败"
            print_warning "请在安装向导中手动配置数据库"
        fi
    else
        print_error "MySQL启动失败，请检查系统日志"
        print_warning "您仍可以继续安装，但需要在安装向导中配置外部数据库"
    fi
}

# 自动配置API节点
auto_configure_api() {
    print_info "自动配置API节点..."

    # 环境变量或自动生成
    ADMIN_USERNAME=${ADMIN_USERNAME:-admin}

    if [ -z "$ADMIN_PASSWORD" ]; then
        ADMIN_PASSWORD=$(openssl rand -base64 12 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)
        print_warning "自动生成管理员密码: ${ADMIN_PASSWORD}"
    fi

    # 配置数据库
    cat > "${INSTALL_DIR}/ling-api/configs/db.yaml" << EOF
db:
  type: mysql
  dsn: lingcdn:${DB_PASSWORD}@tcp(127.0.0.1:3306)/lingcdn?charset=utf8mb4&timeout=30s
  prefix: edge
EOF

    # 初始化数据库
    print_info "初始化数据库..."
    cd "${INSTALL_DIR}/ling-api"
    ./bin/ling-api setup << SETUP_INPUT
y
${DB_PASSWORD}
${ADMIN_USERNAME}
admin@lingcdn.cloud
${ADMIN_PASSWORD}
${ADMIN_PASSWORD}
SETUP_INPUT

    [ $? -eq 0 ] && print_success "API节点配置完成" || print_error "配置失败"
}

# 自动配置Admin节点
auto_configure_admin() {
    print_info "自动配置Admin节点..."

    local node_id=$(openssl rand -hex 16)
    local node_secret=$(openssl rand -base64 24)

    cat > "${INSTALL_DIR}/configs/api.yaml" << EOF
rpc:
    endpoints:
        - http://127.0.0.1:8001
    disableUpdate: false
nodeId: ${node_id}
secret: ${node_secret}
EOF

    print_success "Admin节点配置完成"
}

# 启动服务
start_service() {
    print_info "启动 LingCDN Admin 服务..."

    if systemctl start lingcdnadmin; then
        sleep 3

        if systemctl is-active --quiet lingcdnadmin; then
            print_success "LingCDN Admin 服务启动成功"

            # 检查端口监听
            if netstat -tlnp 2>/dev/null | grep -q ":7788" || ss -tlnp 2>/dev/null | grep -q ":7788"; then
                print_success "LingCDN Admin 服务运行正常 (端口 7788 已监听)"
            else
                print_warning "端口 7788 尚未监听，服务可能正在初始化"
            fi
        else
            print_warning "服务状态异常，查看日志: journalctl -u lingcdnadmin -f"
        fi
    else
        print_error "LingCDN Admin 服务启动失败"
        print_info "查看日志: journalctl -u lingcdnadmin -n 50"
        exit 1
    fi
}

# 上报安装统计
report_install() {
    # 静默上报，不影响安装过程
    curl -s -X POST "${DOWNLOAD_HOST}/api/stats/install.php" \
        -d "os=${OS}" \
        -d "arch=${ARCH}" \
        -d "version=${VERSION}" \
        -d "install_type=script" \
        --max-time 3 &>/dev/null || true
}

# 显示安装信息
show_info() {
    echo ""
    echo "=========================================="
    print_success "安装完成!"
    echo "=========================================="
    echo ""
    echo "版本: $VERSION"
    echo "安装目录: $INSTALL_DIR"
    echo ""

    local server_ip=$(hostname -I | awk '{print $1}')

    echo "LingCDN Admin 访问地址:"
    echo "  内网: http://${server_ip}:7788"

    # 检测公网IP
    local public_ip=$(curl -s --max-time 3 ifconfig.me 2>/dev/null || curl -s --max-time 3 ip.sb 2>/dev/null || echo "")
    if [ -n "$public_ip" ] && [ "$public_ip" != "$server_ip" ]; then
        echo "  外网: http://${public_ip}:7788"
    fi

    echo ""
    echo "管理员账号:"
    echo "  用户名: ${ADMIN_USERNAME}"
    echo "  密码: ${ADMIN_PASSWORD}"
    echo ""
    echo "数据库信息:"
    echo "  数据库: lingcdn"
    echo "  用户: lingcdn"
    echo "  密码: ${DB_PASSWORD}"
    echo ""

    echo "LingCDN Admin 服务管理:"
    echo "  启动: systemctl start lingcdnadmin"
    echo "  停止: systemctl stop lingcdnadmin"
    echo "  重启: systemctl restart lingcdnadmin"
    echo "  状态: systemctl status lingcdnadmin"
    echo "  日志: journalctl -u lingcdnadmin -f"
    echo ""

    echo "重要提示:"
    echo "  - 请确保云服务器安全组已开放端口 7788"
    echo "  - 建议配置 HTTPS 以保障安全"
    echo "  - 首次安装请立即修改默认密码"
    echo ""
    echo "=========================================="
}

# 下载节点安装包
download_node_package() {
    local component="node"
    get_latest_version "$component"

    local file_ext=".tar.gz"
    [[ "$DOWNLOAD_URL" == *.zip ]] && file_ext=".zip"

    local deploy_dir="${INSTALL_DIR}/ling-api/deploy"
    mkdir -p "$deploy_dir"

    local filename="ling-${component}-${OS}-${ARCH}-v${VERSION}${file_ext}"
    local target_file="${deploy_dir}/${filename}"

    print_info "下载节点安装包到本地: ${target_file} ..."
    if ! curl -L -o "$target_file" "$DOWNLOAD_URL" --progress-bar; then
        print_error "节点安装包下载失败"
        exit 1
    fi

    if [ -n "$FILE_MD5" ] && [ "$FILE_MD5" != "null" ]; then
        local actual_md5
        actual_md5=$(md5sum "$target_file" | awk '{print $1}')
        if [ "$actual_md5" != "$FILE_MD5" ]; then
            print_warning "节点安装包 MD5 校验失败，预期:${FILE_MD5} 实际:${actual_md5}"
        else
            print_success "节点安装包 MD5 校验通过"
        fi
    fi

    print_success "节点安装包已保存: ${target_file}"
}

main() {
    show_logo
    check_root
    detect_system
    check_commands

    print_info "开始安装 LingCDN Admin + API..."

    # 检查并安装MySQL
    install_mysql

    # 安装 LingCDN Admin
    download_and_install "admin" "LingCDN Admin"

    # 安装 LingCDN API
    download_and_install "api" "LingCDN API"
    download_node_package

    # 配置
    configure_edge_admin
    auto_configure_api
    auto_configure_admin
    configure_firewall
    setup_systemd

    # 启动服务
    start_service

    # 上报安装统计
    report_install

    # 显示信息
    show_info

    echo ""
    print_success "感谢使用 LingCDN!"
    echo ""
}

# 错误处理
trap 'print_error "安装过程中出现错误"; exit 1' ERR

# 执行主函数
main "$@"
