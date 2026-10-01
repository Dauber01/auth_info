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

## 2026-10-01 配置收敛追问

用户只保留 test 和 line。复用本任务，删除 dev/pre 与旧 config.yaml 环境入口；
将原测试环境实际依赖配置归入 test.yaml，共享 includes/base.yaml 继续仅承载公共参数。
默认命令/空配置路径/目录路径统一选择 test；line 显式选择，仍通过 APP_* 注入部署凭据。
CONFIG_FILE 保留自定义文件入口，Make ENV 仅接受 test/line。同步当前文档，历史分析保留并标记历史。
本轮只做环境收敛，不连接数据库、不部署；此前授权的提交 a4e449b 已推送，本轮改动不自动推送。
实现后检查实际环境文件加载、默认/显式选择、非法 ENV 与命令 dry-run，再执行全量 Go 测试、fmt/vet/build 和文档检查。

配置收敛已完成：CFG-03～05、全量 Go、fmt/vet/build、CLI 默认入口及框架/链接检查通过。

## 2026-10-01 移除文档生成

用户要求删除文档生成功能及相关模块、模板、文件、配置。基线 7c04791。
先移除 biz/data/handler/router/document、document Proto 与对应生成文件、根 templates 中业务模板；
清理 Wire provider/HTTPDeps/路由与 document 专用 deadline，以及配置类型/默认值/校验/YAML。
删除专用 fpdf 依赖，通过 mod-tidy 确認模块闭包；重新生成 Proto/Wire，更新实际装配测试与 API fixture。
同步当前知识文档，保留必要的任务历史并明确模块已删除，清除历史材料中的失效源码链接。
补测未登录/已登录访问两条旧路由均 404，回归 auth/dict/hello、gRPC、MCP 与 test/line 配置。
本轮执行 fmt/generate/vet/build/全量 Go/API、框架和链接检查，不运行数据库迁移或部署。

文档生成删除批次已完成：正式生成与依赖清理通过，13 个 Go 测试包、10 项 API、vet/build 及框架/链接检查通过。
