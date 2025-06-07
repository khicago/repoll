.PHONY: build test clean install run help dev deps fmt lint

# 变量定义
BINARY_NAME=repoll
BUILD_DIR=build
CMD_DIR=cmd/repoll
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# 构建标志
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# 默认目标
all: build

# 显示帮助信息
help:
	@echo "可用的命令："
	@echo "  build      构建二进制文件"
	@echo "  test       运行所有测试"
	@echo "  test-v     运行测试（详细输出）"
	@echo "  coverage   生成测试覆盖率报告"
	@echo "  install    安装到 GOPATH/bin"
	@echo "  clean      清理构建产物"
	@echo "  fmt        格式化代码"
	@echo "  lint       代码静态检查"
	@echo "  deps       下载依赖"
	@echo "  run        运行程序（需要参数：make run ARGS='config.toml'）"
	@echo "  dev        开发模式构建"
	@echo "  release    发布构建"

# 构建
build:
	@echo "🔨 构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 开发构建（快速，无优化）
dev:
	@echo "🔧 开发构建..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# 发布构建（优化）
release:
	@echo "🚀 发布构建..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 详细测试输出
test-v:
	@echo "🧪 运行详细测试..."
	go test -v -race ./...

# 生成覆盖率报告
coverage:
	@echo "📊 生成覆盖率报告..."
	@mkdir -p $(BUILD_DIR)
	go test -coverprofile=$(BUILD_DIR)/coverage.out ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "📊 覆盖率报告: $(BUILD_DIR)/coverage.html"

# 基准测试
bench:
	@echo "⚡ 运行基准测试..."
	go test -bench=. -benchmem ./...

# 安装
install:
	@echo "📦 安装 $(BINARY_NAME)..."
	go install $(LDFLAGS) ./$(CMD_DIR)

# 运行程序
run: build
	@echo "🚀 运行 $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# 清理
clean:
	@echo "🧹 清理构建产物..."
	rm -rf $(BUILD_DIR)
	go clean

# 格式化代码
fmt:
	@echo "💅 格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "🔍 代码静态检查..."
	golangci-lint run

# 下载依赖
deps:
	@echo "📥 下载依赖..."
	go mod download
	go mod tidy

# 更新依赖
deps-update:
	@echo "🔄 更新依赖..."
	go get -u ./...
	go mod tidy

# 生成模拟文件（如果使用 mockgen）
mock:
	@echo "🎭 生成模拟文件..."
	go generate ./...

# 重构项目结构
restructure:
	@echo "🏗️  重构项目结构..."
	chmod +x scripts/restructure.sh
	./scripts/restructure.sh

# 检查代码质量
quality: fmt lint test
	@echo "✅ 代码质量检查完成"

# Docker 构建
docker-build:
	@echo "🐳 构建 Docker 镜像..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# 显示版本信息
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

# 清理 Git 未跟踪的文件
git-clean:
	@echo "🧹 清理未跟踪的文件..."
	git clean -fd 