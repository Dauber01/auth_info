# Agent 项目知识框架

将项目变成 Agent 可以检索、验证和持续维护的知识库，同时适配 Codex 与 Claude Code。
本目录可整体作为独立仓库的根目录，不依赖宿主项目语言、Makefile、Git remote 或绝对路径。
运行环境：Python 3.9+ 标准库，以及文件系统符号链接支持。

## 分层与单一来源

| 层 | 位置（注入后的项目） | 维护方式 |
| --- | --- | --- |
| 通用框架 | `.agents/framework/` | 从框架仓库复制，按版本显式升级 |
| 项目知识与任务 | `docs/`、`docs/tasks/`、`docs/status.md` | 项目维护；升级不覆盖已有文件 |
| 项目入口与配置 | `.agents/project/context.md`、`config.json` | 只放简短引用与 AI 设置 |
| API 级测试 | `tests/api/` | 项目使用 Python 维护；升级不覆盖 |
| 项目 skills | `.agents/skills/<name>/` | 项目维护正文及资源 |
| 通用 skills | `.agents/framework/skills/<name>/` | 唯一正文；通过 `.agents/skills/<name>` 软链接选择启用 |
| 项目 agents | `.agents/agents/*.md` | 项目维护，可覆盖或增加通用角色 |
| 原生入口 | `AGENTS.md`、`CLAUDE.md`、`.codex/`、`.claude/`、`.mcp.json` | 生成；不要直接编辑被管理文件 |
| 升级记录 | `.agents/framework.lock.json`、`.agents/generated.lock.json` | 记录文件哈希与链接，提交 Git，不手改 |
| 个人状态 | `.claude/settings.local.json`、用户级 Codex 设置、凭据、索引 | 留在本机，不进入框架仓库 |

`AGENTS.md` 由 `instructions.md` 与项目 `context.md` 组合，自动载入关键规则。
`CLAUDE.md` 导入 AGENTS.md；详细知识从 `docs/README.md` 导航，任务记忆读 `docs/status.md`。
`.claude/skills -> ../.agents/skills`，框架 skills 再通过逐项软链接指向唯一源文件。
只更改选择列表即可启用/停用通用 skill，不需要复制正文。

知识维护方法见 [知识协议](skills/maintain-knowledge/references/protocol.md)。
框架默认启用 `maintain-knowledge`、`generate-tests`，并提供 test-runner、log-analyzer 两个角色。
financial-analyzing 与 bilibili-transcript 是可选通用 skills，不默认注入其执行入口。
默认没有项目权限放行、MCP 或 hooks；目标项目自行选择。

## 注入与显式升级

独立仓库中执行（目标目录须已存在，路径含空格时加引号）：

```bash
python3 harness.py install --target /path/to/project --dry-run
python3 harness.py install --target /path/to/project
cd /path/to/project
python3 .agents/framework/harness.py check --root .
python3 -B -m unittest discover -s .agents/framework/tests
```

抽离之前，将命令中的 `harness.py` 换成当前项目 `.agents/framework/harness.py` 即可。
安装不要求目标项目有 Git、Go 或 Makefile；不联网、不执行目标项目脚本、不安装依赖。
它建立配置入口、docs 文档/任务约定、tests/api 目录，**不会自动理解新项目**。随后让任一 harness 使用
`maintain-knowledge` 读取目标仓库，补齐来源、架构、业务规则与验证命令。
模板的“待验证”是显式知识缺口，不应被当作完成初始化。

0.2.0 起 docs 是业务知识的唯一归属。从 0.1.x 升级时，将原 `.agents/project/` 中的
业务文档迁入 docs 并修正链接，保留 context.md 的简短引用和 config.json；在项目 skills
名单加入 generate-tests。安装器不擅自移动旧业务文档或覆盖项目配置。

## 任务与测试流程

实现前创建/复用任务，并完成计划、PRD、原型、设计流程和初始测试用例：

```bash
python3 .agents/framework/tasks.py new --id 20260912-example --title "示例任务"
python3 .agents/framework/tasks.py check
```

任务模板保存在框架 `task-templates/`，生成到项目 `docs/tasks/<id>/`；状态和任务记忆
统一在 docs/status.md。页面任务使用 ui-cases.json 保存角色、入口、操作步骤与每步断言；
没有页面时必须明确不适用。格式与状态更新命令见注入后的 docs/tasks/README.md。

每轮代码修改后执行 [generate-tests](skills/generate-tests/SKILL.md)，生成/更新测试，
运行并写回验证记录。该要求由共享主入口约束 agent 工作流；同步脚本或文件保存不会
后台启动模型。新项目 tests/api 初始只有约定，不包含假设其技术栈的可执行用例。

升级时在框架独立仓库切换到经过验证的版本，先 dry-run，再执行相同 install 命令。
`VERSION` 是发布版本，lock 保存实际文件内容哈希，因此同版本本地修改也能被发现。
更新框架时先更新版本并跑测试；项目定制尽量放在项目层，避免修改框架副本。
回退也可从旧版本仓库运行 install；已存在的项目知识继续保留，旧版本 schema 必须兼容。

所有路径和冲突先检查，再写入。已有且未被管理的同名文件、手改的生成文件、修改过的
框架文件、普通 `.claude/skills` 目录均会报告冲突，不自动删除或合并。
不要用删除 lock 的方式绕过：先备份并把定制迁移到项目源文件，再恢复被管理文件至
原记录版本，重试。对首次注入已有 AGENTS.md 的项目，先将原规则人工迁入项目知识，
再移走原入口。工具不猜测应保留哪一条规则。

注入前后查看目标项目 Git diff；把 `.claude/settings.local.json`、`__pycache__/`
以及本机索引加入目标项目自己的忽略规则。工具不覆盖目标项目 `.gitignore` 或 CI。
把 check 和 unittest 命令接入该项目 CI 即可，不要求使用本例宿主的 Makefile。

## 项目配置 schema 1

项目 `config.json` 完整初值见 [配置模板](templates/project/config.json)。字段含义：

| 字段 | 合并方式与用途 |
| --- | --- |
| `schema_version` | 必须为 1，不支持的版本拒绝执行 |
| `skills` | 完整的通用 skill 启用名单；项目自写 skills 自动发现 |
| `agents` | 按 name 覆盖通用角色或增加角色；单个 claude/codex 对象整体替换 |
| `disabled_agents` | 停用已定义角色的名称 |
| `mcp_servers` | 按服务名覆盖；目前仅支持 command、字符串数组 args |
| `claude_settings` | 原生 JSON 顶层字段覆盖，生成 `.claude/settings.json` |
| `codex_settings` | model、model_reasoning_effort、sandbox_mode、approval_policy、features |
| `codex_hooks` | 空对象，或原生 `{"hooks": {...}}`，生成 `.codex/hooks.json` |

新项目角色示例：

```json
{
  "name": "contract-reviewer",
  "description": "Review this project's contract changes.",
  "instructions": "agents/contract-reviewer.md",
  "claude": {"tools": ["Read", "Glob", "Grep"], "model": "sonnet"},
  "codex": {}
}
```

正文路径相对于 `.agents/`，必须位于 `agents/` 或 `framework/agents/` 内。
同名默认 agent 的项目覆盖只需包含 name 和需改变的字段。
Claude agent 允许 model、tools、disallowedTools、permissionMode、skills、maxTurns；
Codex agent 允许 model、model_reasoning_effort、sandbox_mode。空 Codex 对象继承父任务设置。
这些字段表达不同原生能力，模型名称和权限不作跨产品翻译。

Hooks 的共享业务脚本放 `.agents/hooks/`，两边分别包装输入输出。
Claude 的 hooks 放 `claude_settings.hooks`；Codex 放 `codex_hooks.hooks`。
配置只按原生结构透传，不证明该事件或 hook 类型可执行，也不会代替 harness 的信任流程。
实际新增 hook 时，须分别测试事件、工具 matcher、输入、退出码和失败行为。
更复杂的 MCP 或原生配置先扩展适配器及测试，不把专有字段混进共享 skill 正文。

## 修改与检查

```bash
python3 .agents/framework/harness.py sync --root .
python3 .agents/framework/harness.py check --root .
python3 -B -m unittest discover -s .agents/framework/tests
```

check 验证源配置、生成结果、软链接、共享 skill frontmatter、docs 链接、任务材料及状态索引。
它不会证明知识内容正确，不校验所有厂商字段枚举，也不会调用模型或真实 hooks。
业务验证由项目知识库指向的测试和代码检查完成。
新知识文档放 docs，在 README 中加入阅读时机、来源与验证状态；新增 skill
必须遵守生成 AGENTS.md 的兼容约定，并在两种 harness 的新会话中检查发现与实际行为。

框架尚在首个项目内开发时，完成核心修改与测试后，运行
`python3 .agents/framework/harness.py install --target .` 将当前核心版本登记到安装锁。
抽离为独立仓库后，在框架仓库开发测试，再对消费项目显式 install；消费项目的 sync
只生成原生配置，不会重写框架安装基线，因此仍能在升级时识别本地定制。

当前适配依据（验证于 2026-09-12，CLI Codex 0.153.4 / Claude Code 2.1.140）：

- [Codex AGENTS.md](https://learn.chatgpt.com/docs/agent-configuration/agents-md)、[skills](https://learn.chatgpt.com/docs/build-skills)、[subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents)、[hooks](https://learn.chatgpt.com/docs/hooks)。
- [Claude memory](https://code.claude.com/docs/en/memory)、[skills](https://code.claude.com/docs/en/skills)、[subagents](https://code.claude.com/docs/en/sub-agents)、[hooks](https://code.claude.com/docs/en/hooks)。

后续独立仓库只需包含本目录内容；不要携带宿主项目 docs、tests、`.agents/project/`、项目 API skill、
个人配置或锁文件。通用框架的详细知识是“怎样采集和维护项目知识”，具体项目知识留在
各自仓库，并随着该项目代码演进。
