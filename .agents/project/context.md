## 当前项目：auth_info

这是 Go 服务项目，提供 Gin HTTP API 与 gRPC，使用 Wire 装配依赖、GORM 持久化。
修改手写 Go 代码前，必须阅读并遵循
[项目开发规范](.agents/project/development.md)；定位模块先读
[架构与业务入口](.agents/project/architecture.md)。项目事实与通用框架分别维护。

- Proto 契约在 `api/proto/`；生成代码在 `api/gen/api/proto/`，Wire 生成文件为 `internal/app/wire_gen.go`。禁止手改生成 Go 文件。
- 修改 Proto 后执行 `make proto`；修改 Wire 声明后执行 `make wire`。
- 常用质量命令为 `make fmt`、`make lint`、`make test`；受影响包按需运行 `go test`。构建、生成和数据库命令的前提及副作用见项目开发规范。
- handler/service 负责协议转换，biz 负责业务，data 负责持久化；受保护 HTTP 路由依次经过 JWT、Casbin。
- API 修改使用项目 `api-conventions` skill，遵循实际 Proto 与响应边界实现。
- 本项目提供 `make sync-harness`、`make check-harness` 作为框架命令入口。
- 全量测试曾有文档模板相关失败，详情与复核命令见 [已知问题](.agents/project/known-issues.md)；这不是跳过测试的理由。
