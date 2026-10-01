# 架构与业务入口

范围：应用装配、协议、资源与事务边界。2026-10-01 按当前源码核实。
运行验证和历史变化见 [重构任务](tasks/20260930-structure-fusion-review/05-verification.md)。

## 装配与生命周期

- [bootstrap](../internal/app/bootstrap.go) 验证配置、创建 Lifecycle，再调用 Wire 生成的 initializeApp。provider 失败时逆序回收已取得资源，成功后所有权归 App。
- [providers](../internal/app/providers.go) 将部署配置转换为业务所需的 auth.Options，注册 logger、MySQL 的清理函数。
- [App](../internal/app/app.go) 只协调已构建的 HTTP/gRPC：先取得全部监听端口，再启动服务；任一失败关闭监听器及资源。Stop 支持并发/重复调用，所有调用等待相同的关闭结果。
- [Lifecycle](../internal/app/lifecycle.go) 逆序清理并汇总错误；服务先停止，数据库随后关闭，logger 最后关闭。注册发生在初始化阶段，运行时不添加 hook。
- HTTP Shutdown 超时后关闭活动连接；gRPC GracefulStop 超时后强制 Stop。server.shutdown_timeout 为两者共享的优雅关闭预算；业务仍需响应 context 取消。强制关闭不等待忽略取消的 handler；Go 无法强杀这类 goroutine，它们只能在业务返回后结束。清理 hook 不包含在服务 drain 预算中，应保持短且可结束。
- [wire.go](../internal/app/wire.go) 声明构造关系，wire_gen.go 不手改；运行 `make wire` 生成。

## 协议构建与请求链

[server/http.go](../internal/server/http.go) 负责 Gin、中间件、路由、HTTP 参数；
[server/grpc.go](../internal/server/grpc.go) 负责 gRPC 注册、trace、panic/错误边界、期限与 Proto 校验。
HTTP port 与 grpc_port 独立配置，grpc_port 省略或为 0 时回退到 HTTP port + 1000。

```text
HTTP → trace → access → recovery → error → 路由组/鉴权/timeout → handler → biz → data
MCP  → 同一 HTTP 全局中间件 → mcpserver → biz
gRPC → trace/错误/期限 → Protovalidate → service → biz
```

- `/api/v1/auth` 公开注册；hello、dict 使用 JWT → Casbin 保护组。
- `/mcp` 保留独立注册，未加入 JWT/Casbin，也不套用普通 API 请求/写响应期限。改变此边界须单独评审契约。
- gRPC 当前只注册 HelloService；保留已有 Proto，Hello name 本身没有长度校验规则。
- service 是 gRPC 协议适配层，业务编排在 biz。handler 不访问 ORM；业务接口和模型不依赖 Gin/GORM/部署配置。
- HTTP 错误仍返回数值 code；TraceID 使用 X-Trace-ID 响应头。gRPC 使用 x-trace-id metadata。
- HTTP request_timeout、grpc_timeout 是协作式取消，不后台执行或强杀 handler。WriteTimeout 仅给普通 API 设置连接写期限；MCP 不设置短写期限。HTTP ReadHeader/Read/IdleTimeout 仍约束连接。

## 业务模块

| 模块 | 业务边界 | 实现/适配 |
| --- | --- | --- |
| auth | [业务/用户接口](../internal/biz/auth/)；仅接收 JWT Options | [data/auth](../internal/data/auth/)、[HTTP](../internal/handler/auth/) |
| dict | [业务/仓储接口/模型](../internal/biz/dict/) | [data/dict](../internal/data/dict/)、[HTTP](../internal/handler/dict/) |
| hello | [共享业务](../internal/biz/hello/) | HTTP、[gRPC](../internal/service/hello/)、[MCP](../internal/mcpserver/hello.go) |

## 事务

[biz/shared.TxScope](../internal/biz/shared/tx.go) 只定义 `Do(ctx, fn)`；
[data/shared](../internal/data/shared/tx.go) 用 GORM 实现，auth/dict repo 统一通过 `shared.DB(ctx, db)` 选连接。
同一个 scope 中可以原子编排两种仓储，回调失败或 context 取消回滚。禁止嵌套 scope、跨连接池混用；回调 ctx 不得逃逸。
现有单仓储业务保持原行为。新增需要事务的用例再通过 Wire 注入 NewTxScope，不给所有 CRUD 强加外层事务。

## 公共能力、配置与工具

- [pkg/logger](../internal/pkg/logger/) 无全局单例，输出 stdout 与可选本地滚动文件；无 ES/远端上传。
- [pkg/trace](../internal/pkg/trace/) 只处理 context/安全 ID；[pkg/apperr](../internal/pkg/apperr/) 保留 HTTP/gRPC 错误映射。
- 日志记录路由、方法、耗时、状态和 trace，不采集 body/query/Authorization；错误详情由错误边界记录一次。
- [config](../internal/config/) 使用独立 Viper，优先级：默认值 < includes 顺序合并 < 入口文件 < APP_*。LoadConfig 接受具体文件或目录；空路径默认 test 文件，目录读取 test.yaml。环境入口仅 test/line，详情见 [环境配置](configuration.md)。
- `make build` 只编译；`make generate` 显式生成；`make mod-tidy` 显式更新依赖；Wire 版本由 go.mod/tool 锁定。
- 数据库迁移、权限初始化仍是 [migrate](../cmd/migrate/main.go)、[seed](../cmd/seed/main.go) 显式操作，命令也负责关闭连接与日志。

文档生成功能已于 2026-10-01 按用户要求移除，当前业务仅 auth、dict、hello。

复核触发：provider、服务注册、中间件、配置、事务、日志输出或资源访问边界改变。
