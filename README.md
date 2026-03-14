# Auth Info - Gin + Wire + gRPC + Proto 项目框架

这是一个使用 Gin、Google Wire、Protocol Buffers 和 gRPC 的完整项目框架示例。

## 项目特性

✅ **已实现：**
- Gin Web 框架（REST API）
- Google Wire 依赖注入
- Protocol Buffers + gRPC 服务
- gRPC Gateway（REST API 映射）
- Viper 配置管理
- Zap 结构化日志
- Swagger API 文档
- 错误处理中间件
- 优雅关闭机制
- Makefile 工作流自动化

## 项目结构

```
.
├── cmd/main/                    # 应用入口
├── internal/                    # 应用内部代码
│   ├── app/                     # 应用装配（服务启动、依赖注入、路由组装入口）
│   ├── router/                  # HTTP 路由注册（请求路径按模块拆分）
│   ├── handler/                 # REST API 处理器（请求绑定/转换/响应）
│   ├── validation/              # Proto 参数校验（Protovalidate 封装与错误映射）
│   ├── biz/                     # 业务用例层
│   ├── service/                 # gRPC 服务实现与注册
│   ├── data/                    # 持久化层（模型 + 仓储 + DB）
│   ├── middleware/              # 中间件（JWT/Casbin/统一错误处理）
│   ├── config/                  # 配置管理
│   └── logger/                  # 日志组件
├── api/                         # API 定义
│   ├── proto/                   # Proto 契约与依赖（统一目录）
│   │   ├── *.proto              # 业务契约（common/auth/dict/document/hello）
│   │   ├── buf/validate/        # Protovalidate 规则定义
│   │   ├── google/protobuf/     # 仓库内维护的 protobuf（如 struct.proto）
│   │   └── third_party/google/  # 第三方 protobuf 依赖
│   └── gen/                     # 生成的 Proto 代码
├── config/                      # 配置文件
├── docs/                        # Swagger 文档
├── Makefile                     # 构建脚本
└── go.mod                       # Go 模块定义
```

## 契约约定

- 所有对外的请求结构和响应结构统一由 `api/proto/` 生成，HTTP 和 gRPC 共用同一套契约。
- `api/proto/` 统一管理业务 proto 与依赖：业务文件（`common.proto`、`auth.proto`、`dict.proto`、`document.proto`、`hello.proto`）+ `buf/validate` + `third_party/google/protobuf`。
- 所有业务 proto 统一使用 `option go_package = "auth_info/api/gen/api/proto;apipb"`，生成代码集中在 `api/gen/api/proto/`。
- 参数校验统一使用 Protovalidate，规则在业务 proto 中通过 `buf.validate` 注解声明（例如 `(buf.validate.field).string.max_len`）。
- 为了支持 `buf.validate` 导入，仓库内提供 `api/proto/buf/validate/validate.proto`，并在 proto 生成时额外包含 `--proto_path=api/proto/third_party --proto_path=. --proto_path=api/proto`。

## Makefile 命令

### 快速开始

```bash
# 一键启动（自动生成代码并运行）
make run

# 开发模式（不清理已生成的代码）
make dev
```

### 完整命令列表

```bash
# 显示帮助
make help

# 安装必需工具（protoc, protoc-gen-go, protoc-gen-go-grpc）
make install-tools

# 生成 Proto 代码
make proto

# 生成 Wire 依赖注入代码
make wire

# 下载/更新依赖
make mod-tidy

# 编译项目
make build

# 运行项目（生成代码 + 编译 + 启动）
make run

# 开发模式运行（不清理旧代码）
make dev

# 清理生成的文件
make clean

# 运行测试
make test

# 格式化代码
make fmt

# 代码检查
make lint

# 生成 Swagger 文档
make docs

# 执行所有操作（clean, proto, wire, build）
make all
```

## 快速开始

### 1. 安装必需工具

```bash
make install-tools
```

### 2. 生成代码

```bash
make proto      # 生成 Proto 代码
make wire       # 生成 Wire 依赖注入代码
```

### 3. 编译并运行

```bash
make run
```

## 服务端口配置

在 `config/config.yaml` 中配置：

```yaml
server:
  port: 8080      # REST API 监听端口
  mode: debug

log:
  level: debug
  format: json
```

**自动端口分配：**
- REST API HTTP 服务：`http://localhost:8080`
- gRPC 服务：`localhost:9080`（端口 = HTTP 端口 + 1000）

## API 访问

### REST API

```bash
# Hello 接口
curl http://localhost:8080/api/v1/hello

# Swagger 文档
open http://localhost:8080/swagger/index.html
```

### gRPC 服务

使用 `grpcurl` 或其他 gRPC 客户端测试：

```bash
grpcurl -plaintext localhost:9080 list
grpcurl -plaintext localhost:9080 hello.HelloService/SayHello
```

## 添加新的 Proto 定义

### 1. 创建 Proto 文件

在 `api/proto/` 目录下按业务新增 `.proto` 文件，例如 `api/proto/user.proto`：

```protobuf
syntax = "proto3";

package api;

option go_package = "auth_info/api/gen/api/proto;apipb";

message GetUserRequest {
  uint64 id = 1;
}

message User {
  uint64 id = 1;
  string name = 2;
  string email = 3;
}

service UserService {
  rpc GetUser(GetUserRequest) returns (User);
}
```

### 2. 生成代码

```bash
make proto
```

这将在 `api/gen/` 目录下生成对应的 Go 代码。

### 3. 实现服务

在 `internal/service/` 中创建服务实现文件 `user.go`，实现生成的接口。

### 4. 注册服务与路由

- 在 `internal/service/` 实现 gRPC 接口，并在 `internal/app/app.go` 的 `NewApp` 中注册服务。
- 如果需要暴露 HTTP 接口，在 `internal/router/` 中新增路由注册函数，并在 `NewApp` 中挂载。

## 项目配置

### config/config.yaml

```yaml
server:
  port: 8080      # REST API 端口（gRPC 使用 port + 1000）
  mode: debug     # debug 或 release

log:
  level: debug    # debug, info, warn, error
  format: json    # 日志格式
```

## 开发工作流

### 标准开发流程

```bash
# 1. 定义 Proto
vim api/proto/hello.proto

# 2. 生成代码
make proto

# 3. 实现服务
vim internal/service/hello.go

# 4. 注册 HTTP 路由
vim internal/router/hello.go

# 5. 按需调整应用装配（服务/路由挂载）
vim internal/app/app.go

# 6. 生成 Wire 依赖
make wire

# 7. 编译并运行
make run
```

### 快速迭代

```bash
# 修改代码后快速重启
make dev
```

## Wire 依赖注入

Wire 自动生成依赖注入代码。每当修改依赖关系时：

```bash
make wire
```

Wire 会自动分析 `internal/app/wire.go` 中的 `InitializeApp` 函数，并生成 `wire_gen.go`。

## 许可证

MIT
