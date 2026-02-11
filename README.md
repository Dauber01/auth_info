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
├── cmd/main/              # 应用入口
├── internal/              # 应用内部代码
│   ├── app/              # 应用层（Gin + gRPC）
│   ├── config/           # 配置管理
│   ├── handler/          # REST API 处理器
│   ├── service/          # gRPC 服务实现
│   ├── logger/           # 日志组件
│   └── middleware/       # 中间件
├── api/                  # API 定义
│   ├── proto/           # Proto 文件定义
│   └── gen/             # 生成的 Proto 代码
├── pkg/                 # 公共包
├── config/              # 配置文件
├── docs/                # Swagger 文档
├── Makefile             # 构建脚本
└── go.mod              # Go 模块定义
```

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

在 `api/proto/` 目录下创建新的 `.proto` 文件，例如 `api/proto/user.proto`：

```protobuf
syntax = "proto3";

package user;

option go_package = "auth_info/api/gen/user;user";

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}

service UserService {
  rpc GetUser(UserRequest) returns (User);
}
```

### 2. 生成代码

```bash
make proto
```

这将在 `api/gen/` 目录下生成对应的 Go 代码。

### 3. 实现服务

在 `internal/service/` 中创建服务实现文件 `user.go`，实现生成的接口。

### 4. 注册服务

在 `internal/app/app.go` 中的 `NewApp` 函数中注册新的 gRPC 服务。

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

# 4. 注册服务
vim internal/app/app.go

# 5. 生成 Wire 依赖
make wire

# 6. 编译并运行
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

## 常见问题

### Q: 如何在 Windows 上安装 protoc？
A: 访问 https://github.com/protocolbuffers/protobuf/releases 下载 protoc-*-win64.zip，解压后将 bin 目录添加到 PATH 环境变量中。

### Q: 如何添加 gRPC 网关映射？
A: 在 Proto 定义中使用 google.api.http 注解，然后使用 grpc-gateway 代码生成器。参考官方文档：https://grpc-ecosystem.github.io/grpc-gateway/

### Q: 生成的代码为什么没有提交到 Git？
A: `api/gen/` 目录下的生成文件已被添加到 `.gitignore` 中，应该从 Proto 源文件重新生成，避免冲突。

### Q: 如何改变 gRPC 和 HTTP 的端口？
A: 修改 `config/config.yaml` 中的 `server.port` 字段。gRPC 将自动使用 `port + 1000` 作为端口。

