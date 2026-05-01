#!/bin/bash

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_PORT=8080
FRONTEND_PORT=3000
MAX_RETRIES=30
RETRY_INTERVAL=2

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "$1 未安装，请先安装"
        exit 1
    fi
}

check_mysql() {
    log_info "检查 MySQL 连接..."
    if docker exec lattice-mysql mysqladmin ping -h localhost --silent 2>/dev/null; then
        log_info "MySQL 连接正常"
        return 0
    fi
    
    if command -v mysql &> /dev/null; then
        if mysql -h localhost -P 3306 -u admin -ppassword -e "SELECT 1" &>/dev/null 2>&1; then
            log_info "MySQL 连接正常"
            return 0
        fi
    fi
    
    log_warn "MySQL 连接失败，继续尝试..."
    return 0
}

check_redis() {
    log_info "检查 Redis 连接..."
    if docker exec lattice-redis redis-cli ping 2>/dev/null | grep -q "PONG"; then
        log_info "Redis 连接正常"
        return 0
    fi
    
    if command -v redis-cli &> /dev/null; then
        if redis-cli -h localhost -p 6379 ping 2>/dev/null | grep -q "PONG"; then
            log_info "Redis 连接正常"
            return 0
        fi
    fi
    
    log_warn "Redis 连接失败，继续尝试..."
    return 0
}

build_backend() {
    log_info "构建后端..."
    cd "$PROJECT_ROOT"
    
    if ! go build -o bin/api ./cmd/api/; then
        log_error "后端构建失败"
        exit 1
    fi
    
    log_info "后端构建成功"
}

start_backend() {
    log_info "启动后端服务..."
    
    pkill -f "bin/api" 2>/dev/null || true
    sleep 1
    
    cd "$PROJECT_ROOT"
    nohup ./bin/api > logs/api.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > /tmp/lattice-api.pid
    
    log_info "后端服务已启动 (PID: $BACKEND_PID)"
}

wait_for_backend() {
    log_info "等待后端健康检查通过..."
    
    for i in $(seq 1 $MAX_RETRIES); do
        if curl -s "http://localhost:$BACKEND_PORT/health" | grep -q "success"; then
            log_info "后端健康检查通过"
            return 0
        fi
        echo -n "."
        sleep $RETRY_INTERVAL
    done
    
    echo ""
    log_error "后端健康检查超时，请检查日志: logs/api.log"
    exit 1
}

start_frontend() {
    log_info "启动前端开发服务器..."
    
    cd "$PROJECT_ROOT/lattice-coding-web"
    
    if [ ! -d "node_modules" ]; then
        log_info "安装前端依赖..."
        if command -v npm &> /dev/null; then
            npm install
        elif command -v docker &> /dev/null; then
            docker run --rm -v "$(pwd):/app" -w /app node:20-alpine npm install
        else
            log_error "未找到 npm 或 docker，无法安装前端依赖"
            exit 1
        fi
    fi
    
    if command -v npm &> /dev/null; then
        npm run dev -- --host 0.0.0.0 --port $FRONTEND_PORT
    elif command -v docker &> /dev/null; then
        docker run --rm -p $FRONTEND_PORT:$FRONTEND_PORT -v "$(pwd):/app" -w /app node:20-alpine npm run dev -- --host 0.0.0.0 --port $FRONTEND_PORT
    else
        log_error "未找到 npm 或 docker，无法启动前端"
        exit 1
    fi
}

cleanup() {
    log_info "清理..."
    if [ -f /tmp/lattice-api.pid ]; then
        kill $(cat /tmp/lattice-api.pid) 2>/dev/null || true
        rm -f /tmp/lattice-api.pid
    fi
}

trap cleanup EXIT

main() {
    log_info "启动 Lattice Coding..."
    
    mkdir -p "$PROJECT_ROOT/logs"
    
    check_command docker
    
    check_mysql
    check_redis
    build_backend
    start_backend
    wait_for_backend
    start_frontend
}

main "$@"
