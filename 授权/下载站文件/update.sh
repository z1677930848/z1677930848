#!/bin/bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

handle_error() {
    local line=$1
    print_error "update script aborted at line ${line}. Please check logs above."
    exit 1
}

trap 'handle_error $LINENO' ERR

DOWNLOAD_HOST_DEFAULT="https://dl.lingcdn.cloud"
DOWNLOAD_HOST="${LINGCDN_DOWNLOAD_HOST:-$DOWNLOAD_HOST_DEFAULT}"
API_VERSION_URL="${LINGCDN_VERSION_API:-${DOWNLOAD_HOST}/api/boot/versions}"
INSTALL_DIR="${LINGCDN_INSTALL_DIR:-/opt/lingcdn}"
BACKUP_DIR="${LINGCDN_BACKUP_DIR:-$INSTALL_DIR/backups}"
CURL_OPTS=(--fail --silent --show-error --location --connect-timeout 10 --max-time 600)
declare -A VERSION_PAYLOAD_CACHE
INSTALLED_COMPONENTS=()
ADMIN_CURRENT_VERSION="unknown"
API_CURRENT_VERSION="unknown"
NODE_CURRENT_VERSION="unknown"
OS=""
ARCH=""

show_logo() {
    cat <<'EOF'
====================================================
           LingCDN 自动更新脚本
====================================================
EOF
}

ensure_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "please run this script with sudo or as root"
        exit 1
    fi
}

detect_system() {
    print_info "detecting system information..."
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
            print_error "unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    print_success "system: $OS, arch: $ARCH"
}

check_commands() {
    print_info "checking required commands..."
    local commands=("curl" "tar" "unzip" "md5sum")
    local missing=()

    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            missing+=("$cmd")
        fi
    done

    if [ ${#missing[@]} -ne 0 ]; then
        print_warning "missing tools: ${missing[*]}"
        if command -v apt-get &> /dev/null; then
            apt-get update -qq
            apt-get install -y curl tar unzip coreutils
        elif command -v yum &> /dev/null; then
            yum install -y curl tar unzip coreutils
        else
            print_error "cannot auto install dependencies, please install: ${missing[*]}"
            exit 1
        fi
    fi

    if ! command -v sha256sum &> /dev/null && ! command -v shasum &> /dev/null; then
        print_warning "sha256sum/shasum not found, SHA256 validation will be skipped"
    fi

    print_success "all required commands are ready"
}

get_component_version() {
    local bin_path=$1
    if [ ! -x "$bin_path" ]; then
        echo "unknown"
        return
    fi
    local output
    if output=$("$bin_path" -v 2>&1); then
        echo "$output" | grep -oP 'v\K[0-9.]+' | head -1 || echo "unknown"
        return
    fi
    if output=$("$bin_path" version 2>&1); then
        echo "$output" | grep -oP 'v\K[0-9.]+' | head -1 || echo "unknown"
        return
    fi
    echo "unknown"
}

detect_installed_components() {
    print_info "detecting installed components..."
    INSTALLED_COMPONENTS=()

    if [ -x "${INSTALL_DIR}/bin/ling-admin" ]; then
        ADMIN_CURRENT_VERSION=$(get_component_version "${INSTALL_DIR}/bin/ling-admin")
        INSTALLED_COMPONENTS+=("admin")
        print_success "found LingCDN Admin v${ADMIN_CURRENT_VERSION}"
    fi

    if [ -x "${INSTALL_DIR}/ling-api/bin/ling-api" ]; then
        API_CURRENT_VERSION=$(get_component_version "${INSTALL_DIR}/ling-api/bin/ling-api")
        INSTALLED_COMPONENTS+=("api")
        print_success "found LingCDN API v${API_CURRENT_VERSION}"
    fi

    local node_bin="${INSTALL_DIR}/ling-node/bin/ling-node"
    if [ ! -x "$node_bin" ] && [ -x "/usr/local/bin/ling-node" ]; then
        node_bin="/usr/local/bin/ling-node"
    fi
    if [ -x "$node_bin" ]; then
        NODE_CURRENT_VERSION=$(get_component_version "$node_bin")
        INSTALLED_COMPONENTS+=("node")
        print_success "found LingCDN Node v${NODE_CURRENT_VERSION}"
    fi

    if [ ${#INSTALLED_COMPONENTS[@]} -eq 0 ]; then
        print_error "no LingCDN components detected. Please run install script first."
        exit 1
    fi
    echo
}

fetch_versions_payload() {
    local key="${OS}_${ARCH}"
    if [ -n "${VERSION_PAYLOAD_CACHE[$key]:-}" ]; then
        printf "%s" "${VERSION_PAYLOAD_CACHE[$key]}"
        return
    fi
    local url="${API_VERSION_URL}?os=${OS}&arch=${ARCH}"
    print_info "fetching version info from ${url}"
    local payload
    payload=$(curl "${CURL_OPTS[@]}" "$url")
    VERSION_PAYLOAD_CACHE[$key]="$payload"
    printf "%s" "$payload"
}

load_latest_version() {
    local component=$1
    local payload tmp
    payload=$(fetch_versions_payload)
    tmp=$(mktemp)
    printf "%s" "$payload" > "$tmp"

    VERSION=""
    DOWNLOAD_URL=""
    FILE_MD5=""
    FILE_SHA256=""

    if command -v jq &> /dev/null; then
        VERSION=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .version" "$tmp")
        local path
        path=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .url" "$tmp")
        DOWNLOAD_URL="${DOWNLOAD_HOST%/}${path}"
        FILE_MD5=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .md5" "$tmp")
        FILE_SHA256=$(jq -r ".data.versions[] | select(.code==\"${component}\") | .sha256" "$tmp")
    else
        local in_block=0
        while IFS= read -r line; do
            if echo "$line" | grep -q '"code"'; then
                if echo "$line" | grep -q "${component}"; then
                    in_block=1
                else
                    in_block=0
                fi
            fi
            if [ $in_block -eq 1 ]; then
                if echo "$line" | grep -q '"version"'; then
                    VERSION=$(echo "$line" | awk -F'"' '{print $4}')
                elif echo "$line" | grep -q '"url"'; then
                    local rel
                    rel=$(echo "$line" | awk -F'"' '{print $4}' | sed 's|\\/|/|g')
                    DOWNLOAD_URL="${DOWNLOAD_HOST%/}${rel}"
                elif echo "$line" | grep -q '"md5"'; then
                    FILE_MD5=$(echo "$line" | awk -F'"' '{print $4}')
                elif echo "$line" | grep -q '"sha256"'; then
                    FILE_SHA256=$(echo "$line" | awk -F'"' '{print $4}')
                fi
                if echo "$line" | grep -q '}' ; then
                    break
                fi
            fi
        done < "$tmp"
    fi
    rm -f "$tmp"

    if [ -z "$VERSION" ] || [ "$VERSION" = "null" ]; then
        print_error "no available version info for ${component}"
        exit 1
    fi
}

version_is_newer() {
    local current=$1
    local latest=$2
    if [ "$current" = "unknown" ]; then
        return 0
    fi
    if [ "$current" = "$latest" ]; then
        return 1
    fi
    local higher
    higher=$(printf "%s\n%s\n" "$current" "$latest" | sort -V | tail -n1)
    if [ "$higher" = "$latest" ]; then
        return 0
    fi
    return 1
}

create_backup() {
    local component=$1
    mkdir -p "$BACKUP_DIR"
    local ts=$(date +%Y%m%d_%H%M%S)
    local backup_file="${BACKUP_DIR}/ling-${component}-${ts}.tar.gz"
    if [ "$component" = "admin" ]; then
        tar -czf "$backup_file" -C "$INSTALL_DIR" bin/ling-admin web configs 2>/dev/null || true
    elif [ "$component" = "api" ]; then
        tar -czf "$backup_file" -C "$INSTALL_DIR" ling-api 2>/dev/null || true
    elif [ "$component" = "node" ]; then
        if [ -d "${INSTALL_DIR}/ling-node" ]; then
            tar -czf "$backup_file" -C "$INSTALL_DIR" ling-node 2>/dev/null || true
        else
            tar -czf "$backup_file" -C "/usr/local/bin" ling-node 2>/dev/null || true
        fi
    fi
    print_success "backup created: $backup_file"
    ls -t "${BACKUP_DIR}"/ling-${component}-*.tar.gz 2>/dev/null | tail -n +6 | xargs rm -f 2>/dev/null || true
}

stop_service() {
    local component=$1
    local service="ling-${component}"
    if systemctl list-unit-files 2>/dev/null | grep -q "^${service}.service"; then
        if systemctl is-active --quiet "$service"; then
            systemctl stop "$service"
            print_success "stopped ${service} via systemd"
        fi
    fi
    local procs=$(ps aux | grep "ling-${component}" | grep -v grep | grep -v defunct | wc -l)
    if [ "$procs" -gt 0 ]; then
        pkill -9 "ling-${component}" 2>/dev/null || true
        sleep 2
        print_success "terminated running processes for ${component}"
    fi
}

verify_integrity() {
    local file=$1
    local ok=false
    if [ -n "$FILE_SHA256" ] && [ "$FILE_SHA256" != "null" ]; then
        local sha=""
        if command -v sha256sum &> /dev/null; then
            sha=$(sha256sum "$file" | awk '{print $1}')
        elif command -v shasum &> /dev/null; then
            sha=$(shasum -a 256 "$file" | awk '{print $1}')
        fi
        if [ -n "$sha" ] && [ "$sha" = "$FILE_SHA256" ]; then
            print_success "SHA256 verification passed"
            ok=true
        else
            print_warning "SHA256 mismatch (expected $FILE_SHA256 got ${sha:-N/A})"
        fi
    fi
    if [ "$ok" = false ] && [ -n "$FILE_MD5" ] && [ "$FILE_MD5" != "null" ]; then
        local md5
        md5=$(md5sum "$file" | awk '{print $1}')
        if [ "$md5" = "$FILE_MD5" ]; then
            print_success "MD5 verification passed"
            ok=true
        else
            print_warning "MD5 mismatch (expected $FILE_MD5 got $md5)"
        fi
    fi
    if [ "$ok" = false ]; then
        print_warning "package integrity could not be confirmed"
    fi
}

download_and_install() {
    local component=$1
    local display=$2
    local tmp_dir
    tmp_dir=$(mktemp -d)
    local file_ext=".tar.gz"
    if [[ "$DOWNLOAD_URL" == *.zip ]]; then
        file_ext=".zip"
    fi
    local package_file="${tmp_dir}/package${file_ext}"
    print_info "downloading ${display} v${VERSION}"
    curl "${CURL_OPTS[@]}" -o "$package_file" "$DOWNLOAD_URL" --progress-bar
    verify_integrity "$package_file"

    if [ "$file_ext" = ".tar.gz" ]; then
        tar -xzf "$package_file" -C "$tmp_dir"
    else
        unzip -q "$package_file" -d "$tmp_dir"
    fi

    local install_path="$INSTALL_DIR"
    if [ "$component" = "api" ]; then
        install_path="$INSTALL_DIR/ling-api"
    elif [ "$component" = "node" ]; then
        install_path="$INSTALL_DIR/ling-node"
    fi
    mkdir -p "$install_path/bin"

    local binary
    binary=$(find "$tmp_dir" -type f \( -name "ling-${component}" -o -name "edge-${component}" -o -name "sk-${component}" \) | head -1)
    if [ -z "$binary" ]; then
        print_error "unable to locate binary for ${component}"
        rm -rf "$tmp_dir"
        exit 1
    fi
    cp "$binary" "${install_path}/bin/ling-${component}"
    chmod +x "${install_path}/bin/ling-${component}"

    if [ "$component" = "admin" ]; then
        local web_dir
        web_dir=$(find "$tmp_dir" -type d -name "web" | head -1)
        if [ -n "$web_dir" ]; then
            rm -rf "${install_path}/web"
            cp -r "$web_dir" "${install_path}/"
            print_success "web assets updated"
        fi
    fi

    rm -rf "$tmp_dir"
    print_success "${display} files installed"
}

start_service() {
    local component=$1
    local service="ling-${component}"
    if systemctl list-unit-files 2>/dev/null | grep -q "^${service}.service"; then
        if systemctl start "$service"; then
            sleep 3
            if systemctl is-active --quiet "$service"; then
                print_success "${service} started via systemd"
                return
            fi
        fi
        print_warning "failed to start ${service} via systemd"
    else
        local install_path="$INSTALL_DIR"
        if [ "$component" = "api" ]; then
            install_path="$INSTALL_DIR/ling-api"
        elif [ "$component" = "node" ]; then
            install_path="$INSTALL_DIR/ling-node"
        fi
        mkdir -p "${install_path}/logs"
        if nohup "${install_path}/bin/ling-${component}" > "${install_path}/logs/ling-${component}.log" 2>&1 & then
            sleep 3
            if ps aux | grep "ling-${component}" | grep -v grep | grep -v defunct > /dev/null; then
                print_success "${component} started in nohup mode"
                return
            fi
        fi
        print_warning "unable to start ${component}, please check manually"
    fi
}

verify_update() {
    local component=$1
    local expected=$2
    local installed="unknown"
    if [ "$component" = "admin" ]; then
        installed=$(get_component_version "${INSTALL_DIR}/bin/ling-admin")
    elif [ "$component" = "api" ]; then
        installed=$(get_component_version "${INSTALL_DIR}/ling-api/bin/ling-api")
    elif [ "$component" = "node" ]; then
        local node_bin="${INSTALL_DIR}/ling-node/bin/ling-node"
        [ ! -x "$node_bin" ] && node_bin="/usr/local/bin/ling-node"
        installed=$(get_component_version "$node_bin")
    fi
    if [ "$installed" = "$expected" ]; then
        print_success "${component} update verified: v${installed}"
    else
        print_warning "${component} verification mismatch (expected v${expected}, got v${installed})"
    fi
}

update_component() {
    local component=$1
    local current=$2
    local name="LingCDN ${component^}"
    echo "=========================================="
    print_info "checking ${name}"
    echo "=========================================="
    load_latest_version "$component"
    if ! version_is_newer "$current" "$VERSION"; then
        print_success "${name} is already up to date"
        return
    fi
    print_warning "${name} will be updated: v${current} -> v${VERSION}"
    create_backup "$component"
    stop_service "$component"
    download_and_install "$component" "$name"
    start_service "$component"
    verify_update "$component" "$VERSION"
}

show_summary() {
    echo ""
    echo "=========================================="
    print_success "update summary"
    echo "=========================================="
    for component in "${INSTALLED_COMPONENTS[@]}"; do
        local service="ling-${component}"
        local status="stopped"
        if systemctl list-unit-files 2>/dev/null | grep -q "^${service}.service"; then
            if systemctl is-active --quiet "$service" 2>/dev/null; then
                status="running"
            fi
        elif ps aux | grep "ling-${component}" | grep -v grep | grep -v defunct > /dev/null; then
            status="running"
        fi
        echo "${component^}: ${status}"
    done
    echo ""
    echo "Backups: ${BACKUP_DIR}"
    echo "=========================================="
}

report_update() {
    curl -s -X POST "${DOWNLOAD_HOST%/}/api/stats/update.php" \
        -d "os=${OS}" \
        -d "arch=${ARCH}" \
        -d "components=${INSTALLED_COMPONENTS[*]}" \
        -d "update_type=script" \
        --max-time 3 >/dev/null 2>&1 || true
}

main() {
    show_logo
    ensure_root
    detect_system
    check_commands
    detect_installed_components
    print_info "start upgrading components: ${INSTALLED_COMPONENTS[*]}"
    for component in "${INSTALLED_COMPONENTS[@]}"; do
        if [ "$component" = "admin" ]; then
            update_component "$component" "$ADMIN_CURRENT_VERSION"
        elif [ "$component" = "api" ]; then
            update_component "$component" "$API_CURRENT_VERSION"
        elif [ "$component" = "node" ]; then
            update_component "$component" "$NODE_CURRENT_VERSION"
        fi
    done
    report_update
    show_summary
    print_success "thanks for using LingCDN"
}

main "$@"