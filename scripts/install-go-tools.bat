@echo off
REM Go 代码格式化工具安装脚本（Windows 版）

setlocal

echo ========================================
echo 正在安装 Go 代码格式化和检查工具...
echo ========================================
echo.

REM 设置代理（如需要，取消注释）
REM set GOPROXY=https://goproxy.cn,direct

REM 1. gofmt (Go 自带，无需安装)
echo [1/10] gofmt - Go 标准格式化工具（已内置）
echo.

REM 2. goimports - 自动整理导入
echo [2/10] 正在安装 goimports...
go install golang.org/x/tools/cmd/goimports@latest
if %errorlevel% equ 0 (
    echo OK - goimports 安装成功
) else (
    echo FAIL - goimports 安装失败
)
echo.

REM 3. gofumpt - 更严格的格式化工具
echo [3/10] 正在安装 gofumpt...
go install mvdan.cc/gofumpt@latest
if %errorlevel% equ 0 (
    echo OK - gofumpt 安装成功
) else (
    echo FAIL - gofumpt 安装失败
)
echo.

REM 4. golangci-lint - 综合代码检查工具
echo [4/10] 正在安装 golangci-lint...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
if %errorlevel% equ 0 (
    echo OK - golangci-lint 安装成功
) else (
    echo FAIL - golangci-lint 安装失败
)
echo.

REM 5. gopls - Go 语言服务器
echo [5/10] 正在安装 gopls...
go install golang.org/x/tools/gopls@latest
if %errorlevel% equ 0 (
    echo OK - gopls 安装成功
) else (
    echo FAIL - gopls 安装失败
)
echo.

REM 6. delve - Go 调试器
echo [6/10] 正在安装 delve...
go install github.com/go-delve/delve/cmd/dlv@latest
if %errorlevel% equ 0 (
    echo OK - delve 安装成功
) else (
    echo FAIL - delve 安装失败
)
echo.

REM 7. gogetdoc - 文档查看工具
echo [7/10] 正在安装 gogetdoc...
go install github.com/zmb3/gogetdoc@latest
if %errorlevel% equ 0 (
    echo OK - gogetdoc 安装成功
) else (
    echo FAIL - gogetdoc 安装失败
)
echo.

REM 8. gomodifytags - 结构体标签工具
echo [8/10] 正在安装 gomodifytags...
go install github.com/fatih/gomodifytags@latest
if %errorlevel% equ 0 (
    echo OK - gomodifytags 安装成功
) else (
    echo FAIL - gomodifytags 安装失败
)
echo.

REM 9. impl - 接口实现生成工具
echo [9/10] 正在安装 impl...
go install github.com/josharian/impl@latest
if %errorlevel% equ 0 (
    echo OK - impl 安装成功
) else (
    echo FAIL - impl 安装失败
)
echo.

REM 10. gotests - 测试生成工具
echo [10/10] 正在安装 gotests...
go install github.com/cweill/gotests/gotests@latest
if %errorlevel% equ 0 (
    echo OK - gotests 安装成功
) else (
    echo FAIL - gotests 安装失败
)
echo.

echo ========================================
echo 安装完成！
echo ========================================
echo.
echo 已安装的工具：
echo   1. gofmt        - Go 标准格式化
echo   2. goimports    - 自动整理导入
echo   3. gofumpt      - 严格格式化
echo   4. golangci-lint - 代码检查
echo   5. gopls        - 语言服务器
echo   6. delve        - 调试器
echo   7. gogetdoc     - 文档查看
echo   8. gomodifytags - 标签工具
echo   9. impl         - 接口实现生成
echo  10. gotests      - 测试生成
echo.
echo 请在 VS Code 中打开命令面板 (Ctrl+Shift+P) 并运行：
echo   Go: Install/Update Tools
echo.
pause
endlocal
