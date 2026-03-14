.PHONY: help proto wire build run clean install-tools

# 变量定义
PROJECT_NAME := auth_info
BIN_DIR := bin
API_DIR := api
PROTO_DIR := $(API_DIR)/proto
GEN_DIR := $(API_DIR)/gen
MAIN_GO := cmd/main/main.go
OUTPUT := $(BIN_DIR)/$(PROJECT_NAME)
CONFIG_DIR := ./config

# 颜色定义
BLUE := \033[0;34m
GREEN := \033[0;32m
NC := \033[0m # No Color

## help: 显示帮助信息
help:
	@echo "$(BLUE)Available commands:$(NC)"
	@grep -E '##' Makefile | grep -v grep | sed 's/## //' | awk '{print "  $(GREEN)" $$1 "$(NC) " substr($$0, index($$0, $$2))}'

## install-tools: 安装 protoc, protoc-gen-go, protoc-gen-go-grpc 工具
install-tools:
	@echo "$(BLUE)Checking and installing tools...$(NC)"
	@which protoc > /dev/null 2>&1 || (echo "Installing protoc..." && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
	@which protoc-gen-go > /dev/null 2>&1 || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc > /dev/null 2>&1 || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

## proto: 生成 protobuf Go 代码
proto: install-tools
	@echo "$(BLUE)Generating proto code...$(NC)"
	@mkdir -p $(GEN_DIR)
	@protoc \
		--proto_path=$(PROTO_DIR)/third_party \
		--proto_path=. \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "$(GREEN)✓ Proto code generated$(NC)"

## wire: 生成 Wire 依赖注入代码
wire:
	@echo "$(BLUE)Generating Wire code...$(NC)"
	@go run github.com/google/wire/cmd/wire@latest ./internal/app
	@echo "$(GREEN)✓ Wire code generated$(NC)"

## mod-tidy: 更新依赖
mod-tidy:
	@echo "$(BLUE)Running go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## build: 编译项目
build: mod-tidy proto wire
	@echo "$(BLUE)Building project...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(OUTPUT) $(MAIN_GO)
	@echo "$(GREEN)✓ Build complete: $(OUTPUT)$(NC)"

## run: 生成代码并运行项目
run: clean build
	@echo "$(BLUE)Starting application...$(NC)"
	@$(OUTPUT) -config $(CONFIG_DIR)

## dev: 快速开发模式（不清理生成的代码）
dev: proto wire build
	@echo "$(BLUE)Starting application in dev mode...$(NC)"
	@$(OUTPUT) -config $(CONFIG_DIR)

## clean: 清理生成的文件和可执行文件
clean:
	@echo "$(BLUE)Cleaning up...$(NC)"
	@rm -rf $(BIN_DIR)
	@rm -rf $(GEN_DIR)/*.pb.go $(GEN_DIR)/*_grpc.pb.go
	@echo "$(GREEN)✓ Clean complete$(NC)"

## test: 运行单元测试
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

## fmt: 格式化代码
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## lint: 运行代码检查
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Lint complete$(NC)"

## docs: 生成 Swagger 文档
docs:
	@echo "$(BLUE)Generating Swagger docs...$(NC)"
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main/main.go
	@echo "$(GREEN)✓ Swagger docs generated$(NC)"

## all: 执行所有操作（clean, proto, wire, build）
all: clean proto wire build
	@echo "$(GREEN)✓ All operations complete$(NC)"
