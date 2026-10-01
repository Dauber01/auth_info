# auth_info 项目概览

由原根 README 的业务说明迁入。开发约束见 [development.md](development.md)，
任务流程见 [tasks/README.md](tasks/README.md)，验证入口见 [testing.md](testing.md)。

✅ **已实现：**
- Gin Web 框架（REST API）
- Google Wire 依赖注入
- Protocol Buffers + gRPC 服务
- Gin HTTP 与 gRPC 服务共享 Proto 契约
- Viper 配置管理
- Zap 结构化日志
- 错误处理中间件
- 优雅关闭机制
- Makefile 工作流自动化

## 项目结构

```
.
├── cmd/main/                    # 应用入口
├── internal/                    # 应用内部代码
│   ├── app/                     # 应用运行、依赖注入与资源生命周期
│   ├── server/                  # HTTP/gRPC 构建、中间件与路由挂载
│   ├── router/                  # HTTP 路由注册（请求路径按模块拆分）
│   ├── handler/                 # REST API 处理器（请求绑定/转换/响应）
│   ├── validation/              # Proto 参数校验（Protovalidate 封装与错误映射）
│   ├── biz/                     # 业务用例层
│   ├── service/                 # gRPC 服务实现与注册
│   ├── data/                    # 持久化层（模型 + 仓储 + DB）
│   ├── middleware/              # 中间件（JWT/Casbin/统一错误处理）
│   ├── config/                  # 配置管理
│   └── pkg/                     # logger、trace、apperr 公共能力
├── api/                         # API 定义
│   ├── proto/                   # Proto 契约与依赖（统一目录）
│   │   ├── *.proto              # 业务契约（common/auth/dict/document/hello）
│   │   ├── buf/validate/        # Protovalidate 规则定义
│   │   ├── google/protobuf/     # 仓库内维护的 protobuf（如 struct.proto）
│   │   └── third_party/google/  # 第三方 protobuf 依赖
│   └── gen/                     # 生成的 Proto 代码
├── config/                      # 配置文件
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
# 构建并启动（使用已生成代码）
make run

# 开发模式（显式生成并运行）
make dev
```

### 完整命令列表

```bash
# 显示帮助
make help

# 安装固定版本 Go 生成插件（protoc 需预装）
make install-tools

# 生成 Proto 代码
make proto

# 生成 Wire 依赖注入代码
make wire

# 下载/更新依赖
make mod-tidy

# 编译项目
make build

# 运行项目（仅编译 + 启动）
make run

# 开发模式运行（显式生成后启动）
make dev

# 清理构建产物，保留生成源码
make clean

# 运行测试
make test

# 格式化代码
make fmt

# 代码检查
make lint

# 编译与测试
make all
```

## 快速开始

### 1. 安装必需工具

先通过系统包管理器安装 protoc，再安装 Go 插件：

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

在 `config/test.yaml` 或 `config/line.yaml` 中覆盖共享参数：

```yaml
server:
  port: 8080      # REST API 监听端口
  mode: test

log:
  level: warn
  format: json
```

**自动端口分配：**
- REST API HTTP 服务：`http://localhost:8080`
- gRPC 服务：`localhost:9080`（默认 HTTP 端口 + 1000，可用 server.grpc_port 覆盖）

## API 访问

### REST API

```bash
# Hello 接口
curl -H "Authorization: Bearer ${TEST_JWT}" http://localhost:8080/api/v1/hello
```

### gRPC 服务

当前注册的服务为 `hello.HelloService`，契约位于 `api/proto/hello.proto`。
使用具备该契约的 gRPC 客户端调用；当前服务没有注册 reflection，不能依赖服务反射自动列举接口。

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

在 `internal/service/user/` 中创建服务实现文件 `service.go`，实现生成的接口。

### 4. 注册服务与路由

- 在 `internal/service/` 实现 gRPC 接口，并在 `internal/service/service.go` 的注册聚合函数中注册服务。
- 如果需要暴露 HTTP 接口，在 `internal/router/` 中新增路由注册函数，并在 `internal/server/http.go` 中挂载。

## 项目配置

仅保留 test 和 line，默认 test；通过 `make run ENV=line` 选择线上环境。
文件职责、覆盖顺序及部署参数见 [环境配置](configuration.md)。

### test / line

```yaml
server:
  port: 8080      # REST API 端口（gRPC 使用 port + 1000）
  mode: test      # test 环境；line 使用 release

log:
  level: warn     # test 默认 warn；line 默认 info
  format: json    # 日志格式
```

## 开发工作流

### 标准开发流程

先按任务约定完成计划、PRD、原型、设计和初始用例。以下为代码阶段示例；
每轮修改后执行 generate-tests skill，再运行测试并更新任务验证与 status。

```bash
# 1. 定义 Proto
vim api/proto/hello.proto

# 2. 生成代码
make proto

# 3. 实现服务
vim internal/service/hello/service.go

# 4. 注册 HTTP 路由
vim internal/router/hello/router.go

# 5. 按需调整协议装配
vim internal/server/http.go

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

Wire 会自动分析 `internal/app/wire.go` 中的 `initializeApp` 函数（外层 InitializeApp 负责失败回收），并生成 `wire_gen.go`。

## 许可证

MIT

2026-10-01：多环境/includes/APP 配置、生命周期、TraceID、事务和文档资源的完整说明见 [架构](architecture.md)。日志仅本地输出，不接 ES。
