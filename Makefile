.PHONY: all deps proto build build-engine build-ui run-engine run-ui clean help

# 默认目标
all: deps proto build

# 安装依赖
deps:
	@echo "Installing Go dependencies..."
	cd macro-engine && go mod tidy
	cd macro-ui && go mod tidy
	@echo "Installing frontend dependencies..."
	cd macro-ui/frontend && npm install

# 生成 proto 文件
proto:
	@echo "Generating gRPC code from proto files..."
	@if command -v protoc >/dev/null 2>&1; then \
		protoc --go_out=./macro-engine --go_opt=paths=source_relative \
		--go-grpc_out=./macro-engine --go-grpc_opt=paths=source_relative \
		-I./macro-engine/proto ./macro-engine/proto/*.proto; \
		echo "Proto files generated successfully"; \
	else \
		echo "protoc not found, please install protoc and protoc-gen-go"; \
		echo "Windows: choco install protoc"; \
		echo "Or download from: https://github.com/protocolbuffers/protobuf/releases"; \
	fi

# 构建引擎服务
build-engine:
	@echo "Building macro-engine..."
	cd macro-engine && go build -o ../bin/macro-engine.exe ./cmd

# 构建 UI 应用
build-ui: deps proto
	@echo "Building macro-ui..."
	cd macro-ui/frontend && npm run build
	cd macro-ui && wails build

# 构建所有
build: build-engine build-ui

# 运行引擎服务（开发模式）
run-engine:
	@echo "Starting macro-engine gRPC server on port 50051..."
	cd macro-engine && go run ./cmd

# 运行 UI 应用
run-ui:
	@echo "Starting macro-ui (Wails)..."
	cd macro-ui && wails dev

# 启动完整应用（需要先启动引擎）
run-all: run-engine &
	@sleep 3
	run-ui

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf macro-ui/frontend/dist/
	rm -f macro-ui/macro-ui.exe
	go clean -cache

# 初始化项目（首次克隆后运行）
init: deps proto build

# 旧版命令兼容
deps-legacy:
	go mod download
	go mod tidy

build-legacy:
	@echo "Legacy build command - use 'make build' for new architecture"

run-legacy:
	@echo "Legacy run command - use 'make run-ui' for new architecture"

# 帮助信息
help:
	@echo "Macro Operation - 项目构建脚本 (gRPC 分离架构)"
	@echo ""
	@echo "使用说明:"
	@echo "  make all          - 安装依赖、生成 proto、构建所有模块"
	@echo "  make deps         - 安装 Go 和前端依赖"
	@echo "  make proto        - 从 proto 文件生成 gRPC 代码"
	@echo "  make build        - 构建所有模块"
	@echo "  make build-engine - 只构建引擎服务 (macro-engine)"
	@echo "  make build-ui     - 只构建 UI 应用 (macro-ui)"
	@echo "  make run-engine   - 运行引擎服务 (端口 50051)"
	@echo "  make run-ui       - 运行 UI 应用 (Wails)"
	@echo "  make run-all      - 运行引擎和 UI（需要两个终端）"
	@echo "  make clean        - 清理构建产物"
	@echo "  make init         - 初始化项目（首次使用）"
	@echo "  make help         - 显示此帮助信息"
	@echo ""
	@echo "注意: 运行前请确保已安装:"
	@echo "  - Go 1.21+"
	@echo "  - Node.js 18+"
	@echo "  - protoc (用于生成 gRPC 代码)"
	@echo "  - Wails CLI: go install github.com/wailsapp/wails/cmd/wails@latest"
