.PHONY: help proto wire build run clean install-tools install-swag migrate seed mod-tidy docs swagger test fmt lint all

# 变量定义
PROJECT_NAME := auth_info
BIN_DIR := bin
API_DIR := api
PROTO_DIR := $(API_DIR)/proto
GEN_DIR := $(API_DIR)/gen
MAIN_GO := cmd/main/main.go
MIGRATE_GO := cmd/migrate/main.go
SEED_GO := cmd/seed/main.go
OUTPUT := $(BIN_DIR)/$(PROJECT_NAME)
CONFIG_DIR := ./config
SWAG_VERSION := v1.16.6
SWAG_CMD := github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
SWAG_DIRS := ./cmd/main,./internal/handler/auth,./internal/handler/dict,./internal/handler/document,./internal/handler/hello,./api/gen/api/proto
GOPATH_BIN := $(shell go env GOPATH)/bin
SWAG_BIN := $(GOPATH_BIN)/swag
export PATH := $(GOPATH_BIN):$(PATH)

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

## install-swag: 安装 Swagger 文档生成工具
install-swag:
	@if [ ! -x "$(SWAG_BIN)" ]; then \
		echo "$(BLUE)Installing swag $(SWAG_VERSION)...$(NC)"; \
		go install $(SWAG_CMD); \
	fi

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

## migrate: 执行数据库迁移（独立命令）
migrate: mod-tidy
	@echo "$(BLUE)Running migrations...$(NC)"
	@go run $(MIGRATE_GO) -config $(CONFIG_DIR)
	@echo "$(GREEN)✓ Migrations complete$(NC)"

## seed: 初始化默认权限策略（独立命令）
seed: mod-tidy
	@echo "$(BLUE)Seeding default policies...$(NC)"
	@go run $(SEED_GO) -config $(CONFIG_DIR)
	@echo "$(GREEN)✓ Policies seeded$(NC)"

## build: 编译项目
build: mod-tidy proto wire docs
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
docs: install-swag
	@echo "$(BLUE)Generating Swagger docs...$(NC)"
	@$(SWAG_BIN) init -g main.go -d $(SWAG_DIRS) --parseInternal --parseDependency
	@echo "$(GREEN)✓ Swagger docs generated$(NC)"

## swagger: docs 的别名，生成 Swagger 文档
swagger: docs

## all: 执行所有操作（clean, proto, wire, build）
all: clean proto wire build
	@echo "$(GREEN)✓ All operations complete$(NC)"
