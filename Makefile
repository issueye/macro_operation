.PHONY: all build run test clean deps install

# 变量定义
BINARY_NAME=macro-recorder
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}

# 默认目标
all: deps build

# 安装依赖
deps:
	go mod download
	go mod tidy

# 构建
build:
	@echo "Building ${BINARY_NAME}..."
	go build -ldflags "${LDFLAGS}" -o bin/${BINARY_NAME} ./cmd/macro

# 运行
run:
	go run ./cmd/macro

# 测试
test:
	go test -v -cover ./...

# 测试覆盖率
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码检查
lint:
	go vet ./...
	gofmt -l .

# 格式化代码
fmt:
	gofmt -w .

# 清理
clean:
	rm -rf bin/
	rm -rf coverage.out coverage.html

# 运行所有检查
check: fmt lint test

# 构建所有平台
build-all:
	@echo "Building for all platforms..."
	GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/${BINARY_NAME}-windows-amd64.exe ./cmd/macro
	GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/${BINARY_NAME}-darwin-amd64 ./cmd/macro
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/${BINARY_NAME}-linux-amd64 ./cmd/macro

# 帮助
help:
	@echo "可用命令:"
	@echo "  make deps        - 安装依赖"
	@echo "  make build       - 构建项目"
	@echo "  make run         - 运行项目"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理构建"
	@echo "  make fmt         - 格式化代码"
	@echo "  make lint        - 代码检查"
	@echo "  make build-all   - 构建所有平台"
