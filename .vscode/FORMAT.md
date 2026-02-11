# Go 代码自动格式化配置指南

本文档介绍如何在 VS Code 中配置 Go 代码自动格式化功能。

## 快速开始

### 1. 安装格式化工具

#### Windows 环境
```bash
cd d:\GoProject\auth_info
.\scripts\install-go-tools.bat
```

#### Linux/Mac 环境
```bash
cd /path/to/auth_info
chmod +x scripts/install-go-tools.sh
./scripts/install-go-tools.sh
```

### 2. 在 VS Code 中安装/更新工具

打开命令面板（`Ctrl+Shift+P` 或 `Cmd+Shift+P`），输入：
```
Go: Install/Update Tools
```
勾选所有工具并点击 "OK" 安装。

### 3. 测试自动格式化

1. 打开任意 `.go` 文件
2. 修改代码（例如删除一些空格或换行）
3. 按 `Ctrl+S` 保存
4. **代码会自动格式化！** ✨

## 已配置的格式化工具

### 1. **goimports** (推荐，默认使用)

**功能**：
- ✅ 自动格式化代码（包含 gofmt 的所有功能）
- ✅ 自动添加缺失的导入
- ✅ 自动删除未使用的导入
- ✅ 按规范排序导入语句

**配置**：
```json
"go.formatTool": "goimports",
"go.formatFlags": ["-local", "auth_info"]
```

**示例**：
```go
// 保存前
package main
func main() {
fmt.Println("hello")
}

// 保存后（自动添加导入）
package main

import "fmt"

func main() {
	fmt.Println("hello")
}
```

### 2. **gofmt** (Go 官方标准)

**功能**：
- ✅ 标准的 Go 代码格式化
- ✅ 统一缩进、空格、换行

**切换方式**：
修改 `.vscode/settings.json`：
```json
"go.formatTool": "gofmt"
```

### 3. **gofumpt** (更严格的格式化)

**功能**：
- ✅ 包含 gofmt 的所有功能
- ✅ 更严格的格式化规则
- ✅ 删除多余的空行
- ✅ 优化结构体字段对齐

**切换方式**：
修改 `.vscode/settings.json`：
```json
"go.formatTool": "gofumpt"
```

## 自动化功能

### 保存时自动执行

当您按 `Ctrl+S` 保存文件时，会自动执行：

1. **格式化代码** - 使用 goimports
2. **整理导入** - 添加/删除/排序导入语句
3. **修复问题** - 自动修复简单的代码问题

配置位于 `.vscode/settings.json`：
```json
"[go]": {
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll": "explicit",
    "source.organizeImports": "explicit"
  }
}
```

### 粘贴时自动格式化

粘贴代码时也会自动格式化：
```json
"editor.formatOnPaste": true
```

## 代码质量检查

### golangci-lint

配置了 golangci-lint 进行代码质量检查：

```json
"go.linter": "golangci-lint",
"go.linterFlags": [
  "--fast",
  "--max-issues-per-linter=0",
  "--max-same-issues=0"
]
```

**检查内容**：
- ❌ 未使用的变量
- ❌ 未使用的参数
- ❌ 潜在的 nil 指针
- ❌ 不可达的代码
- ❌ 代码风格问题

### 手动运行检查

在终端中执行：
```bash
golangci-lint run ./...
```

## 格式化选项对比

| 特性 | gofmt | goimports | gofumpt |
|-----|-------|-----------|---------|
| 基础格式化 | ✅ | ✅ | ✅ |
| 自动导入管理 | ❌ | ✅ | ✅ |
| 严格格式化 | ❌ | ❌ | ✅ |
| Go 官方标准 | ✅ | ✅ | ❌ |
| 推荐使用 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

## 常用命令

### VS Code 命令面板

- `Ctrl+Shift+P` → `Format Document` - 手动格式化当前文件
- `Ctrl+Shift+P` → `Go: Install/Update Tools` - 安装/更新工具
- `Ctrl+Shift+P` → `Organize Imports` - 整理导入

### 终端命令

```bash
# 格式化单个文件
goimports -w main.go

# 格式化整个项目
goimports -w -local auth_info .

# 使用 gofumpt 格式化
gofumpt -w .

# 检查代码（不修改）
gofmt -l .

# 查看格式化差异
gofmt -d main.go
```

## 自定义配置

### 修改本地包路径

如果您的项目模块名不是 `auth_info`，请修改：

`.vscode/settings.json`：
```json
"go.formatFlags": ["-local", "your-module-name"]
```

### 使用不同的格式化工具

切换到 gofumpt：
```json
"go.formatTool": "gofumpt"
```

### 禁用特定文件的自动格式化

在文件顶部添加注释：
```go
//go:build tools
// +build tools

// gofmt: off
package main
```

## 团队协作

### 统一格式化标准

确保团队所有成员：
1. 使用相同的 VS Code 配置（`.vscode/settings.json`）
2. 安装相同的工具版本
3. 提交代码前运行格式化

### Git 提交前自动格式化

创建 `.git/hooks/pre-commit`：
```bash
#!/bin/bash
# 格式化所有修改的 Go 文件
FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
if [ -n "$FILES" ]; then
    goimports -w -local auth_info $FILES
    git add $FILES
fi
```

## 故障排查

### 问题 1: 保存时没有自动格式化

**解决方案**：
1. 检查 VS Code 设置是否正确
2. 打开命令面板，运行 `Go: Install/Update Tools`
3. 重启 VS Code

### 问题 2: goimports 找不到

**解决方案**：
1. 确保 `$GOPATH/bin` 在 PATH 中
2. 运行安装脚本重新安装：`.\scripts\install-go-tools.bat`

### 问题 3: 导入没有自动整理

**解决方案**：
检查 `.vscode/settings.json` 中是否配置：
```json
"editor.codeActionsOnSave": {
  "source.organizeImports": "explicit"
}
```

### 问题 4: 格式化很慢

**解决方案**：
1. 排除大型文件夹（如 vendor、node_modules）
2. 使用 gofmt 替代 gofumpt
3. 关闭不必要的 linter

### 问题 5: Windows 上工具无法执行

**解决方案**：
确保 Go 安装正确，检查环境变量：
```powershell
go env GOPATH
go env GOROOT
```

添加到 PATH：
```powershell
$env:PATH += ";$(go env GOPATH)\bin"
```

## 最佳实践

### 1. 保存即格式化
养成按 `Ctrl+S` 保存的习惯，代码自动保持整洁。

### 2. 提交前检查
提交代码前运行：
```bash
make fmt
golangci-lint run ./...
```

### 3. 统一团队配置
将 `.vscode/settings.json` 提交到版本控制。

### 4. 使用 goimports
推荐使用 goimports，自动管理导入语句。

### 5. 定期更新工具
每月运行一次：
```bash
.\scripts\install-go-tools.bat
```

## 参考链接

- [gofmt 官方文档](https://pkg.go.dev/cmd/gofmt)
- [goimports 官方文档](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [gofumpt GitHub](https://github.com/mvdan/gofumpt)
- [golangci-lint 官方文档](https://golangci-lint.run/)
- [VS Code Go 扩展文档](https://github.com/golang/vscode-go)

## 快速测试

创建测试文件 `test_format.go`：

```go
package main
import("fmt";"time")
func   main(){
fmt.Println("hello")
time.Sleep(1)
}
```

保存文件（`Ctrl+S`），代码会自动格式化为：

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello")
	time.Sleep(1)
}
```

🎉 如果看到这样的效果，说明配置成功！
