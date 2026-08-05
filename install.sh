#!/bin/bash
# ============================================================
#  锚点财务 (AnchorFinance) 安装脚本
#  用法: bash install.sh
#  安装目录: /opt/anchorfinance
# ============================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

INSTALL_DIR="/opt/anchorfinance"

# ============================================================
# 随机密码生成
# ============================================================
generate_password() {
    local len=${1:-16}
    if command -v openssl &> /dev/null; then
        openssl rand -base64 32 | tr -dc 'a-zA-Z0-9!@#$%^&*' | head -c "$len"
    else
        head -c 64 /dev/urandom 2>/dev/null | tr -dc 'a-zA-Z0-9!@#$%^&*' | head -c "$len" || \
        cat /proc/urandom 2>/dev/null | tr -dc 'a-zA-Z0-9!@#$%^&*' | head -c "$len"
    fi
}

# ============================================================
# 1. 检测环境
# ============================================================
check_env() {
    info "检测编译环境..."

    if ! command -v go &> /dev/null; then
        error "未检测到 Go，请先安装: https://go.dev/dl/"
    fi
    ok "Go $(go version | grep -oP 'go\K[0-9]+\.[0-9]+')"

    if ! command -v node &> /dev/null; then
        error "未检测到 Node.js，请先安装: https://nodejs.org/"
    fi
    ok "Node.js $(node --version)"

    if ! command -v mysql &> /dev/null; then
        warn "未检测到 mysql 客户端，将跳过自动建库（需手动导入）"
        HAS_MYSQL=false
    else
        ok "MySQL 客户端 $(mysql --version | grep -oP '[0-9]+\.[0-9]+\.[0-9]+')"
        HAS_MYSQL=true
    fi
}

# ============================================================
# 2. 收集配置
# ============================================================
collect_config() {
    # ── 数据库 ──
    echo ""
    info "══════════════════════════════"
    info "  数据库配置"
    info "══════════════════════════════"

    read -p "  数据库主机 [localhost]: " DB_HOST
    DB_HOST=${DB_HOST:-localhost}

    read -p "  数据库端口 [3306]: " DB_PORT
    DB_PORT=${DB_PORT:-3306}

    read -p "  数据库用户名 [root]: " DB_USER
    DB_USER=${DB_USER:-root}

    read -s -p "  数据库密码: " DB_PASS
    echo ""

    while true; do
        read -p "  数据库名 [anchorfinance]: " DB_NAME
        DB_NAME=${DB_NAME:-anchorfinance}
        if [[ ! "$DB_NAME" =~ ^[a-zA-Z0-9_]{1,64}$ ]]; then
            warn "数据库名只能包含字母、数字和下划线，最长64字符"
            continue
        fi
        break
    done

    # ── 站点 ──
    echo ""
    info "══════════════════════════════"
    info "  站点配置"
    info "══════════════════════════════"

    while true; do
        read -p "  系统名称: " SITE_NAME
        [ -n "$SITE_NAME" ] && break
        warn "系统名称不能为空"
    done

    while true; do
        read -p "  网站域名 [http://localhost]: " SITE_URL
        SITE_URL=${SITE_URL:-http://localhost}
        if [[ ! "$SITE_URL" =~ ^https?:// ]]; then
            warn "网站域名必须以 http:// 或 https:// 开头"
            continue
        fi
        if [[ "$SITE_URL" == */ ]]; then
            warn "网站域名不能以 / 结尾"
            continue
        fi
        break
    done

    generate_admin_path() {
        cat /dev/urandom 2>/dev/null | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1 || echo "admin1panel"
    }

    DEFAULT_ADMIN_PATH=$(generate_admin_path)
    while true; do
        read -p "  后台路径 [$DEFAULT_ADMIN_PATH]: " ADMIN_PATH
        ADMIN_PATH=${ADMIN_PATH:-$DEFAULT_ADMIN_PATH}
        if [[ ! "$ADMIN_PATH" =~ ^[a-zA-Z0-9]{2,32}$ ]]; then
            warn "后台路径只能包含字母和数字，长度2-32位"
            continue
        fi
        if [[ "$ADMIN_PATH" =~ ^[0-9]+$ ]]; then
            warn "后台路径不能全是数字"
            continue
        fi
        if [[ "$ADMIN_PATH" =~ ^[a-zA-Z]+$ ]]; then
            warn "后台路径不能全是字母，必须包含数字"
            continue
        fi
        break
    done
    ok "后台路径: /$ADMIN_PATH"

    read -p "  服务端口 [8080]: " SERVER_PORT
    SERVER_PORT=${SERVER_PORT:-8080}

    read -p "  监听地址 [0.0.0.0]: " SERVER_HOST
    SERVER_HOST=${SERVER_HOST:-0.0.0.0}

    # ── 管理员 ──
    echo ""
    info "══════════════════════════════"
    info "  管理员账号"
    info "══════════════════════════════"

    while true; do
        read -p "  管理员用户名: " ADMIN_USER
        if [ -z "$ADMIN_USER" ]; then
            warn "用户名不能为空"
            continue
        fi
        if [ ${#ADMIN_USER} -lt 4 ] || [ ${#ADMIN_USER} -gt 20 ]; then
            warn "用户名长度必须在4-20位之间"
            continue
        fi
        break
    done

    echo "  密码（留空则随机生成，长度6-32位）:"
    read -s -p "  > " ADMIN_PASS
    echo ""
    if [ -z "$ADMIN_PASS" ]; then
        ADMIN_PASS=$(generate_password 16)
        GENERATED_PASS=true
        ok "已随机生成密码: $ADMIN_PASS"
    else
        GENERATED_PASS=false
        while true; do
            if [ ${#ADMIN_PASS} -lt 6 ]; then
                warn "密码长度至少6位"
                read -s -p "  重新输入: " ADMIN_PASS
                echo ""
                continue
            fi
            if [ ${#ADMIN_PASS} -gt 32 ]; then
                warn "密码长度最多32位"
                read -s -p "  重新输入: " ADMIN_PASS
                echo ""
                continue
            fi
            break
        done
    fi

    while true; do
        read -p "  管理员邮箱: " ADMIN_EMAIL
        [[ "$ADMIN_EMAIL" == *@*.* ]] && break
        warn "请输入有效的邮箱地址"
    done

    # ── 可选：Redis ──
    echo ""
    info "══════════════════════════════"
    info "  可选配置（回车跳过）"
    info "══════════════════════════════"

    REDIS_ENABLED="false"
    read -p "  是否启用 Redis? [y/N]: " ENABLE_REDIS
    if [[ "$ENABLE_REDIS" == "y" || "$ENABLE_REDIS" == "Y" ]]; then
        REDIS_ENABLED="true"
        read -p "  Redis 主机 [127.0.0.1]: " REDIS_HOST
        REDIS_HOST=${REDIS_HOST:-127.0.0.1}
        read -p "  Redis 端口 [6379]: " REDIS_PORT
        REDIS_PORT=${REDIS_PORT:-6379}
        read -s -p "  Redis 密码: " REDIS_PASS
        echo ""
        read -p "  Redis DB [0]: " REDIS_DB
        REDIS_DB=${REDIS_DB:-0}
    fi

    read -p "  安装为 systemd 服务? [Y/n]: " INSTALL_SERVICE
    INSTALL_SERVICE=${INSTALL_SERVICE:-Y}

    # ── 确认 ──
    echo ""
    info "══════════════════════════════"
    info "  确认配置"
    info "══════════════════════════════"
    echo "  安装目录: $INSTALL_DIR"
    echo "  数据库:   $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
    echo "  站点名:   $SITE_NAME"
    echo "  域  名:   $SITE_URL"
    echo "  后  台:   $SITE_URL/$ADMIN_PATH"
    echo "  端  口:   $SERVER_HOST:$SERVER_PORT"
    echo "  管理员:   $ADMIN_USER / $ADMIN_EMAIL"
    if [ "$REDIS_ENABLED" = "true" ]; then
        echo "  Redis:    $REDIS_HOST:$REDIS_PORT (DB $REDIS_DB)"
    fi
    echo ""

    read -p "  确认以上信息? [Y/n]: " CONFIRM
    [[ "$CONFIRM" == "n" || "$CONFIRM" == "N" ]] && error "安装已取消"
}

# ============================================================
# 3. 编译项目
# ============================================================
build_project() {
    info "编译后端..."
    cd "$INSTALL_DIR/backend"
    go mod tidy 2>&1 | tail -1
    CGO_ENABLED=0 go build -ldflags="-s -w" -o "$INSTALL_DIR/anchorfinance" . 2>&1 || error "后端编译失败"
    ok "后端编译完成"

    info "编译用户端前端..."
    cd "$INSTALL_DIR/frontend"
    npm install --legacy-peer-deps 2>&1 | tail -1
    npm run build 2>&1 || error "用户端前端编译失败"
    ok "用户端前端编译完成"

    info "编译管理端前端..."
    cd "$INSTALL_DIR/frontend/admin"
    npm install --legacy-peer-deps 2>&1 | tail -1
    npm run build 2>&1 || error "管理端前端编译失败"
    ok "管理端前端编译完成"

    cd "$INSTALL_DIR"
}

# ============================================================
# 4. 生成 .env 文件
# ============================================================
generate_env() {
    info "生成 .env 配置文件..."

    cat > "$INSTALL_DIR/.env" << EOF
# 锚点财务 环境配置
# 自动生成于 $(date '+%Y-%m-%d %H:%M:%S')
# 此文件仅存储数据库连接信息，其他配置均在数据库中

# 数据库
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASS=$DB_PASS
DB_NAME=$DB_NAME

# 服务器
SERVER_HOST=$SERVER_HOST
SERVER_PORT=$SERVER_PORT
EOF

    chmod 600 "$INSTALL_DIR/.env"
    ok ".env 文件已生成（权限 600）"
}

# ============================================================
# 5. 创建数据库并导入表结构
# ============================================================
setup_database() {
    info "创建数据库..."

    if [ "$HAS_MYSQL" = true ]; then
        export MYSQL_PWD="$DB_PASS"

        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" \
            -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null \
            || error "创建数据库失败，请检查连接信息"
        ok "数据库 $DB_NAME 创建成功"

        info "导入表结构和默认数据..."
        if [ -f "$INSTALL_DIR/init.sql" ]; then
            mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
                < "$INSTALL_DIR/init.sql" 2>/dev/null \
                || error "导入SQL失败"
            ok "表结构和默认数据导入成功"
        else
            error "未找到 init.sql"
        fi

        unset MYSQL_PWD
    else
        warn "无 mysql 客户端，请手动执行："
        echo "    CREATE DATABASE \`$DB_NAME\` DEFAULT CHARACTER SET utf8mb4;"
        echo "    mysql -u $DB_USER -p $DB_NAME < $INSTALL_DIR/init.sql"
    fi
}

# ============================================================
# 6. 写入安装配置到数据库
# ============================================================
write_install_config() {
    info "写入安装配置到数据库..."

    if [ "$HAS_MYSQL" != true ]; then
        warn "无 mysql 客户端，请手动导入管理员和站点配置"
        return
    fi

    # 编译密码哈希工具
    HASH_BIN="$INSTALL_DIR/tools/bcrypt_hash"
    if [ ! -f "$HASH_BIN" ] && [ ! -f "${HASH_BIN}.exe" ]; then
        info "编译密码哈希工具..."
        mkdir -p "$INSTALL_DIR/tools"
        cat > "$INSTALL_DIR/tools/bcrypt_hash.go" << 'GOEOF'
package main

import (
    "fmt"
    "os"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: bcrypt_hash <password>")
        os.Exit(1)
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), 10)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Print(string(hash))
}
GOEOF
        (cd "$INSTALL_DIR/backend" && go build -o "$INSTALL_DIR/tools/bcrypt_hash" "$INSTALL_DIR/tools/bcrypt_hash.go") 2>/dev/null \
            || error "编译密码哈希工具失败"
    fi
    [ -f "${HASH_BIN}.exe" ] && HASH_BIN="${HASH_BIN}.exe"

    ADMIN_PASS_HASH=$("$HASH_BIN" "$ADMIN_PASS") \
        || error "密码哈希计算失败"

    # SQL 转义
    sql_escape() { echo "${1//\'/\\\'}"; }
    E_USER=$(sql_escape "$ADMIN_USER")
    E_HASH=$(sql_escape "$ADMIN_PASS_HASH")
    E_EMAIL=$(sql_escape "$ADMIN_EMAIL")
    E_NAME=$(sql_escape "$SITE_NAME")
    E_URL=$(sql_escape "$SITE_URL")
    E_PATH=$(sql_escape "$ADMIN_PATH")

    export MYSQL_PWD="$DB_PASS"

    # 插入管理员
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
        -e "INSERT INTO users (username, password, email, is_admin, status, created_at, updated_at) VALUES ('${E_USER}', '${E_HASH}', '${E_EMAIL}', 1, 1, NOW(), NOW());" 2>/dev/null \
        || error "插入管理员失败（用户名可能已存在）"
    ok "管理员账号已写入"

    # 更新站点配置
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
        -e "UPDATE system_configs SET value='${E_NAME}' WHERE \`key\`='company_name';" 2>/dev/null
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
        -e "UPDATE system_configs SET value='${E_URL}' WHERE \`key\`='system_url';" 2>/dev/null
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
        -e "UPDATE system_configs SET value='${E_PATH}' WHERE \`key\`='admin_path';" 2>/dev/null
    ok "站点配置已更新"

    # Redis 配置写入数据库
    if [ "$REDIS_ENABLED" = "true" ]; then
        R_HOST=$(sql_escape "$REDIS_HOST")
        R_PORT=$(sql_escape "$REDIS_PORT")
        R_PASS=$(sql_escape "$REDIS_PASS")
        R_DB=$(sql_escape "$REDIS_DB")

        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
            -e "UPDATE system_configs SET value='true' WHERE \`key\`='redis_enabled';" 2>/dev/null
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
            -e "UPDATE system_configs SET value='${R_HOST}' WHERE \`key\`='redis_host';" 2>/dev/null
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
            -e "UPDATE system_configs SET value='${R_PORT}' WHERE \`key\`='redis_port';" 2>/dev/null
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
            -e "UPDATE system_configs SET value='${R_PASS}' WHERE \`key\`='redis_password';" 2>/dev/null
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" \
            -e "UPDATE system_configs SET value='${R_DB}' WHERE \`key\`='redis_db';" 2>/dev/null
        ok "Redis 配置已写入数据库"
    fi

    unset MYSQL_PWD
}

# ============================================================
# 7. 安装 systemd 服务
# ============================================================
install_systemd() {
    [[ "$INSTALL_SERVICE" != "y" && "$INSTALL_SERVICE" != "Y" && "$INSTALL_SERVICE" != "" ]] && return

    info "安装 systemd 服务..."

    cat > /etc/systemd/system/anchorfinance.service << EOF
[Unit]
Description=AnchorFinance Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/anchorfinance
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable anchorfinance
    ok "systemd 服务已安装"
    echo "  启动: systemctl start anchorfinance"
    echo "  状态: systemctl status anchorfinance"
    echo "  日志: journalctl -u anchorfinance -f"
}

# ============================================================
# 8. 完成
# ============================================================
finish() {
    echo ""
    echo "============================================================"
    echo -e "  ${GREEN}安装完成!${NC}"
    echo "============================================================"
    echo ""
    echo "  安装目录: $INSTALL_DIR"
    echo "  管理后台: $SITE_URL/$ADMIN_PATH"
    echo "  用户名:   $ADMIN_USER"
    echo "  密  码:   $ADMIN_PASS"
    echo ""
    if [ "$GENERATED_PASS" = true ]; then
        warn "密码为随机生成，请妥善保存！"
    fi
    echo "  启动方式:"
    echo "    cd $INSTALL_DIR && ./anchorfinance"
    echo ""
    if [[ "$INSTALL_SERVICE" == "y" || "$INSTALL_SERVICE" == "Y" || "$INSTALL_SERVICE" == "" ]]; then
        echo "  或: systemctl start anchorfinance"
        echo ""
    fi
    echo "============================================================"
}

# ============================================================
# 主流程
# ============================================================
main() {
    echo ""
    echo "============================================================"
    echo "  锚点财务 (AnchorFinance) 安装向导"
    echo "============================================================"
    echo ""

    check_env
    collect_config
    build_project
    setup_database
    generate_env
    write_install_config
    install_systemd
    finish
}

main
