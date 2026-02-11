# Go 代码格式化快速参考

## 🚀 一句话总结
**按 `Ctrl+S` 保存代码，它会自动格式化并整理导入！**

## 📦 安装工具

### Windows
```bash
.\scripts\install-go-tools.bat
```

### Linux/Mac
```bash
chmod +x scripts/install-go-tools.sh
./scripts/install-go-tools.sh
```

## ✨ 自动格式化功能

| 操作 | 快捷键 | 功能 |
|------|--------|------|
| 保存 | `Ctrl+S` | 自动格式化 + 整理导入 |
| 粘贴 | `Ctrl+V` | 自动格式化 |
| 手动格式化 | `Ctrl+Shift+P` → `Format` | 格式化当前文件 |

## 🛠️ 已安装的工具

| # | 工具 | 功能 | 状态 |
|----|------|------|------|
| 1 | **goimports** | 自动导入管理 + 格式化 | 🔴 默认 |
| 2 | **gofmt** | Go 官方标准格式化 | ⚪ 备选 |
| 3 | **gofumpt** | 更严格的格式化 | ⚪ 备选 |
| 4 | golangci-lint | 代码质量检查 | ✅ 启用 |
| 5 | gopls | 语言服务器 | ✅ 启用 |
| 6 | delve | 调试器 | ✅ 启用 |
| 7 | gogetdoc | 文档查看 | ✅ 启用 |
| 8 | gomodifytags | 标签生成 | ✅ 启用 |
| 9 | impl | 接口实现 | ✅ 启用 |
| 10 | gotests | 测试生成 | ✅ 启用 |

## 🔄 切换格式化工具

修改 `.vscode/settings.json`：

### 使用 gofmt（最基础）
```json
"go.formatTool": "gofmt"
```

### 使用 goimports（推荐）✨
```json
"go.formatTool": "goimports"
```

### 使用 gofumpt（最严格）
```json
"go.formatTool": "gofumpt"
```

## 🧪 测试格式化

### 在 VS Code 中测试
1. 新建 `test.go` 文件
2. 粘贴以下代码（格式混乱）：
```go
package main
import("fmt";"os")
func   main(){
fmt.Println("hello")
}
```
3. 按 `Ctrl+S` 保存
4. ✅ 代码自动格式化！

### 在终端中测试
```bash
# 查看哪些文件需要格式化
goimports -l .

# 格式化所有文件
goimports -w -local auth_info .

# 查看格式化差异
goimports -d main.go
```

## 📋 常用命令

```bash
# 安装/更新工具
.\scripts\install-go-tools.bat

# 检查代码质量
golangci-lint run ./...

# 格式化单个文件
goimports -w main.go

# 格式化整个项目
goimports -w -local auth_info .

# 查看差异但不修改
goimports -d .

# 使用不同工具格式化
gofmt -w main.go     # gofmt
gofumpt -w main.go   # gofumpt
```

## ⚙️ 配置位置

所有配置都在 `.vscode/settings.json`：
```json
// 格式化工具选择
"go.formatTool": "goimports",

// 格式化标志
"go.formatFlags": ["-local", "auth_info"],

// 保存时自动格式化
"editor.formatOnSave": true,

// 粘贴时自动格式化
"editor.formatOnPaste": true,

// 保存时整理导入
"editor.codeActionsOnSave": {
  "source.organizeImports": "explicit"
}
```

## 🎯 设置本地包路径

如果项目模块名不是 `auth_info`，修改：
```json
"go.formatFlags": ["-local", "your-module-name"]
```

## 📝 提交前检查清单

- [ ] 运行 `.\scripts\install-go-tools.bat` 更新工具
- [ ] 打开 VS Code 并编辑代码
- [ ] 按 `Ctrl+S` 保存（自动格式化）
- [ ] 运行 `golangci-lint run ./...` 检查质量
- [ ] 提交代码

## 🐛 常见问题

### Q: 保存时没有格式化？
A: 重启 VS Code，或运行 `Go: Install/Update Tools`

### Q: goimports 找不到？
A: 运行 `.\scripts\install-go-tools.bat` 重新安装

### Q: 导入没有整理？
A: 检查 `editor.codeActionsOnSave` 配置

### Q: 我想用 gofumpt？
A: 修改 `"go.formatTool": "gofumpt"`

### Q: 想禁用自动格式化？
A: 将 `"editor.formatOnSave": false`

## 🎓 学习更多

详细文档见：[FORMAT.md](FORMAT.md)

快捷键参考：
- `Ctrl+Shift+P` - 命令面板
- `Ctrl+S` - 保存（自动格式化）
- `Ctrl+.` - 快速修复
- `F12` - 跳转定义
- `Shift+F12` - 查看引用

---

**💡 提示**: 团队开发时，确保所有成员使用相同的配置文件（`.vscode/settings.json`），保持代码风格一致！
