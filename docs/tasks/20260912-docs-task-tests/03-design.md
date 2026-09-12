# 设计与流程

docs 保存业务文档和任务事实；.agents/project 保存简短入口引用、config.json。
框架提供 docs 和任务模板；只初始化缺失文档，升级不覆盖已存在内容。
status 表是任务状态的唯一来源，任务文件补充决策、验收细节与证据。

```mermaid
flowchart TD
    A[读取 status 和项目知识] --> B[建立或复用任务]
    B --> C[计划、PRD、原型、设计、初始用例]
    C --> D[实现一轮代码变更]
    D --> E[执行 generate-tests skill]
    E --> F[补充测试并运行]
    F --> G{检查结果}
    G -->|失败| D
    G -->|通过或明确受阻| H[更新验证记录和 status 记忆]
```

API 测试：Python runner 构建/启动 Go fixture，读取随机本机端口，通过 HTTP 发送请求；
fixture 使用真实路由/handler/业务用例，依赖内存用户存储，退出后释放进程与资源。
该测试覆盖 HTTP 到业务边界，不能替代数据库、完整部署或浏览器测试。

兼容性：不添加专有 skill frontmatter，也不通过保存文件的 hook 启动模型；
共享 AGENTS 工作流程要求 agent 在每轮代码修改后实际读取并执行 generate-tests。
