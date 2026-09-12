# 测试约定与运行方式

本项目每轮代码修改完成后执行 generate-tests skill，更新任务 Case ID 与自动化映射，
再运行相应测试。测试结果与必要证据写任务 05-verification.md 和 docs/status.md。

## Python API 级测试

```bash
make test-api
# 等价命令：python3 -B tests/run_api.py
```

需要 Python 3.9+、Go 及已下载的项目 Go 依赖，不需要额外 Python 包。
runner 构建临时 Go fixture，在 127.0.0.1 的随机端口运行，再由 Python unittest
通过真实 HTTP 请求验证 API，结束时关闭子进程并删除临时构建产物。

fixture 位于 [tests/api/server](../tests/api/server/main.go)，使用真实 auth/hello 路由、
handler、业务用例、Proto 校验、JWT/Casbin 与错误中间件；用户 repository 使用内存实现。
它不启动项目数据库，不运行迁移或 seed。JWT 测试密钥在启动时临时生成，不写入仓库。

当前范围：注册/登录成功、重复注册、输入非法、错误凭据、缺失/无效 Token，以及
登录后访问受保护 hello 接口。用例位于 [tests/api/test_auth.py](../tests/api/test_auth.py)。
每个用例使用独立用户，整个 fixture 随每次 runner 执行重建。
结果不能代表 MySQL、完整应用启动、gRPC、所有业务模块或浏览器 UI 已验证。

`API_BASE_URL` 由 runner 传给 Python 测试进程，正常使用无需手动设置。
直接执行 unittest 而未提供测试服务地址会失败；空用例集合也会让 runner 失败。
测试只使用 fixture 账号；不得对生产服务执行这些写入场景。

## 单元测试与质量检查

Go 单元测试继续与被测包同目录。按变更运行受影响包，公共行为变化运行 `make test`。
提交前执行 `make fmt`、`make lint`；框架与任务工具执行 `make check-harness`。
既有文档模板问题见 [known-issues.md](known-issues.md)，必须与本次结果分开记录。

## 页面用例

用例及原型保存在 docs/tasks 的对应任务目录，格式见 [任务约定](tasks/README.md)。
API 测试通过不能代替页面点击与逐步断言。没有页面时标记不适用及原因；受阻时记录
缺少的环境与下一步，不编造截图、控件或成功结果。
