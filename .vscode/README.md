# VS Code 配置说明

本目录包含了项目的 VS Code 开发环境配置。

## 文件说明

### 1. settings.json
工作区设置，包含：
- **Go 测试配置**：默认使用 `-v` 参数显示详细输出
- **代码格式化**：保存时自动格式化，使用 gofmt
- **代码提示**：启用自动导入未引用的包
- **文件排除**：排除 bin、vendor、生成的 proto 文件等
- **多语言配置**：Go、Proto、YAML、JSON、Makefile 的编辑器设置

### 2. launch.json
调试配置，包含：
- **启动主程序**：调试主应用程序
- **调试当前文件**：调试当前打开的 Go 文件
- **调试当前测试**：调试当前测试文件
- **调试单个测试函数**：调试选中的测试函数
- **附加到进程**：附加到正在运行的进程
- **调试 Wire 生成**：调试 Wire 代码生成过程

### 3. tasks.json
任务配置，包含：
- **生成 Proto 代码** (Ctrl+Shift+B → 选择)
- **生成 Wire 代码**
- **编译项目** (默认构建任务)
- **运行项目**
- **清理构建文件**
- **运行测试 (verbose)**
- **运行测试并生成覆盖率报告**
- **查看测试覆盖率**
- **格式化代码**
- **代码检查 (vet)**
- **更新依赖**
- **完整构建 (all)**

### 4. extensions.json
推荐扩展列表：
- `golang.go` - Go 语言支持
- `zxh404.vscode-proto3` - Protocol Buffers 支持
- `humao.rest-client` - REST API 测试
- `eamodio.gitlens` - Git 增强
- 其他辅助工具

### 5. go.code-snippets
Go 代码片段：
- `handler` - 创建 Gin Handler（带 Swagger 注释）
- `service` - 创建 Service 函数
- `test` - 创建表驱动测试
- `iferr` - 错误检查
- `grpcservice` - 创建 gRPC 服务实现
- `log` - 添加日志
- `struct` - 创建带 JSON 标签的结构体
- `wire` - 创建 Wire Provider
- `protomessage` - 创建 Proto 消息
- `protoservice` - 创建 Proto 服务

### 6. api-test.http
REST Client 测试文件：
- 使用 REST Client 扩展快速测试 API
- 包含常见的 HTTP 请求示例

## 快速开始

### 1. 安装推荐扩展
打开命令面板 (Ctrl+Shift+P)，输入：
```
Extensions: Show Recommended Extensions
```
点击 "Install All" 安装所有推荐扩展。

### 2. 运行构建任务
按 `Ctrl+Shift+B` 打开构建任务列表，选择要执行的任务。

默认构建任务是 "编译项目"，直接按 `Ctrl+Shift+B` 即可快速编译。

### 3. 启动调试
按 `F5` 或点击侧边栏的"运行和调试"按钮，选择调试配置：
- **启动主程序**：最常用，调试完整应用
- **调试当前测试**：在测试文件中使用

### 4. 运行测试
在测试文件中，函数名上方会显示 "run test" 和 "debug test" 按钮。
点击即可快速运行或调试单个测试。

测试输出会显示详细信息（因为配置了 `-v` 参数）。

### 5. 使用代码片段
在 Go 文件中输入片段前缀（如 `handler`），然后按 `Tab` 键：
- 输入 `handler` → 生成完整的 Gin Handler 函数
- 输入 `test` → 生成表驱动测试模板
- 输入 `iferr` → 生成错误检查代码

### 6. 测试 API
打开 `api-test.http` 文件，点击请求上方的 "Send Request" 按钮即可测试 API。

## 常用快捷键

### 调试相关
- `F5` - 启动调试
- `Shift+F5` - 停止调试
- `F9` - 切换断点
- `F10` - 单步跳过
- `F11` - 单步进入
- `Shift+F11` - 单步跳出

### 构建相关
- `Ctrl+Shift+B` - 运行构建任务
- `Ctrl+Shift+P` → 输入 "task" - 运行任何任务

### 测试相关
在测试文件中：
- 函数上方出现 "run test" 链接
- 点击运行/调试单个测试
- 测试结果显示在输出面板

### 代码导航
- `F12` - 跳转到定义
- `Alt+F12` - 查看定义
- `Shift+F12` - 查看所有引用
- `Ctrl+T` - 转到符号
- `Ctrl+Shift+O` - 转到文件中的符号

### 代码编辑
- `Alt+Shift+F` - 格式化文档
- `Ctrl+.` - 快速修复
- `F2` - 重命名符号

## 自定义配置

如果需要个人化配置，可以：
1. 在用户设置中覆盖工作区设置
2. 创建 `.vscode/settings.local.json`（已被 .gitignore 忽略）

## 故障排查

### Go 扩展未正常工作
1. 打开命令面板 (Ctrl+Shift+P)
2. 输入 "Go: Install/Update Tools"
3. 选择所有工具并安装

### 调试器无法启动
1. 确保安装了 Delve：`go install github.com/go-delve/delve/cmd/dlv@latest`
2. 检查 GOPATH/bin 是否在 PATH 中

### Proto 文件语法高亮不正常
安装 `zxh404.vscode-proto3` 扩展。

### 测试输出不显示详细信息
检查 `settings.json` 中是否配置了：
```json
"go.testFlags": ["-v"]
```

## 相关链接

- [VS Code Go 扩展文档](https://github.com/golang/vscode-go)
- [VS Code 调试指南](https://code.visualstudio.com/docs/editor/debugging)
- [REST Client 扩展文档](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)
