# 架构融合重构验证与交接

verification_status: passed

日期：2026-10-01。已完成用户授权的架构融合，排除 ES 与其他远端日志，不迁入参考项目业务代码。
架构批次验收覆盖 AC-01～09；后附 2026-09-30 分析阶段记录仅作为历史证据。
环境：macOS arm64 / Go 1.25.4 / Python 3；当前工作区未提交。

## 实施结果与测试生成

- app/server 拆分，bootstrap + Lifecycle 统一初始化失败回收、端口失败回收、并发幂等关闭与逆序资源释放。
- 配置独立加载、多环境/includes/APP 覆盖；保留旧目录入口；新增校验与分协议超时。
- pkg/logger/trace/apperr、HTTP/gRPC 错误与 panic 边界；日志只接 stdout/可选本地滚动文件。
- document.Resources 由 data provider 注入，替换硬编码路径；业务检查取消，外部图片请求有期限与大小上限。
- TxScope 通过 context 支持 auth/dict 仓储原子编排，拒绝嵌套/跨池；现有简单 CRUD 保持行为。
- 构建、代码生成、依赖更新分离，Wire 版本由 go.mod tool 锁定；Proto 与业务协议未修改。
- 已按每轮修改执行 generate-tests：核对 AC 与实际 diff，新增生命周期/config/trace/server/资源/事务用例，保留既有 auth/dict 与文档行为断言，更新真实 HTTP fixture 与 2 项 API 回归。
- Case ID、前置/输入/预期/清理与测试函数见 [04-test-cases](04-test-cases.md)。APP、CFG、HTTP、RPC、LOG、TRACE、DOC、TX、API、BUILD、DOCS 各组均通过下述对应检查。
- 无页面，ui-cases.json 标记 applicable=false；没有浏览器结果。

## 最终实际运行

所有 Go 命令使用以下环境（离线依赖与独立缓存）：

```sh
export GOCACHE=/private/tmp/auth-info-structure-review-go-cache
export GOPROXY=off
export GOSUMDB=off
```

| 命令或检查 | 退出码 | 结果 |
| --- | --- | --- |
| `make fmt wire` | 0 | 格式化完成，Wire 正常重新生成，未手改 wire_gen.go |
| `make test` | 0 | 全量 Go 测试通过：15 个有测试的包；无测试的包只编译，不算测试覆盖 |
| `make lint build` | 0 | go vet 与只读模块构建通过，产物 bin/auth_info |
| `make test-api` | 0 | 9 项真实 HTTP API 测试通过，使用独立内存 fixture |
| 下列 `go test -race` | 0 | 9 个有测试的包通过；apperr 无独立测试，映射由边界测试覆盖 |
| 重复 `make wire` 与 `make build` 的 SHA-256 比较 | 0 | wire_gen.go、go.mod、go.sum 内容未改变，生成可重复，build 不改输入 |
| `python3 .agents/framework/harness.py sync --root .` | 0 | 更新 0 个生成文件 |
| `python3 .agents/framework/harness.py check --root .` | 0 | harness 配置与知识链接一致 |
| `python3 -B -m unittest discover -s .agents/framework/tests` | 0 | 24 项框架/任务测试通过 |
| `python3 .agents/framework/tasks.py check` | 0 | 任务材料与状态格式有效 |
| docs 本地 Markdown 相对链接检查 / `git diff --check` | 0 | 链接目标存在，无空白错误；最终检查覆盖完成后的文档 |

```sh
go test -mod=readonly -race -count=1 -timeout=120s \
  ./internal/app ./internal/server ./internal/config ./internal/middleware \
  ./internal/pkg/... ./internal/data/shared ./internal/data/document ./internal/biz/document
```

本机网络测试在允许监听 loopback 的执行环境运行。它们不连接生产或部署数据库。

## 验证中发现并修复的问题

1. 初次日志清理测试发现 stdout 在测试进程中不支持文件式 Sync（bad file descriptor）。stdout 改为不拥有的 Writer，仅关闭应用拥有的滚动文件；重复关闭测试与全量/race 均通过。
2. 新增 APP-04 先复现 gRPC handler 忽略取消时 Stop 超过时限：gRPC 的 GracefulStop 等待 handler 时持有内部锁，同步调用 Stop/等待 Serve 仍可能阻塞。强制 Stop 与服务退出等待改为受 drain 预算约束；保留原失败断言后测试通过。
3. 事务取消测试采用临时文件 SQLite，避免 database/sql 丢弃取消连接时内存库随连接消失；仍断言取消后没有已提交记录。
4. 文档原有 3 项历史失败已修复：用可注入资源与内存 JSON/DOCX 夹具保留内容断言，不再依赖个人路径或写仓库产物。
5. 早期 sandbox 本机绑定失败属于环境限制；在允许本机端口后重跑通过。未将受阻命令算作成功。Wire 早期离线工具依赖问题已通过 go.mod tool 固定工具依赖解决，最终生成/构建验证通过。

## 覆盖边界与维护事项

- 未运行真实 MySQL 故障/连接池集成、数据迁移、seed、完整生产配置启动或部署。跨仓储事务用临时 SQLite 验证，不代表所有 MySQL 方言/锁行为已验证。
- Go 不能终止不响应取消的业务 goroutine。drain 超时后不再等待它们；新业务仍须传递 context 并及时返回。资源 cleanup hook 应能结束，cleanup 本身不在服务 drain 时限内。
- 真实 PDF 中文字体、业务模板视觉效果需用部署资源验收；仓库没有附带中文字体。通过 APP_DOCUMENT_FONT_PATH / document.font_path 配置，未配置时用渲染库内置字体。
- 文件日志接线与关闭已验证，未执行长期滚动/压缩压力测试。无 ES、其他远端日志、Redis 或 OBS 依赖引入。
- 未新增无业务用途的协议 DTO/业务命令；保留现有 Proto/domain/data 转换边界。TxScope 已可用，后续真正的跨仓储用例通过 Wire 注入。
- 参考项目保持只读；未提交 git。实施与验收已完成，后续扩展按 [架构](../../architecture.md) 的模块边界进行。

## 2026-09-30 分析阶段历史证据

日期：2026-09-30。验收范围是源码比较与分析交付，不代表建议已实施或全项目测试通过。
源码基线：auth_info HEAD 5cffad6；参考项目 HEAD 76becec，以当前磁盘文件为准。参考项目原有未提交变更保留。
环境：macOS arm64，go1.25.4。Go 命令使用独立临时构建缓存、禁用模块下载及 sumdb，-mod=readonly 防止修改依赖清单。

## 需求与用例结论

- REVIEW-001 / AC-01：通过。双方先查 CodeGraph，对未返回源码的指定文件补充读取；报告提供实际路径及职责比较。
- REVIEW-002 / AC-02、AC-03：通过。报告区分应吸收/需改造/应保留，明确多协议、认证、错误 code 类型与旧配置入口兼容边界；建议分四阶段。
- REVIEW-003 / AC-04：通过。以下两个限定包测试命令均退出 0，7 + 7 个包通过；未宣称全量通过。
- DOC-001 / AC-04：通过。任务、harness、链接检查及 24 项框架测试通过。
- generate-tests：本轮仅新增/更新分析文档，没有代码修改，不触发每轮代码修改后的 skill 要求；未编写与实现无关的占位测试。
- API：未运行。没有接口变更，本轮不调用真实服务或执行外部写入；参考的旧新对照能力只做源码检查。
- 页面：不适用；无前端页面或页面交互变更，ui-cases.json 明确 applicable=false。

## 实际命令与结果

在参考项目根目录执行（退出 0，7 包通过）：

```sh
GOCACHE=/private/tmp/auth-info-structure-review-go-cache GOPROXY=off GOSUMDB=off go test -mod=readonly -count=1 ./internal/app ./internal/server ./internal/config ./internal/middleware ./internal/pkg/trace ./internal/pkg/txctx ./internal/data/shared
```

涉及的基础设施行为包括逆序关闭一次、服务路由与超时字段、配置入口与默认值、TraceID/错误/recovery、事务 context 与回滚。事务测试使用内存 SQLite，不连接 MySQL。

在 auth_info 根目录执行（退出 0，7 包通过）：

```sh
GOCACHE=/private/tmp/auth-info-structure-review-go-cache GOPROXY=off GOSUMDB=off go test -mod=readonly -count=1 ./internal/config ./internal/data ./internal/biz/auth ./internal/biz/dict ./internal/data/dict ./internal/handler/auth ./internal/handler/dict
```

文档与框架检查（均退出 0）：

```sh
python3 .agents/framework/harness.py sync --root .
python3 .agents/framework/harness.py check --root .
python3 .agents/framework/tasks.py check
python3 -B -m unittest discover -s .agents/framework/tests
git diff --check
```

sync 更新 0 个生成文件；框架单元测试 24 项通过。额外 Python 本地链接检查校验任务 Markdown 与 docs/README.md 的 54 个链接，均存在。

## 限制与交接

- 未运行全量 Go 测试、document 包、Python API、MySQL/Redis/OBS 集成、完整服务启动或生产部署。
- document 的历史模板失败见 docs/known-issues.md；本轮只重新核实硬编码资源路径及忽略 PDF ctx 的源码，没有重跑这些失败测试。
- 首次宽泛 CodeGraph 和一次长文件读取输出被截断；用精确路径/符号补齐需要的证据，没有将工具的测试覆盖提示当作实际覆盖率。
- 报告建议首先拆 server 并处理生命周期，再处理配置与请求追踪；事务和 pkg 归位按真实需求推进。
- 未做代码迁移、依赖升级、数据库变更、提交或参考项目文件修改。后续实现复用此任务材料，细化实际批次验收并执行 generate-tests。

## 配置收敛批次（2026-10-01）

AC-10～12 已实现并验证，以上架构回归记录属于前一批次。本批次基线 a4e449b，当前修改未提交。

- 删除 config/config.yaml、dev.yaml、pre.yaml；原 test 的有效配置归入 test.yaml，仍直接包含公共 base；line 不继承测试配置。
- CLI/config loader 默认 test；目录形式读取 test.yaml。Make 默认 ENV=test，只允许 test/line，CONFIG_FILE/CONFIG_DIR 覆盖保留。
- generate-tests 已执行：新增 TestLoadConfigDefaultAndDirectorySelectTest 与 TestRepositoryEnvironmentProfiles，映射 CFG-03/04；复用配置合并、隔离及校验用例。临时目录/环境自动恢复，不连接数据库。
- CFG-05：25 项实际 Make 配置/命令 dry-run 断言通过，包含默认/test/line、自定义目录/文件、非法与空 ENV；run/migrate/seed 的路径均正确。
- `go run -mod=readonly ./cmd/{main,migrate,seed} -h` 分别执行，均退出 0，帮助信息确认默认 ./config/test.yaml。只运行帮助，不启动或写库。
- 使用本记录的离线 Go 环境执行 `make fmt`、`make test`、`make lint build`，均退出 0；15 个有测试的 Go 包通过。
- harness sync/check、tasks check、24 项框架测试均退出 0；sync 更新 0 文件；134 个 docs 本地链接有效，git diff --check 通过。
- 未更改 API/并发/生成协议，本批次不重复 API/race/Wire；前一批次结果保留，未算作本次新执行。无页面，ui-cases 继续不适用。
- 未运行真实 MySQL 联调、服务部署或数据库迁移；本批次未提交或推送。

## 移除文档生成批次（2026-10-01）

AC-13～15 已完成。基线 7c04791；以上架构/配置/文档渲染结果仅是前批次历史证据。

- 已删除 biz/data/handler/router/document 的完整模块与专用测试、document.proto/document.pb.go、4 个业务模板及专用 known-issues 文档。
- 已清理 HTTPDeps、路由、Wire/provider、文档资源所有权、DocumentConfig/DocumentTimeout、默认值/校验和 base.yaml 对应配置；test/line 及剩余服务不依赖模板目录。
- `make generate` 退出 0，正式生成 Proto/Wire；剩余 Proto 的生成结果仅 protoc 版本注释从原环境变为本机版本，消息/服务代码未变，已按元数据行过滤比较确认。
- `make mod-tidy` 退出 0；只移除 github.com/go-pdf/fpdf 及其两个 sum 条目，没有升级剩余依赖。
- generate-tests 已按实际删除 diff 与 AC 更新 REMOVE-01～04：更新真实装配测试/fixture，新增 `test_document_generation_routes_are_removed`；APP-only 配置覆盖用例改用 log.file.compress，保留原独立加载/校验行为断言。
- 使用前述离线 Go 环境执行 `make fmt`、`make test`、`make lint build`，均退出 0；13 个有测试的 Go 包通过。
- `make test-api` 退出 0，10 项 API 测试通过。新增用例实际验证未登录/已登录时 POST PDF/Word 旧路径均 404，且无附件头；原 9 项鉴权/trace/业务回归通过。
- 源码/配置/模块清单扫描退出 0：无文档模块导入、类型、provider、专用超时/模板/字体/图片配置、生成消息或 fpdf 依赖；业务模板目录不存在。
- harness sync/check、tasks check 均退出 0；sync 更新 0 文件；24 项框架测试通过；123 个 docs 本地链接有效，git diff --check 通过。
- 当前架构、配置、API 和测试说明已更新；旧任务材料仅保留历史观察，已消除指向删除源码的失效链接。无页面，ui-cases 不适用。
- 本轮不重复 race（没有增加并发逻辑），不连接真实 MySQL、执行数据库迁移或部署；当前删除批次未提交/推送。
