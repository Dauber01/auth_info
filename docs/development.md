# auth_info 开发规范

来源：原项目 AGENTS.md；2026-09-12 迁入 docs。规范描述预期约束，
实际实现以代码为准。命令来源为 [Makefile](../Makefile)，更改工具链或工作流时复核。

## Commands

```bash
# 首次环境初始化（protoc 需预先通过系统包管理器安装）
make install-tools    # 检查工具，安装 Go 的 protoc 插件
make sync-harness     # 同步两种 harness 的配置和 skills 入口
make check-harness    # 检查配置一致性、skills 资源和同步脚本

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
make test-api         # Python 调用本机隔离服务，验证真实 HTTP API
```

运行后访问：
- HTTP API：`http://localhost:8080`
- gRPC：`localhost:9080`（= HTTP port + 1000）

`make proto` 的工具安装和 `make wire` 可能联网获取 Go 工具；`make build` 还会 tidy、
生成 Proto 与 Wire。执行前确认本机依赖，普通测试使用 `make test` 或对应包的 `go test`。
当前 `install-tools` 在缺少 protoc 时也仅调用 Go 插件安装，不能代替系统安装 protoc。

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

## Go 代码规范

以下规则适用于所有手写 Go 代码。自动生成文件（如 `*.pb.go`、`wire_gen.go`）不适用，且禁止手动修改。

### 格式与行宽

- 遵循 Go 官方惯例、Effective Go 和 Go Code Review Comments，优先采用标准库及仓库已有模式。
- 所有 Go 文件提交前必须经过 `gofmt`；导入分组和排序优先使用 `goimports`。
- 手写代码每行最多 120 个字符。函数调用、参数列表、复合字面量和日志字段超过限制时按语义换行。
- URL、不可拆分的导入路径、结构体标签、生成代码和无法合理断行的错误文本可以例外。
- 不为追求行数限制使用含糊缩写、拆碎日志消息或降低可读性。
- 文件和包保持职责单一；避免循环依赖、无意义的 `utils` 包和过早抽象。

### 命名与可读性

- 包名使用简短的小写单词，不使用下划线、复数堆叠或含糊缩写。
- 导出标识符必须有以该标识符名称开头的 GoDoc 注释。
- 缩写保持一致，例如 `ID`、`HTTP`、`URL`、`JWT`、`API`，不要写成 `Id`、`Http`、`Url`。
- 接口在使用方定义，并保持最小化；除非确实需要替换实现或便于测试，不要预先创建接口。
- 优先使用早返回减少嵌套；复杂函数应拆分为具有明确业务含义的小函数。
- 注释解释“为什么”以及约束，不复述代码正在做什么。
- 禁止遗留调试输出、注释掉的代码、无用途的 TODO 或空实现。

### 错误处理

- 每个错误都必须处理；禁止无说明地丢弃错误。确需忽略时使用 `_ =` 并写明原因。
- 底层错误使用 `%w`、`fmt.Errorf` 或 `apperr.Wrap` 保留错误链，禁止仅拼接 `err.Error()`。
- 业务层使用 `apperr.New` 或 `apperr.Wrap` 返回稳定错误码，不向客户端泄露数据库、文件路径、密钥等内部细节。
- 不要使用 panic 处理可预期错误；panic 仅用于真正不可恢复的程序不变量。
- 错误信息使用小写、无句号的简洁短语，并说明失败动作，例如 `query user: %w`。
- 同一个错误只记录一次：返回错误的底层函数负责补充上下文，最终处理该错误的边界负责记录日志。

### 日志规范

- 使用项目注入的 `zap.Logger`，禁止在业务代码中使用 `fmt.Print*`、标准库 `log` 或直接打印到 stdout。
- 在离错误发生最近、同时拥有足够业务上下文且能决定如何处理错误的边界记录日志。
- HTTP 请求错误统一由 `middleware/error.go` 记录；handler、biz、data 返回带上下文的错误，不重复打印同一错误。
- 启动、后台任务和无法继续向上返回的错误，应在对应命令或任务边界立即记录。
- 日志使用稳定的英文事件名和结构化字段，例如：

```go
logger.Error(
	"create user failed",
	zap.String("username", username),
	zap.Error(err),
)
```

- `Error` 用于请求或任务失败，`Warn` 用于可恢复的异常，`Info` 用于重要状态变化，`Debug` 用于开发诊断。
- 禁止记录密码、JWT、密钥、完整 Authorization 头、数据库连接串以及其他敏感信息。
- 记录可定位问题的字段，如操作名、资源 ID、请求路径、耗时；不要把整份请求对象或大段文本塞入日志。

### Context、并发与资源

- 接收 `context.Context` 的函数将其作为第一个参数，并沿调用链传递。
- handler 以下的代码不得用 `context.Background()` 替换请求上下文，除非明确创建独立生命周期任务。
- 数据库、HTTP 和其他外部调用必须支持超时或取消；循环和长任务应检查 context。
- 启动 goroutine 时必须明确其退出条件、错误传递、资源所有者和关闭方式，避免 goroutine 泄漏。
- 打开的文件、响应体、数据库事务等资源必须在成功获取后立即安排关闭或回滚。
- 共享状态必须通过不可变数据、channel 或明确的同步机制保护；不得依赖未同步的 map 或字段读写。

### 分层与依赖

- handler 只负责协议转换、绑定、校验和响应，不承载业务规则。
- biz 负责业务规则和用例编排，不依赖 Gin、GORM 等传输或持久化实现细节。
- data 负责持久化和外部数据访问，不决定 HTTP 或 gRPC 状态码。
- 新依赖优先通过 Wire 注入，不使用可变全局变量或隐藏的单例。
- 不跨层直接访问实现，例如 handler 不直接操作 GORM，data 不直接构造 HTTP 响应。

### 安全与数据

- 所有外部输入都必须校验，包括 JSON、路径参数、查询参数、文件名、URL 和配置。
- SQL 查询使用 GORM/参数化查询，禁止拼接用户输入生成 SQL。
- 文件路径必须防止路径穿越；远程 URL 获取必须考虑 SSRF、大小限制、超时和内容类型。
- 敏感配置从配置或环境注入，不得写死在代码、测试快照或日志中。
- 涉及认证、授权、密码、Token 或权限策略的改动必须增加失败路径测试。

### 测试与交付

- 开始实现前准备 docs/tasks 的计划、PRD、原型、设计与用例；每轮代码修改后执行 generate-tests skill，更新任务用例、验证记录和 docs/status.md。
- 修复 bug 时优先先添加可复现测试，再实现修复。
- 新业务逻辑至少覆盖成功、参数非法、资源不存在以及关键依赖失败路径。
- 测试应确定性执行，不依赖真实时间、随机网络、共享数据库或执行顺序。
- 修改 `.proto` 后运行 `make proto`；修改 `wire.go` 后运行 `make wire`。
- 提交前至少运行受影响包测试、`make fmt` 和 `make lint`；共享行为或跨模块改动应运行 `make test`。
- 不通过降低断言、跳过测试或吞掉错误来让检查通过。

## 认证与授权

- JWT：HS256 签名，Claims 含 `UserID/Username/Role`，有效期由 `config.yaml` 的 `jwt.expire` 控制
- RBAC：Casbin v3 + GORM 适配器，策略存储在数据库 `casbin_rule` 表；`make seed` 初始化默认策略（admin 全访问，user 访问 GET 路由）
- 中间件顺序固定：ErrorHandler → JWTAuth → CasbinAuth

## 配置

主配置文件：`config/config.yaml`，通过 Viper 加载。
[LoadConfig](../internal/config/config.go) 使用 `AddConfigPath`，所以 `-config` 接收目录。
Makefile 的构建产物为 `bin/auth_info`：

```bash
./bin/auth_info -config ./config
```

生产环境务必修改 `jwt.secret` 和数据库密码。

## 测试

测试文件与被测代码同包，分布在 `internal/biz/` 和 `internal/handler/` 下。运行单个包的测试：

```bash
go test -v ./internal/biz/auth/...
go test -v ./internal/handler/...
```
