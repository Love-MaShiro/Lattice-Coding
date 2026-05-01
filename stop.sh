#!/bin/bash

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_PID_FILE="/tmp/lattice-api.pid"
FRONTEND_PID_FILE="/tmp/lattice-web.pid"
GRACEFUL_TIMEOUT=10

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

stop_process() {
    local pid_file=$1
    local name=$2
    
    if [ ! -f "$pid_file" ]; then
        log_warn "$name PID 文件不存在: $pid_file"
        return 0
    fi
    
    local pid=$(cat "$pid_file")
    
    if ! kill -0 "$pid" 2>/dev/null; then
        log_warn "$name 进程已不存在 (PID: $pid)"
        rm -f "$pid_file"
        return 0
    fi
    
    log_info "停止 $name (PID: $pid)..."
    kill -TERM "$pid" 2>/dev/null
    
    local count=0
    while [ $count -lt $GRACEFUL_TIMEOUT ]; do
        if ! kill -0 "$pid" 2>/dev/null; then
            log_info "$name 已优雅停止"
            rm -f "$pid_file"
            return 0
        fi
        sleep 1
        count=$((count + 1))
        echo -n "."
    done
    
    echo ""
    log_warn "$name 优雅停止超时，强制终止..."
    kill -KILL "$pid" 2>/dev/null
    rm -f "$pid_file"
    log_info "$name 已强制停止"
}

stop_docker_containers() {
    log_info "停止 Docker 容器中的前端服务..."
    docker stop lattice-web 2>/dev/null || true
    docker rm lattice-web 2>/dev/null || true
}

stop_by_name() {
    local name=$1
    local pattern=$2
    
    local pids=$(pgrep -f "$pattern")
    if [ -n "$pids" ]; then
        log_info "停止 $name 进程..."
        for pid in $pids; do
            if kill -0 "$pid" 2>/dev/null; then
                log_info "停止 $name (PID: $pid)..."
                kill -TERM "$pid" 2>/dev/null
                
                local count=0
                while [ $count -lt $GRACEFUL_TIMEOUT ]; do
                    if ! kill -0 "$pid" 2>/dev/null; then
                        break
                    fi
                    sleep 1
                    count=$((count + 1))
                done
                
                if kill -0 "$pid" 2>/dev/null; then
                    kill -KILL "$pid" 2>/dev/null
                fi
            fi
        done
        log_info "$name 已停止"
    fi
}

main() {
    log_info "停止 Lattice Coding 服务..."
    
    stop_process "$BACKEND_PID_FILE" "后端"
    
    stop_process "$FRONTEND_PID_FILE" "前端"
    
    stop_docker_containers
    
    stop_by_name "后端" "bin/api"
    stop_by_name "前端" "vite"
    
    log_info "所有服务已停止"
}

main "$@"
