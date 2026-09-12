# 架构与业务入口

范围：定位运行入口、协议边界和模块归属。状态：下列入口已于 2026-09-12 阅读确认；
分层约束属于项目开发规则，不代表所有历史实现均已满足。

## 启动与协议

- [应用装配](../internal/app/app.go) 的 `NewApp` 组装 Gin、错误处理中间件、路由和 gRPC。
- HTTP 的 `/api/v1` 下认证路由在公开组注册；hello、dict、document 在 JWT 与 Casbin 保护组注册。具体方法和子路径以 [router](../internal/router/) 的各模块注册函数为准。
- `/mcp` 在该保护组外独立注册，使用注入的 Hello MCP handler。HTTP 业务路由的 JWT/Casbin 保护不能推导为 MCP 也有相同保护。
- gRPC 安装 Protovalidate unary interceptor；[service.go](../internal/service/service.go) 当前仅注册 HelloService。其他模块是否提供 gRPC 要检查实际注册，不能只看 Proto 定义。
- `Run` 使用配置 HTTP 端口，gRPC 端口为 HTTP 端口加 1000；`Stop` 处理服务关闭。
- 依赖声明在 [wire.go](../internal/app/wire.go)，生成的 wire_gen.go 不手改。

## 层与模块

HTTP 调用方向为 router → handler → biz → data；gRPC 从 service 进入 biz。
协议转换、业务规则、存储访问分别由这些边界承担，详细规范见 [development.md](development.md)。

| 模块 | 导航 | 进一步核实的内容 |
| --- | --- | --- |
| auth | [业务](../internal/biz/auth/)、[HTTP 路由](../internal/router/auth/) | 登录、注册、JWT 与用户数据约束 |
| dict | [业务](../internal/biz/dict/)、[HTTP 路由](../internal/router/dict/) | 字典类型/项操作与唯一性规则 |
| document | [业务](../internal/biz/document/)、[HTTP 路由](../internal/router/document/) | PDF/Word 模板、图片访问及资源前提 |
| hello | [业务](../internal/biz/hello/)、[gRPC](../internal/service/hello/) | 示例协议接口与复用逻辑 |

接口契约来源为 [api/proto](../api/proto/)；错误映射来源为
[apperr](../internal/apperr/apperr.go)、[HTTP 错误中间件](../internal/middleware/error.go)。
避免在知识文档复制整份契约；变更接口时更新对应 skill 与规则引用。

数据库迁移入口为 [cmd/migrate](../cmd/migrate/main.go)，策略初始化入口为
[cmd/seed](../cmd/seed/main.go)。这些操作修改数据库，不是一般验证命令。

复核触发：路由注册、服务注册、中间件顺序、Wire 依赖、端口或模块边界改变。
