# auth_info 项目知识索引

本目录由项目维护。通用框架升级不覆盖这些知识；变更相关代码时同步复核来源。

| 任务 | 阅读入口 | 证据与状态 |
| --- | --- | --- |
| 开发 Go、生成代码、运行检查 | [开发规范](development.md) | 原项目 AGENTS.md 规范迁入；命令来源 Makefile |
| 定位 HTTP/gRPC/MCP 与模块依赖 | [架构与业务入口](architecture.md) | 已阅读启动装配及服务注册，2026-09-12 |
| 修改接口契约、校验和响应 | [API skill](../skills/api-conventions/SKILL.md) | 项目专用约定；实现以 Proto、httpx 和中间件为准 |
| 分析文档测试失败 | [已知问题](known-issues.md) | 有日期的历史验证结果，使用前重新确认 |
| 修改 Agent 框架或注入其他项目 | [框架说明](../framework/README.md) | 可独立抽离的通用资产与配置格式 |
| 增补知识、处理文档与代码不一致 | [知识维护 skill](../framework/skills/maintain-knowledge/SKILL.md) | 通用维护流程 |

尚未完成逐条验证的内容：各业务模块完整状态机、生产部署及回滚流程、所有鉴权失败路径。
需要这些信息时，从相关源码、测试和用户决策补充；当前索引不代表已全面审计业务行为。
