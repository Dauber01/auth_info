# Agent 开发入口

Codex 自动读取本文件；Claude Code 通过 `CLAUDE.md` 的 `@AGENTS.md` 导入同一份说明。
本文件由框架规范与项目上下文生成。修改源文件后同步，不直接编辑生成结果。

## 知识库使用与维护

- 开始任务时先阅读 [.agents/project/index.md](.agents/project/index.md)，再按任务读取对应知识文档。
- 项目事实、开发约束、验证结果分别记录，事实附代码路径或命令来源；未核实的内容明确标记“待验证”。知识与代码不一致时核实并更新，不把旧文档当作实现事实。
- 保留项目已有术语和架构边界。跨模块修改、接口变更、命令变化和问题修复后，更新受影响知识及索引；不要把原始会话、密钥或整份源码复制进知识库。
- 知识不足时使用 `maintain-knowledge` skill 采集与任务相关的证据。没有代码证据时不臆造业务规则，也不为填满模板而扩大任务范围。
- 如果仓库根目录存在 `.codegraph/`，理解或定位代码时先调用 `codegraph_explore` 或 `codegraph explore`；没有索引则跳过，不自行创建索引。

## 配置和资源归属

- `.agents/framework/` 是可独立抽离的通用框架：规范、共享 agent、通用 skill、适配器和测试。项目知识与覆盖配置放在 `.agents/project/`；项目自定义 agent 放在 `.agents/agents/`。
- **同一 skill 只维护一份正文和资源。** 项目 skill 创建于 `.agents/skills/<name>/`；框架 skill 的正文在 `.agents/framework/skills/<name>/`，由 `.agents/skills/<name>` 软链接暴露。`.claude/skills -> ../.agents/skills` 让两种 harness 读取同一文件。不要复制或替换成普通目录。
- 新 skill 的共享 frontmatter 只使用单行 `name`、`description`；name 与目录名一致，使用小写字母、数字和连字符。正文使用 Markdown，以相对于 `SKILL.md` 的链接引用真实资源。
- 不在共享 skill 中使用 `context: fork`、`agent`、`hooks`、`allowed-tools`、`disable-model-invocation`、`$ARGUMENTS`、动态命令展开或 `${CLAUDE_SKILL_DIR}` 等专有语义。描述工具能力，不写死 harness 工具名。脚本路径从已加载的 `SKILL.md` 解析，不依赖当前工作目录。
- Agent 的正文共用 Markdown；模型、工具和权限按 harness 分别配置。Claude agent 生成 Markdown/YAML，Codex agent 生成 TOML，不能直接共用原生配置文件。
- `.agents/project/config.json` 选择框架 skills、agent 覆盖、MCP 和各 harness 设置；通用默认值在框架 `defaults.json`。不要手动编辑生成的 `.codex/`、`.claude/agents/`、`.claude/settings.json`、`.mcp.json`。
- MCP 只共享已验证的 `command` / `args` 子集。使用 PATH 中的程序，不写个人绝对路径或凭据。注入不会安装 MCP、创建索引或改变用户全局配置。
- Hooks 按两套原生格式分别配置，业务逻辑可共用项目 `.agents/hooks/` 中的脚本；事件名、matcher、输入和退出码语义必须分别适配测试。框架不自动转换 hooks，也不默认启用。Claude 的 `prompt` / `agent` hook 不能假定在 Codex 执行。
- `.claude/settings.local.json` 和用户级 Codex 配置保存个人设置；不要提交或通过框架升级覆盖。项目权限覆盖必须显式维护，不从当前开发机历史权限推导通用默认值。

## 同步与交付

从项目根目录执行（Python 3.9+ 标准库，无包安装）：

```bash
python3 .agents/framework/harness.py sync --root .
python3 .agents/framework/harness.py check --root .
python3 -B -m unittest discover -s .agents/framework/tests
```

生成文件发生手工修改时，同步会报告冲突。先把定制移到项目源配置，再恢复生成文件到已记录版本后同步。
升级仅更新框架管理的内容，项目知识和个人配置保留。两种 CLI 均从项目根目录启动；更改发现配置后新建会话验证。
Windows 必须启用符号链接支持并让 Git 保留软链接，不能用复制 skills 代替。
验证命令以项目知识库为准；汇报实际执行与失败原因，不能将受阻或未执行的检查说成通过。
