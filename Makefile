.PHONY: help wire build run test clean dev install lint fmt

# 默认目标
.DEFAULT_GOAL := help

# 项目信息
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)
CMD_PATH=cmd/server

# Go 相关变量
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt

# Wire
WIRE=wire

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装项目依赖
	@echo "📦 安装依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ 依赖安装完成"

wire: ## 生成 Wire 依赖注入代码
	@echo "🔧 生成 Wire 代码..."
	@if ! command -v $(WIRE) > /dev/null; then \
		echo "⚠️  Wire 未安装，正在安装..."; \
		$(GO) install github.com/google/wire/cmd/wire@latest; \
	fi
	cd $(CMD_PATH) && $(WIRE)
	@echo "✅ Wire 代码生成完成"

build: wire ## 编译项目
	@echo "🔨 编译项目..."
	@mkdir -p bin
	$(GO) build -o $(BINARY_PATH) ./$(CMD_PATH)
	@echo "✅ 编译完成: $(BINARY_PATH)"

run: wire ## 运行项目
	@echo "🚀 启动服务..."
	$(GO) run ./$(CMD_PATH)/main.go ./$(CMD_PATH)/wire_gen.go

dev: ## 开发模式（自动重载需要额外工具）
	@echo "🔥 开发模式..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "💡 建议安装 air 实现热重载: go install github.com/cosmtrek/air@latest"; \
		$(GO) run ./$(CMD_PATH)/main.go ./$(CMD_PATH)/wire_gen.go; \
	fi

test: ## 运行测试
	@echo "🧪 运行测试..."
	$(GOTEST) -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "🧪 生成测试覆盖率..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告: coverage.html"

lint: ## 运行代码检查
	@echo "🔍 代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint 未安装，使用 go vet"; \
		$(GOVET) ./...; \
	fi

fmt: ## 格式化代码
	@echo "✨ 格式化代码..."
	@$(GOFMT) -w .
	@echo "✅ 代码格式化完成"

clean: ## 清理编译文件
	@echo "🧹 清理编译文件..."
	@rm -rf bin/
	@rm -rf $(CMD_PATH)/wire_gen.go
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

docker-build: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	docker build -t $(BINARY_NAME):latest .

docker-run: ## 运行 Docker 容器
	@echo "🐳 运行 Docker 容器..."
	docker run -p 8080:8080 $(BINARY_NAME):latest

mod-tidy: ## 整理依赖
	@echo "📦 整理依赖..."
	$(GO) mod tidy
	@echo "✅ 依赖整理完成"

mod-download: ## 下载依赖
	@echo "📦 下载依赖..."
	$(GO) mod download
	@echo "✅ 依赖下载完成"

update-deps: ## 更新依赖
	@echo "⬆️  更新依赖..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "✅ 依赖更新完成"

check: lint test ## 执行所有检查（代码检查 + 测试）

all: clean wire build test ## 执行完整构建流程

.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	$(GO) install github.com/google/wire/cmd/wire@latest
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ 工具安装完成"
