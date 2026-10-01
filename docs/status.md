# 项目任务状态与记忆

开始工作先读本文件与相关任务的计划、PRD、设计和用例。一个逻辑任务复用同一目录，
后续追问记录在该任务中；变更范围形成独立交付时再建立新任务。
状态可取 `planned`、`active`、`blocked`、`review`、`done`。done 必须有验证证据；
受阻记录原因和恢复动作，不用勾选完成代替验证。

<!-- TASKS:START -->
| 任务 | 状态 | 标题 | 任务记忆与下一步 |
| --- | --- | --- | --- |
| [20260912-docs-task-tests](tasks/20260912-docs-task-tests/00-plan.md) | done | 文档、任务流程与 API 测试规范 | docs 归属、任务流程和测试 skill 已完成；24 项框架/任务及 7 项 API 测试通过，后续任务沿用此流程。 |
| [20260930-structure-fusion-review](tasks/20260930-structure-fusion-review/00-plan.md) | done | 架构融合重构（不含 ES） | 架构及配置收敛已交付；文档生成已完整删除，13 个 Go 测试包、10 项 API、生成/vet/build 和 24 项框架检查通过；删除批次未提交。 |
<!-- TASKS:END -->

## 长期任务记忆

- 2026-09-12：通用框架采用复制注入、显式升级，保留项目定制；首版提交 `69aec4e`。
- 2026-09-12：用户明确 docs 存全部业务信息，tasks 保存每项任务过程，根 tests 存 Python API 测试。AI 配置目录只维护框架与入口引用。
- 2026-09-12：已落地六个任务步骤文档和 UI case 集合，generate-tests 每轮代码修改后执行；本项目当前无前端，不伪造页面验证。
- 2026-10-01：架构融合完成，保留 Proto/鉴权/多协议；不接 ES、不迁入参考业务。文档历史测试曾通过资源注入修复；后续已按用户要求删除该功能。

- 2026-10-01：环境入口收敛为 test/line，默认 test；公共 includes/base.yaml 保留。make ENV 仅接受这两个值，目录配置入口读取 test.yaml。

- 2026-10-01：按用户要求删除 PDF/Word 生成功能，包括 Proto、业务四层、资源/模板、专用配置和 fpdf 依赖；旧接口返回 404，当前保留 auth/dict/hello。
