# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# 首次环境初始化
make install-tools    # 安装 protoc、protoc-gen-go、protoc-gen-go-grpc

# 代码生成（修改 .proto 文件或 wire.go 后必须执行）
make proto            # 从 .proto 文件生成 Go 代码
make wire             # 生成 Wire 依赖注入代码（wire_gen.go）

# 构建与运行
make run              # 清理 → proto → wire → 构建 → 运行（首次或完整重建）
make dev              # 快速重建（跳过 clean，适合开发迭代）

# 数据库
make migrate          # 执行数据库迁移（cmd/migrate）
make seed             # 初始化 Casbin 默认策略（cmd/seed，首次部署必须运行）

# 质量检查
make fmt              # go fmt 格式化
make lint             # go vet 静态分析
make test             # go test -v ./...（运行全部测试）
make docs             # 生成 Swagger 文档
```

运行后访问：
- HTTP API：`http://localhost:8080`
- gRPC：`localhost:9080`（= HTTP port + 1000）
- Swagger UI：`http://localhost:8080/swagger/index.html`

## 架构概览

本项目是一个同时提供 **REST API（Gin）** 和 **gRPC** 服务的 Go 微服务框架，采用清洁架构分层：

```
HTTP Request
    → middleware/error.go (全局错误捕获)
    → middleware/auth.go (JWT 验证 → Casbin RBAC 鉴权)
    → handler/ (请求绑定、参数校验)
    → biz/ (业务逻辑、用例)
    → data/ (数据库操作、GORM)

gRPC Request
    → validation/grpc_interceptor.go (protovalidate 校验)
    → service/ (实现 gRPC 接口)
    → biz/ (同上业务逻辑层)
```

## 依赖注入（Google Wire）

所有依赖关系在 `internal/app/wire.go` 中声明，`wire_gen.go` 是自动生成的，**不要手动编辑**。修改 `wire.go` 后必须运行 `make wire` 重新生成。

## Proto 文件与代码生成

- Proto 定义位于 `api/proto/`，生成代码位于 `api/gen/api/proto/`（自动生成，不要手动编辑）
- 字段校验使用 `buf.build/validate` 注解（`protovalidate`），不是 `protoc-gen-validate`
- 新增业务模块需要：① 写 `.proto` → ② `make proto` → ③ 实现 handler/biz/data → ④ `make wire`

## 错误处理规范

使用 `internal/apperr/apperr.go` 中定义的错误码，该包负责将应用错误映射到 HTTP 状态码和 gRPC 状态码：

```go
apperr.New(apperr.CodeNotFound, "user not found")
apperr.Wrap(apperr.CodeInternal, "db query failed", err)
```

`middleware/error.go` 全局拦截 panic 和错误，统一转换响应格式。

## 认证与授权

- JWT：HS256 签名，Claims 含 `UserID/Username/Role`，有效期由 `config.yaml` 的 `jwt.expire` 控制
- RBAC：Casbin v3 + GORM 适配器，策略存储在数据库 `casbin_rule` 表；`make seed` 初始化默认策略（admin 全访问，user 访问 GET 路由）
- 中间件顺序固定：ErrorHandler → JWTAuth → CasbinAuth

## 配置

主配置文件：`config/config.yaml`，通过 Viper 加载。启动时用 `-config` 标志指定路径：

```bash
./bin/server -config ./config/config.yaml
```

生产环境务必修改 `jwt.secret` 和数据库密码。

## 测试

测试文件与被测代码同包，分布在 `internal/biz/` 和 `internal/handler/` 下。运行单个包的测试：

```bash
go test -v ./internal/biz/auth/...
go test -v ./internal/handler/...
```
