# 文档、任务流程与 API 测试规范：计划

任务 ID：20260912-docs-task-tests。建立于 2026-09-12，实施前记录。

## 目标与判断

将用户提出的五项目录与开发流程约定落入当前项目，并随通用框架注入新项目。
业务知识集中到 docs；.agents/project 只保留入口和原生设置覆盖，避免双份业务说明。
任务材料先于实现建立。测试生成由可被两种 harness 发现的 skill 执行，主入口明确要求。

采用 Python 标准库编写 API 测试；本机 Go fixture 使用真实 handler、校验和中间件，
以内存 repository 隔离数据库。这样能重复验证 HTTP 行为；MySQL 集成需另行覆盖。
当前没有前端页面，原型与页面操作用例应说明适用范围，不能伪造浏览器验证结果。

## 实施步骤

1. 完成 [PRD](01-prd.md)、[原型](02-prototype.md)、[设计](03-design.md)、[用例](04-test-cases.md)。
2. 将已有开发、架构和问题知识迁到 docs，更新所有入口与可注入模板。
3. 增加任务创建工具，校验任务材料、status 索引和链接；升级保留项目文档。
4. 新增 generate-tests skill，完成 Python API 测试和隔离服务。
5. 使用 skill 根据实际 diff 补充回归用例，执行检查并写入 [验证记录](05-verification.md)。

约束：不修改业务 API 契约，不凭空扩展 UI；已有失败需区分本次引入与历史原因。
