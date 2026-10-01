# 结构融合重构计划

2026-09-30 用户确认全部架构重构，明确跳过日志 ES；只借鉴架构，不复制参考项目业务代码。
沿用前轮源码分析与同一任务目录。保留 Proto、鉴权、HTTP/gRPC/MCP、业务语义与响应契约。

## 阶段

1. 已完成：server/app/lifecycle 分离；可清理的资源 provider；命令入口统一资源所有权。
2. 已完成：隔离配置实例、环境文件/includes/APP 环境变量、配置校验；公共 pkg、TraceID、日志、超时。
3. 已完成：文档资源注入与取消；事务接口、仓储上下文复用；可重复构建命令。
4. 已完成：按 generate-tests 从 AC 与 diff 补测试；运行单元/API/race/vet/build/生成及框架检查；更新 docs。

阶段可作为一个连贯的代码批次实现；进入验证前必须重读并执行 generate-tests。不引入 Redis/OBS/参考业务或任何远端日志输出。

## 验证

- 保持 API 成功/错误 code 的数值类型、公开/保护边界、文件响应与 Proto 规则。
- 测试 App 端口失败、停止幂等/并发、资源逆序释放与初始化失败清理。
- 测试配置覆盖/隔离/非法输入，TraceID、panic、取消与 MCP 流式例外。
- 事务 rollback/commit/嵌套约束使用隔离数据库；文档使用临时模板与可替换资源。
- make fmt/lint/test/test-api/build/wire、相关 race、harness sync/check 及 24 项框架测试。

## 工具记录

- CodeGraph 已先行查询，未返回的指定文件继续按需读。
- 一次 skill 资源通配符没有匹配，已改为读取确切 SKILL.md，未影响实现。

2026-10-01：以上阶段全部实施；单元/API/race、生成/构建和文档检查通过，详见 [验证记录](05-verification.md)。
