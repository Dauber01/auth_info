# 已知问题与验证记录

范围：解释已有文档生成测试的环境依赖，避免将历史问题归因于无关配置改动。
状态：2026-09-12 执行 `make test` 的历史观察；使用前重新验证，不是永久结论。

当时全量测试中的三个失败：

- `TestGeneratePDF_WithImage`：`template not found: example_template`。
- `TestGenerateWord_WithTextAndImage`：`template not found: word_template_test`。
- `TestGenerateWord_RichTextAndImageOptions`：`template not found: word_template_test`。

来源：[document UseCase](../internal/biz/document/document.go) 的 `NewUseCase`
含 Windows 模板/字体路径；相关测试位于 [document 包](../internal/biz/document/)。
当时 `make fmt` 与 `make lint` 通过；这些测试失败在 AI 配置改造前已复现。

复核命令：`go test -v ./internal/biz/document/...`，共享行为验证运行 `make test`。
测试涉及本机 HTTP 测试服务时需要允许监听本地端口；权限受阻与断言失败须分别报告。
后续修复应为资源路径提供可配置来源和确定性测试夹具，并补充回归验证；本次配置改造
没有修改文档业务代码、模板或断言，也没有通过跳过测试规避问题。

复核触发：文档资源路径、模板加载、图片获取或上述测试改变。修复后更新本记录的状态。
