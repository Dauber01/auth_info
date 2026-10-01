.PHONY: help config generate proto wire build run dev clean install-tools migrate seed mod-tidy test test-api fmt lint all sync-harness check-harness

PROJECT_NAME := auth_info
BIN_DIR := bin
OUTPUT := $(BIN_DIR)/$(PROJECT_NAME)
PROTO_DIR := api/proto
GEN_DIR := api/gen
PYTHON ?= python3
ENV ?=
CONFIG_DIR ?= ./config
CONFIG_FILE ?= $(if $(ENV),$(CONFIG_DIR)/$(ENV).yaml,$(CONFIG_DIR))
PROTOC_GO_VERSION := v1.36.11
PROTOC_GRPC_VERSION := v1.5.1
GOPATH_BIN := $(shell go env GOPATH)/bin
export PATH := $(GOPATH_BIN):$(PATH)

## help: 显示可用命令
help:
	@awk '/^## / {sub(/^## /, ""); print}' Makefile

## config: 显示环境与配置入口
config:
	@echo "ENV=$(ENV) CONFIG_FILE=$(CONFIG_FILE)"

## sync-harness: 同步共享 agent 入口
sync-harness:
	@$(PYTHON) .agents/framework/harness.py sync --root .

## check-harness: 检查配置、知识链接与框架测试
check-harness:
	@$(PYTHON) .agents/framework/harness.py check --root .
	@$(PYTHON) -B -m unittest discover -s .agents/framework/tests

## test-api: 使用隔离 HTTP fixture 运行 Python 测试
test-api:
	@$(PYTHON) -B tests/run_api.py

## install-tools: 显式安装固定版本的 Go 生成插件（protoc 需预装）
install-tools:
	@command -v protoc >/dev/null || { echo "Install protoc first"; exit 1; }
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GRPC_VERSION)

## proto: 显式生成 Proto，不自动安装工具
proto:
	@command -v protoc >/dev/null && command -v protoc-gen-go >/dev/null && command -v protoc-gen-go-grpc >/dev/null
	@mkdir -p $(GEN_DIR)
	protoc --proto_path=$(PROTO_DIR)/third_party --proto_path=. --proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative $(PROTO_DIR)/*.proto

## wire: 使用固定版本生成 Wire
wire:
	go run -mod=readonly github.com/google/wire/cmd/wire ./internal/app

## generate: 显式更新所有生成代码
generate: proto wire

## mod-tidy: 显式更新依赖清单
mod-tidy:
	go mod tidy

## build: 仅编译，不更改依赖或生成代码
build:
	@mkdir -p $(BIN_DIR)
	go build -mod=readonly -o $(OUTPUT) ./cmd/main

## run: 构建并使用 CONFIG_FILE 启动
run: build
	$(OUTPUT) -config "$(CONFIG_FILE)"

## dev: 显式生成、构建并启动
dev: generate
	@$(MAKE) run

## migrate: 显式数据库迁移
migrate:
	go run -mod=readonly ./cmd/migrate -config "$(CONFIG_FILE)"

## seed: 显式初始化权限策略
seed:
	go run -mod=readonly ./cmd/seed -config "$(CONFIG_FILE)"

## clean: 仅清理构建产物
clean:
	rm -rf $(BIN_DIR)

## test: 全部 Go 测试
test:
	go test -mod=readonly -count=1 ./...

## fmt: 格式化 Go 文件
fmt:
	go fmt ./...

## lint: Go 静态检查
lint:
	go vet -mod=readonly ./...

## all: 编译并测试
all: build test
