# auth_info

项目业务说明统一维护在 [docs](docs/README.md)。开始或恢复工作先阅读
[任务状态与记忆](docs/status.md)，再进入对应任务材料。

- [项目概览](docs/overview.md)
- [架构与业务入口](docs/architecture.md)
- [开发规范](docs/development.md)
- [任务材料与流程](docs/tasks/README.md)
- [测试运行方式](docs/testing.md)
- [通用 Agent 框架与注入说明](.agents/framework/README.md)

Codex 与 Claude Code 共用生成的 [AGENTS.md](AGENTS.md)。Skills 保持单一来源，
项目配置在 `.agents/project/config.json`，个人配置不提交。

```bash
make sync-harness
make check-harness
make test-api
```

框架修改、注入其他项目与升级方法见框架说明。具体业务命令和前提以 docs 为准。

## 许可证

MIT
