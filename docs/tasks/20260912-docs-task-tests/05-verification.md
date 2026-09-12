# 验证与交付

verification_status: passed

验证日期：2026-09-12。范围：docs/任务流程、测试生成 skill、Python API 验证设施；
没有修改业务 API 契约或应用实现。结果对应本任务 AC-01 至 AC-06。

## generate-tests 执行记录

已实际读取 skill，结合 PRD、diff、现有用例和接口源码生成/复核 [用例映射](04-test-cases.md)。
复用原框架 17 项测试，新增 7 项任务流程回归及 7 个 Python API 测试方法（包含参数子案例）。

回归曾发现：状态更新会改动任务表外复制的历史记录；临时包路径中的系统目录软链接
会被错误地当作包内软链接拒绝。对应测试先失败，修复后通过，保留为回归。
没有通过删除断言、跳过失败或只匹配生成文案来完成验证。

## 实际执行结果

| 命令/检查 | 退出码 | 结果 |
| --- | --- | --- |
| `make check-harness` | 0 | 配置、知识链接、任务材料一致；24 项框架/任务测试通过 |
| `make test-api` | 0 | 7 个 API 测试方法通过；实际使用 Python 请求隔离 Go HTTP 服务 |
| `go test ./internal/handler/... ./internal/biz/auth/... ./internal/biz/hello/... ./tests/api/server` | 0 | 有测试的受影响包通过；其余包构建通过，未当作额外测试数量 |
| `make fmt` | 0 | Go 格式化通过 |
| `make lint` | 0 | go vet 通过 |
| skill-creator `quick_validate.py` 与 samefile 检查 | 0 | 五个 skills 格式有效，两种 harness 入口指向同一正文 |
| `git diff --check` | 0 | 无空白错误 |

Go 检查使用临时 GOCACHE 及已缓存依赖，GOPROXY/GOSUMDB 关闭；API fixture 仅监听本机随机端口。
测试账号与密钥仅在隔离服务中使用，任务记录不保存 Token 或真实凭据。

## 覆盖边界与交接

- 无前端页面，原型文件已说明不适用，ui-cases.json 的 applicable 为 false；没有声称运行浏览器测试。
- API fixture 使用真实 auth/hello handler、校验和 JWT/Casbin，内存替代用户数据库；MySQL、gRPC、完整部署及其它业务模块不在本次 API 覆盖内。
- 本轮未重复运行全部 Go 业务测试；既有三个文档模板失败保留在 [已知问题](../../known-issues.md)，不能据本次结果声称全项目测试全部通过。
- 通用框架为 0.2.0，支持注入 docs 和 tests 约定、任务初始化及状态维护；升级保留已有项目文档、任务记忆和测试。
- 测试生成由共享 AGENTS 工作流程要求 agent 执行；没有安装文件保存即启动模型的后台 hook。
- 本次交付包含文档迁移、框架 0.2.0、任务流程与 API 测试设施；后续业务任务使用任务工具建立材料，再进入实现。
