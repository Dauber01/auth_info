#!/bin/bash
# Go 代码格式化工具安装脚本

echo "========================================"
echo "正在安装 Go 代码格式化和检查工具..."
echo "========================================"
echo ""

# 设置代理（如果需要）
# export GOPROXY=https://goproxy.cn,direct

# 1. gofmt (Go 自带，无需安装)
echo "✓ gofmt - Go 标准格式化工具（已内置）"
echo ""

# 2. goimports - 自动整理导入
echo "正在安装 goimports..."
go install golang.org/x/tools/cmd/goimports@latest
if [ $? -eq 0 ]; then
    echo "✓ goimports 安装成功"
else
    echo "✗ goimports 安装失败"
fi
echo ""

# 3. gofumpt - 更严格的格式化工具
echo "正在安装 gofumpt..."
go install mvdan.cc/gofumpt@latest
if [ $? -eq 0 ]; then
    echo "✓ gofumpt 安装成功"
else
    echo "✗ gofumpt 安装失败"
fi
echo ""

# 4. golangci-lint - 综合代码检查工具
echo "正在安装 golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
if [ $? -eq 0 ]; then
    echo "✓ golangci-lint 安装成功"
else
    echo "✗ golangci-lint 安装失败"
fi
echo ""

# 5. gopls - Go 语言服务器
echo "正在安装 gopls..."
go install golang.org/x/tools/gopls@latest
if [ $? -eq 0 ]; then
    echo "✓ gopls 安装成功"
else
    echo "✗ gopls 安装失败"
fi
echo ""

# 6. delve - Go 调试器
echo "正在安装 delve..."
go install github.com/go-delve/delve/cmd/dlv@latest
if [ $? -eq 0 ]; then
    echo "✓ delve 安装成功"
else
    echo "✗ delve 安装失败"
fi
echo ""

# 7. gogetdoc - 文档查看工具
echo "正在安装 gogetdoc..."
go install github.com/zmb3/gogetdoc@latest
if [ $? -eq 0 ]; then
    echo "✓ gogetdoc 安装成功"
else
    echo "✗ gogetdoc 安装失败"
fi
echo ""

# 8. gomodifytags - 结构体标签工具
echo "正在安装 gomodifytags..."
go install github.com/fatih/gomodifytags@latest
if [ $? -eq 0 ]; then
    echo "✓ gomodifytags 安装成功"
else
    echo "✗ gomodifytags 安装失败"
fi
echo ""

# 9. impl - 接口实现生成工具
echo "正在安装 impl..."
go install github.com/josharian/impl@latest
if [ $? -eq 0 ]; then
    echo "✓ impl 安装成功"
else
    echo "✗ impl 安装失败"
fi
echo ""

# 10. gotests - 测试生成工具
echo "正在安装 gotests..."
go install github.com/cweill/gotests/gotests@latest
if [ $? -eq 0 ]; then
    echo "✓ gotests 安装成功"
else
    echo "✗ gotests 安装失败"
fi
echo ""

echo "========================================"
echo "安装完成！"
echo "========================================"
echo ""
echo "已安装的工具："
echo "  1. gofmt        - Go 标准格式化"
echo "  2. goimports    - 自动整理导入"
echo "  3. gofumpt      - 严格格式化"
echo "  4. golangci-lint - 代码检查"
echo "  5. gopls        - 语言服务器"
echo "  6. delve        - 调试器"
echo "  7. gogetdoc     - 文档查看"
echo "  8. gomodifytags - 标签工具"
echo "  9. impl         - 接口实现生成"
echo " 10. gotests      - 测试生成"
echo ""
echo "请确保 \$GOPATH/bin 在您的 PATH 中："
echo "  export PATH=\$PATH:\$(go env GOPATH)/bin"
echo ""
echo "在 VS Code 中，打开命令面板 (Ctrl+Shift+P) 并运行："
echo "  Go: Install/Update Tools"
echo ""
