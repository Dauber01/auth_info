# 架构重构 PRD

用户决策：完成分析方案中的架构融合，跳过 ES 日志，不迁入参考项目业务代码。

- AC-01：App 只协调运行与关闭；HTTP/gRPC 构建归 server，协议适配继续复用 biz。
- AC-02：资源所有权明确，初始化/启动失败清理，停止有界、逆序、幂等且可并发调用。
- AC-03：文件/旧目录配置入口、多环境/includes、明确覆盖顺序、APP 环境变量、独立加载与校验。
- AC-04：TraceID 贯穿 HTTP/gRPC/业务日志；统一 panic/错误边界；协作式取消，MCP 不套普通请求超时。
- AC-05：控制台 JSON/console 与可选滚动文件日志，消除全局可变 Logger；不接 ES 或其他远端日志。
- AC-06：文档模板/字体/HTTP 资源由配置与 provider 注入，测试不依赖个人绝对路径，支持取消。
- AC-07：biz 定义 TxScope，data 实现，auth/dict 仓储共享事务 context；失败回滚且嵌套语义明确。
- AC-08：内部公共包归位，Proto/domain/data 模型边界不变，构建与依赖更新/生成分离。
- AC-09：原有路径、响应结构、鉴权、Proto 校验和业务行为兼容，单元/API 与文档验证实际通过。

无页面改动，不部署、不迁移数据库、不修改参考项目。

## 配置收敛补充验收（2026-10-01）

- AC-10：仓库环境入口仅 test.yaml 和 line.yaml；test 不再依赖已删除的 config/dev/pre，保留其现有有效设置；公共基础配置可复用。
- AC-11：主服务、migrate、seed 默认 test；LoadConfig 的空路径与目录路径选择 test.yaml；line 必须显式选择；不回退到旧配置或另一个环境。
- AC-12：Make 仅接受 ENV=test/line，支持 CONFIG_DIR/CONFIG_FILE 覆盖；line 的部署配置仍由 APP_* 注入，文档与命令保持一致。

AC-10～12 取代 AC-03 中保留 config.yaml 目录入口的旧决策，其余架构约束保持。

## 删除功能验收（2026-10-01）

- AC-13：移除全部文档生成模块、Proto 及生成文件、模板和专用依赖；无运行时文档资源 provider。
- AC-14：删除 document 配置块、字体/模板/图片参数与 document_timeout，包括默认值/校验；两环境可继续加载。
- AC-15：旧 PDF/Word 路径在未登录和已登录时均 404；剩余 auth/dict/hello/gRPC/MCP 与公共 deadline 行为正确，生成/测试/知识记录同步。

本次明确取代 AC-06 文档生成资源注入需求及 AC-09 中保留文档接口的旧边界；其余业务规则保持。
