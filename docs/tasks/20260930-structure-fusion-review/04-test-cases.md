# 重构验收用例

2026-10-01 按 PRD AC-01～09 与实际 diff 更新；沿用原有行为断言，不从新实现推导业务预期。
公共清理规则：所有临时文件由 t.TempDir 删除；测试服务器/连接显式关闭；事务库为独立临时 SQLite。
Go 测试不连接部署数据库；Python fixture 用内存账号，runner 退出时清理进程和临时产物。
以下自动化路径均相对仓库根目录，实际执行结果见 [验证记录](05-verification.md)。

| Case ID | AC | 前置、输入与操作 | 预期 | 自动化映射 |
| --- | --- | --- | --- | --- |
| APP-01 | 01/02 | 空闲 App 注入有错误的资源 hook；8 个并发 Stop，再 Run | 逆序释放一次；每次返回相同错误；已停止 App 不重启 | internal/app/app_test.go: TestStopIsReverseConcurrentAndIdempotent |
| APP-02 | 01/02 | 第二监听端口失败；正常双服务启动后 Stop | 保留启动错误，释放首个 listener/资源；正常 Run 退出 | 同文件 TestStartupFailureClosesFirstListenerAndResources / TestRunAndStopWaitForServerCompletion |
| APP-03 | 02 | 构造中成功获得两个资源后 provider 失败 | 逆序释放已取得资源，保留失败原因 | 同文件 TestInitializeFailureReleasesAlreadyAcquiredResources |
| APP-04 | 02/04 | 活动 HTTP 请求；gRPC 请求收到取消但暂不返回；100ms drain | 超时关闭连接，客户端中断；Stop/Run 不因 gRPC handler 忽略取消而无限等待；测试最终释放 handler | 同文件 TestHTTPShutdownDeadlineClosesActiveConnection / TestGRPCShutdownDeadlineCancelsActiveRPC |
| CFG-01 | 03 | 目录/文件入口，多个 includes、入口、APP 覆盖；连续加载两套配置 | 默认值 < includes 顺序 < 入口 < APP；env-only 字段可用；实例互不污染；旧连接池 duration 保持 | internal/config/config_test.go: TestLoadConfig_MySQLPoolDuration / TestLoadConfigPrecedenceIsolationAndEnvironmentOnlyFields / TestIncludeOrderAndEntryOverride |
| CFG-02 | 03 | 缺 include、循环 include、非法端口/负超时/写期限/日志格式/空 secret；显式零超时 | 非法配置在开资源前报错；显式关闭超时不被默认值覆盖 | 同文件 TestLoadConfigIncludesAndValidationFailures |
| HTTP-01 | 04/09 | 带 trace、敏感 query/body/Authorization 的失败请求 | 数值 500、响应头与日志 trace 一致；错误记录一次，不采集请求凭据 | internal/middleware/request_test.go: TestTraceAccessAndFailureEnvelope |
| HTTP-02 | 04/09 | deadline、取消、panic，以及已经写出 202 的超时请求 | 504/408/500；已写响应保留；均不双写 | 同文件 TestTimeoutPanicAndAlreadyWrittenResponses / TestCanceledErrorMapping |
| HTTP-03 | 01/04/09 | 真实 HTTP 装配；公开 auth、保护 hello/dict/document、MCP；记录 writer deadline | 路由边界不变；MCP 能 Flush，无 API context/write deadline；document 写预算更长 | internal/server/server_test.go: TestHTTPRoutingAndMCPDeadlineIsolation |
| RPC-01 | 01/04/09 | bufconn 发真实 Hello RPC 和复用 RegisterRequest 规则的测试 RPC；触发 deadline/panic | Hello 原 code=0 与消息保持；trace metadata 透传；Proto 非法值 InvalidArgument；超时/异常状态正确 | 同文件 TestGRPCHelloValidationAndTrace / TestGRPCBoundaryTimeoutAndPanic |
| LOG-01 | 05 | 两个不同配置 logger；带 trace 输出临时 JSON 文件，多次 cleanup | 实例独立，JSON/console 可构建，文件日志字段正确，stdout 不被当成拥有的文件关闭/同步 | internal/pkg/logger/logger_test.go: TestLoggerInstancesAndLocalFileCleanup |
| TRACE-01 | 04 | 合法/非法/空/过长 ID，不同 context | 安全透传或生成 ID，context 互不污染 | internal/pkg/trace/trace_test.go: TestNormalizeAndContextIsolation |
| DOC-01 | 06/09 | 内存 JSON/DOCX 模板，文字、富文本、图片输入 | 保留原 PDF 与 DOCX 内容/图片断言，取消、缺模板及资源失败向上传递 | internal/biz/document/document_test.go 的原 3 项生成测试；resources_test.go: TestResourcesFailureMissingAndCancellation |
| DOC-02 | 06 | 临时模板目录，正常/缺失模板、路径穿越、向外 symlink、已取消 ctx | 仅读取目录内模板；非法路径、缺失资源、取消可区分 | internal/data/document/resources_test.go: TestLocalResourcesConfineTemplatesAndObserveCancellation |
| DOC-03 | 06 | 本机图片 HTTP：成功、404、取消、超 Content-Length、分块超限；过大 base64 | 合法内容返回；错误与取消传递；10MiB 约束保留，不读取无界远端正文 | 同目录 TestImageSuccessFailureAndCancellation、image_test.go 的两项大小测试；biz/document/fetch_image_test.go |
| TX-01 | 07 | 同一事务调用真实 auth/dict repo，成功或第二阶段失败，事务内外读取 | 两仓储一起提交或一起回滚，读写共用事务 | internal/data/shared/tx_test.go: TestCrossRepositoryCommitAndRollback |
| TX-02 | 07 | 嵌套、跨池仓储、调用前取消、写入后取消 | 拒绝嵌套/跨池；已取消不执行 callback；中途取消回滚 | 同文件 TestTransactionRejectsNestedForeignPoolAndCanceledContext / TestCancellationDuringTransactionRollsBack |
| API-01 | 09 | 独立 fixture 用户注册/登录/hello；重复、非法输入、错误凭据、缺失/错误 token | 复用 7 项既有 API 测试；HTTP/响应 code/data 与鉴权兼容 | tests/api/test_auth.py 原 7 项 test_* |
| API-02 | 04/09 | 自定义 X-Trace-ID 访问保护 API；已登录但无策略访问 dict | 401 数值契约与 trace 保持；无授权返回 403 | 同文件 test_trace_header_preserves_error_contract / test_authenticated_user_without_policy_is_forbidden |
| BUILD-01 | 08 | fmt、Wire、纯 build、vet、全量测试、相关 race；重复 Wire 比较摘要 | 格式正确、生成可重复、编译检查通过、未在 build 中更新模块或生成代码 | Makefile 与验证命令 |
| DOCS-01 | 09 | harness sync/check、tasks check、框架单元测试及相对链接校验 | 文档/入口同步，无坏链接，无生成文件手改 | .agents/framework 与本地链接检查 |

不为现有 Proto/domain/data 转换再写重复用例，原 auth/dict biz、handler、repo 测试纳入全量回归。
未修改 Proto，无需重新生成或修改字段规则；gRPC 校验测试复用已有 RegisterRequest 注解。
页面不适用，ui-cases.json 为 applicable=false；不运行浏览器，也不以 API 结果充当页面结果。
文件滚动算法依赖 lumberjack，本轮验证输出接线及关闭，未宣称长期保留/压缩的压力测试通过。

## 配置收敛新增用例

| Case ID | AC | 前置、输入与操作 | 预期 | 自动化与清理 |
| --- | --- | --- | --- | --- |
| CFG-03 | 10/11 | 临时 config 下存在 test/line/旧 config；空路径、目录、显式文件加载；移除 test 后再加载目录 | 默认与目录选 test，显式 line 正常；缺 test 报错，不回退旧环境 | internal/config/config_test.go: TestLoadConfigDefaultAndDirectorySelectTest；t.TempDir/t.Chdir 自动恢复 |
| CFG-04 | 10/12 | 直接加载仓库 test/line；line 未注入/注入 APP 凭据 | test 独立可加载；line 未配置密钥报错，注入后 release/info，数据库参数来自环境 | 同包 TestRepositoryEnvironmentProfiles；t.Setenv 自动恢复，不连接数据库 |
| CFG-05 | 11/12 | make config 及 run/migrate/seed dry-run，默认/test/line/非法 ENV，自定义目录/文件 | 路径选择一致，非法环境早报错，覆盖有效 | 手动命令断言，无服务启动或数据库写入 |

本轮不修改 API/页面/Proto 行为，复用已有 Go 回归与前轮 API 记录；不增加重复 API 测试。
