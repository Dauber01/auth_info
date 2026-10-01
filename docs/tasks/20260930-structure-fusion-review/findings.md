# 证据工作记录

- 双方都有 CodeGraph 索引，已使用索引先行查询。
- auth_info 已有 handler/biz/data、Wire、多协议适配和集中 docs 任务流程。
- 参考项目目录增加 server、lifecycle、biz/shared、data/shared、pkg/txctx、trace 与多环境配置，需继续核实实现。
- 参考项目 AGENTS/README 的 Hello-only 依赖示例与当前 brand/asset 目录不完全对应，文档不能直接当作完整事实。

- 已核实参考 server.NewHTTPServer 收拢路由、中间件及 HTTP 超时；auth_info.NewApp 同时承担这些职责并建立 gRPC。
- 参考 App.Stop 逆序执行 lifecycle hooks（MySQL、Redis、日志）；auth_info.Stop 关闭 HTTP/gRPC 并 Sync 日志，但未持有或关闭 DB。
- 参考 LoadConfig 使用独立 Viper 实例、includes 与入口覆盖；auth_info 使用全局 Viper 和固定 config 文件名。
- 参考 TxScope 把业务事务接口与 GORM 实现分开，通过 context 传递事务；适合有跨仓储写入时引入。
- CodeGraph 多次命中相关依赖却未返回所请求函数，已对未覆盖的指定文件/区段使用直接读取，不重复核验已返回源码。

- 已核实错误 code 类型差异、参考超时中间件的协作式语义、参考 Wire 初始化失败缺少 cleanup，以及双方全局 logger 状态。
- 当前 document.NewUseCase 的 Windows 资源路径仍在，GeneratePDF 忽略 ctx。
- 双方各 7 个所选 Go 包测试通过，未运行 API/全量/外部依赖测试。
- 详细判断与目标结构归入 analysis.md；本文件只保留采集进度，不作为另一份长期知识。

## 架构重构进行中（2026-10-01）

- 公共 logger/apperr 已归 pkg；新 trace、server、lifecycle、独立配置、data/document 与 TxScope 已编写。没有迁入参考业务或 ES 输出。
- generate-tests 已按新 AC 与实际 diff 开始执行，既有 document 测试依赖本机模板/写仓库产物，改为隔离 fixture 并保留业务断言。
- task status 工具在 done 时先校验验证记录，因此先保留前轮 passed 将状态恢复 active，再将本轮 verification 改 pending。

## 最终证据（2026-10-01）

以上旧“当前”描述属于分析阶段；现状以 docs/architecture.md 和 05-verification.md 为准。
资源初始化失败通过 bootstrap 统一清理，不依赖 App 已成功构建。
服务 drain 预算不能依赖同步 gRPC Stop 一定返回；忽略取消的 handler 已有复现与回归测试。
事务取消使用临时 SQLite 文件验证回滚，避免取消连接被丢弃造成内存库消失的夹具干扰。
Wire 生成可重复、纯 build 不更改生成文件和模块清单；最终验证完成。

用户后续指定仅保留 test/line；test 原先通过 config.yaml 间接继承本地参数，已归入独立 test.yaml；共享 base 保持无部署凭据。目录入口和三个命令默认统一选 test，不再读取旧 config.yaml。详见 docs/configuration.md。
