# 已知问题与验证记录

范围：解释已有文档生成测试的环境依赖，避免将历史问题归因于无关配置改动。
状态：2026-09-12 历史失败；2026-10-01 已通过资源注入与隔离夹具修复，document 包测试已重跑通过。

当时全量测试中的三个失败：

- `TestGeneratePDF_WithImage`：`template not found: example_template`。
- `TestGenerateWord_WithTextAndImage`：`template not found: word_template_test`。
- `TestGenerateWord_RichTextAndImageOptions`：`template not found: word_template_test`。

历史来源：[document UseCase](../internal/biz/document/document.go) 的 `NewUseCase`
含 Windows 模板/字体路径；相关测试位于 [document 包](../internal/biz/document/)。
当时 `make fmt` 与 `make lint` 通过；这些测试失败在 AI 配置改造前已复现。

复核命令：`go test -v ./internal/biz/document/...`，共享行为验证运行 `make test`。
测试涉及本机 HTTP 测试服务时需要允许监听本地端口；权限受阻与断言失败须分别报告。
2026-09-12 的配置改造没有修改文档业务代码、模板或断言，也没有通过跳过测试规避问题。
当时建议资源注入与确定性夹具，已在下述 2026-10-01 重构中实施。

复核触发：文档资源路径、模板加载、图片获取或上述测试改变。修复后更新本记录的状态。

## 2026-10-01 修复

NewUseCase 接收 Resources；data/document 通过配置打开模板目录、可选字体和 HTTP client，已移除开发机 Windows 路径。
上述三项测试现在创建独立夹具并保留原文本/图片/富文本断言，输出不再写 templates/test*.docx。
本轮 `go test -count=1 ./...` 中 document 包已通过，完整最终验证见 [架构任务](tasks/20260930-structure-fusion-review/05-verification.md)。
部署仍需提供实际业务模板与中文字体；单元夹具不等同于真实模板视觉验收。
