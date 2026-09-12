# 任务目录与执行约定

每个逻辑任务一个 `YYYYMMDD-slug` 目录。实现前先分析目标、约束、影响和方案，
写出可复核的计划与设计理由，然后完善 PRD、原型、流程和初始测试集合。
这一步是分析和文档产出，不自动增加用户确认关卡。

| 文件 | 内容 |
| --- | --- |
| 00-plan.md | 问题、约束、方案判断、实施步骤与验证方式 |
| 01-prd.md | 用户需求、业务规则、范围、AC 编号验收标准 |
| 02-prototype.md | 可查看原型及状态说明；图片/HTML 放任务 assets；无 UI 明确不适用 |
| 03-design.md | 组件与数据职责、接口/状态、Mermaid 或图片设计流程、失败分支 |
| 04-test-cases.md | 需求到 Case ID 的映射，包含单元、API、页面场景与预期 |
| ui-cases.json | AI 操作页面的结构化步骤集合；不得用 API 测试结果替代 |
| 05-verification.md | 测试生成结论、实际运行证据、问题、交接记忆 |

```bash
python3 .agents/framework/tasks.py new --id 20260912-example --title "示例任务"
python3 .agents/framework/tasks.py status --id 20260912-example --state active --memory "需求和设计已完成，开始实现"
python3 .agents/framework/tasks.py check
```

命令从项目根目录执行，也可在子命令前传 `--root /path/to/project`。
新建任务是 planned；实现前完善材料，确认 UI 是否适用，再切为 active。
恢复旧任务时先读 docs/status.md 和任务记录；追加需求复用目录，不为每条消息重复建任务。
阶段完成/暂停时更新状态与记忆；done 需要 `verification_status: passed` 和实际验证证据。
检查器验证文件结构与记录格式，不能证明文档已充分思考或测试真的通过。

## 页面用例格式

初建时 applicable 为 null，表示尚未分析。开始实现前必须选择 true 或 false。
false 要给出原因且 cases 为空；true 至少包含一个完整场景。示意格式如下，实际任务
必须换成 PRD/原型中的真实页面与断言：

```json
{
  "schema_version": 1,
  "applicable": true,
  "reason": "该任务改变登录页面交互",
  "cases": [{
    "id": "UI-001",
    "title": "空表单显示校验提示",
    "role": "未登录访客",
    "start_url": "${WEB_BASE_URL}/login",
    "preconditions": ["使用隔离浏览器会话，页面已加载"],
    "steps": [{"action": "click", "target": "可访问名称为登录的提交按钮", "expected": "用户名和密码字段显示必填提示，仍在登录页"}],
    "cleanup": ["关闭该隔离会话"]
  }]
}
```

执行前由 agent 从当前测试环境配置解析 WEB_BASE_URL；它不是 harness 自动变量。
定位优先使用可访问名称、标签和稳定测试标识，不猜测页面坐标或不存在的控件。
每步校验可观察结果，失败记录实际状态；必要截图放 assets 并在验证记录链接。
