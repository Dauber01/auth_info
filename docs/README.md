# auth_info 文档与业务知识

docs 是全部项目业务信息的文档归属。AI 框架在 .agents/framework，项目 AI 配置在
.agents/project/config.json；它们引用这里的事实，不再维护另一份项目知识库。

| 阅读目的 | 文档 | 状态 |
| --- | --- | --- |
| 项目用途、启动与扩展示例 | [项目概览](overview.md) | 由原 README 迁入并修正过时路径 |
| 当前任务、历史决策、交接与下一步 | [任务状态](status.md) | 随每个任务阶段维护 |
| 开始新任务、准备 PRD/原型/流程与用例 | [任务约定](tasks/README.md) | 项目流程要求 |
| Go 开发、生成代码与命令前提 | [开发规范](development.md) | 原规范迁入，事实以源码复核 |
| 协议入口、模块归属与调用边界 | [架构](architecture.md) | 2026-10-01 架构重构后核实 |
| test/line 选择、配置覆盖与字段职责 | [环境配置](configuration.md) | 仅保留两个环境，默认 test |
| 接口契约、校验和错误响应 | [API 约定](api-conventions.md) | 结合 Proto 和 handler 使用 |
| Python API、Go 单元及页面验证 | [测试约定](testing.md) | 记录执行边界与命令 |
| 参考项目结构比较与融合顺序 | [结构融合分析](tasks/20260930-structure-fusion-review/analysis.md) | 保留重构前分析；实施见任务设计与验证 |
| 知识采集范围与缺口 | [采集记录](onboarding.md) | 部分已核实 |

每份业务文档记录来源、适用范围、验证日期或版本及复核触发条件。
业务实现和决策发生变化时同步更新文档、受影响任务用例及 status 记忆。
