#!/bin/bash
# ============================================================
#  锚点财务 (AnchorFinance) 编译脚本
#  用法: bash build.sh [--skip-frontend] [--skip-backend]
#  输出: /opt/anchorfinance
# ============================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_DIR="/opt/anchorfinance"

SKIP_FRONTEND=false
SKIP_BACKEND=false

for arg in "$@"; do
    case $arg in
        --skip-frontend) SKIP_FRONTEND=true ;;
        --skip-backend)  SKIP_BACKEND=true ;;
        --help|-h)
            echo "用法: bash build.sh [选项]"
            echo ""
            echo "选项:"
            echo "  --skip-frontend  跳过前端编译"
            echo "  --skip-backend   跳过后端编译"
            echo "  --help           显示帮助"
            echo ""
            echo "输出目录: $INSTALL_DIR"
            exit 0
            ;;
    esac
done

# 读取版本号
APP_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null | tr -d '[:space:]' || echo "1.0.0")
info "版本号: $APP_VERSION"

# 创建安装目录
mkdir -p "$INSTALL_DIR"

info "检测编译环境..."

if [ "$SKIP_BACKEND" = false ]; then
    if ! command -v go &> /dev/null; then
        error "未检测到 Go，请先安装: https://go.dev/dl/"
    fi
    ok "Go $(go version | grep -oP 'go\K[0-9]+\.[0-9]+')"
fi

if [ "$SKIP_FRONTEND" = false ]; then
    if ! command -v node &> /dev/null; then
        error "未检测到 Node.js，请先安装: https://nodejs.org/"
    fi
    ok "Node.js $(node --version)"
fi

# 编译后端
if [ "$SKIP_BACKEND" = false ]; then
    info "编译后端..."
    cd "$SCRIPT_DIR/backend"
    go mod tidy 2>&1 | tail -1
    CGO_ENABLED=0 go build -ldflags="-s -w" -o "$INSTALL_DIR/anchorfinance" . 2>&1 || error "后端编译失败"
    ok "后端编译完成: $INSTALL_DIR/anchorfinance"
    cd "$SCRIPT_DIR"
fi

# 编译前端（用户端 + 管理端）
if [ "$SKIP_FRONTEND" = false ]; then
    info "编译用户端前端..."
    cd "$SCRIPT_DIR/frontend"
    npm install --legacy-peer-deps 2>&1 | tail -1
    VITE_VERSION="$APP_VERSION" npm run build 2>&1 | tail -5 || error "用户端前端编译失败"
    mkdir -p "$INSTALL_DIR/frontend"
    cp -r dist/* "$INSTALL_DIR/frontend/"
    ok "用户端前端编译完成: $INSTALL_DIR/frontend/"

    info "编译管理端前端..."
    cd "$SCRIPT_DIR/frontend/admin"
    npm install --legacy-peer-deps 2>&1 | tail -1
    VITE_VERSION="$APP_VERSION" npm run build 2>&1 | tail -5 || error "管理端前端编译失败"
    mkdir -p "$INSTALL_DIR/admin"
    cp -r dist/* "$INSTALL_DIR/admin/"
    ok "管理端前端编译完成: $INSTALL_DIR/admin/"
    cd "$SCRIPT_DIR"
fi

# 复制配置文件到安装目录
info "复制配置文件..."
cp "$SCRIPT_DIR/init.sql" "$INSTALL_DIR/init.sql" 2>/dev/null || true
cp "$SCRIPT_DIR/install.sh" "$INSTALL_DIR/install.sh" 2>/dev/null || true
cp "$SCRIPT_DIR/README.md" "$INSTALL_DIR/README.md" 2>/dev/null || true

# 复制logo
cp "$SCRIPT_DIR/frontend/public/logo.png" "$INSTALL_DIR/frontend/logo.png" 2>/dev/null || true
cp "$SCRIPT_DIR/frontend/admin/public/logo.png" "$INSTALL_DIR/admin/logo.png" 2>/dev/null || true

# 清理
info "清理临时文件..."
rm -rf "$SCRIPT_DIR/backend/tmp" 2>/dev/null || true

echo ""
echo "============================================================"
echo -e "  ${GREEN}编译完成!${NC}"
echo "============================================================"
echo ""
echo "  安装目录: $INSTALL_DIR"
echo "  可执行文件: $INSTALL_DIR/anchorfinance"
echo "  用户端资源: $INSTALL_DIR/frontend/"
echo "  管理端资源: $INSTALL_DIR/admin/"
echo ""
echo "  安装方式: cd $INSTALL_DIR && bash install.sh"
echo ""
echo "============================================================"
