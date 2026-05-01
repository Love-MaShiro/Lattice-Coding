.PHONY: start stop restart build clean package help

PROJECT_NAME := lattice-coding
VERSION := 1.0.0
BUILD_DIR := bin
DIST_DIR := dist
FRONTEND_DIR := lattice-coding-web

help:
	@echo "Lattice Coding Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make start     - 启动服务"
	@echo "  make stop      - 停止服务"
	@echo "  make restart   - 重启服务"
	@echo "  make build     - 构建后端和前端"
	@echo "  make clean     - 清理构建产物"
	@echo "  make package   - 打包成 tar.gz"
	@echo ""

start:
	@echo "[Make] 启动服务..."
	@./start.sh

stop:
	@echo "[Make] 停止服务..."
	@./stop.sh

restart: stop
	@echo "[Make] 重启服务..."
	@sleep 2
	@./start.sh

build: build-backend build-frontend
	@echo "[Make] 构建完成"

build-backend:
	@echo "[Make] 构建后端..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/api ./cmd/api/
	@go build -o $(BUILD_DIR)/worker ./cmd/worker/
	@echo "[Make] 后端构建完成: $(BUILD_DIR)/api, $(BUILD_DIR)/worker"

build-frontend:
	@echo "[Make] 构建前端..."
	@cd $(FRONTEND_DIR) && npm run build
	@echo "[Make] 前端构建完成: $(FRONTEND_DIR)/dist"

clean:
	@echo "[Make] 清理构建产物..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -f /tmp/lattice-api.pid
	@rm -f /tmp/lattice-web.pid
	@echo "[Make] 清理完成"

package: build
	@echo "[Make] 打包..."
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)
	@cp -r $(BUILD_DIR) $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)/
	@cp -r configs $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)/
	@cp -r $(FRONTEND_DIR)/dist $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)/web
	@cp start.sh stop.sh Makefile $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)/
	@cp README.md $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)/ 2>/dev/null || true
	@cd $(DIST_DIR) && tar -czvf $(PROJECT_NAME)-$(VERSION).tar.gz $(PROJECT_NAME)-$(VERSION)
	@rm -rf $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)
	@echo "[Make] 打包完成: $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION).tar.gz"
