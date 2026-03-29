@echo off
REM Build script for Windows (Alternative to Makefile)

setlocal

set PROJECT_NAME=auth_info
set BIN_DIR=bin
set API_DIR=api
set PROTO_DIR=%API_DIR%\proto
set GEN_DIR=%API_DIR%\gen
set MAIN_GO=cmd\main\main.go
set MIGRATE_GO=cmd\migrate\main.go
set SEED_GO=cmd\seed\main.go
set OUTPUT=%BIN_DIR%\%PROJECT_NAME%.exe
set CONFIG_DIR=.\config

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="install-tools" goto install_tools
if "%1"=="proto" goto proto
if "%1"=="wire" goto wire
if "%1"=="mod-tidy" goto mod_tidy
if "%1"=="migrate" goto migrate
if "%1"=="seed" goto seed
if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="clean" goto clean
if "%1"=="test" goto test
if "%1"=="fmt" goto fmt
if "%1"=="all" goto all

echo Unknown command: %1
goto help

:help
echo.
echo Available commands:
echo   make help           - Show this help message
echo   make install-tools  - Install protoc-gen-go tools
echo   make proto          - Generate protobuf Go code
echo   make wire           - Generate Wire dependency injection code
echo   make mod-tidy       - Update dependencies
echo   make migrate        - Run database migrations
echo   make seed           - Seed default casbin policies
echo   make build          - Build project
echo   make run            - Generate code and run project
echo   make clean          - Clean generated files
echo   make test           - Run tests
echo   make fmt            - Format code
echo   make all            - Execute all operations
echo.
goto end

:install_tools
echo Installing tools...
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo Tools installed
goto end

:proto
echo Generating proto code...
if not exist %GEN_DIR% mkdir %GEN_DIR%
protoc --proto_path=%PROTO_DIR%\third_party --proto_path=. --proto_path=%PROTO_DIR% --go_out=%GEN_DIR% --go_opt=paths=source_relative --go-grpc_out=%GEN_DIR% --go-grpc_opt=paths=source_relative %PROTO_DIR%\*.proto
if %errorlevel% neq 0 (
    echo Error: protoc not installed or generation failed
    echo Please install protoc: https://github.com/protocolbuffers/protobuf/releases
    goto end
)
echo Proto code generated
goto end

:wire
echo Generating Wire code...
go run github.com/google/wire/cmd/wire@latest .\internal\app
echo Wire code generated
goto end

:mod_tidy
echo Updating dependencies...
go mod tidy
echo Dependencies updated
goto end

:migrate
echo Running migrations...
go run %MIGRATE_GO% -config %CONFIG_DIR%
echo Migrations complete
goto end

:seed
echo Seeding default policies...
go run %SEED_GO% -config %CONFIG_DIR%
echo Policies seeded
goto end

:build
call :mod_tidy
call :wire
echo Building project...
if not exist %BIN_DIR% mkdir %BIN_DIR%
go build -o %OUTPUT% %MAIN_GO%
echo Build complete: %OUTPUT%
goto end

:run
call :clean
call :build
echo Starting application...
%OUTPUT% -config %CONFIG_DIR%
goto end

:clean
echo Cleaning up...
if exist %BIN_DIR% rd /s /q %BIN_DIR%
if exist %GEN_DIR% (
    del /q %GEN_DIR%\*.pb.go 2>nul
    del /q %GEN_DIR%\*_grpc.pb.go 2>nul
)
echo Cleanup complete
goto end

:test
echo Running tests...
go test -v .\...
goto end

:fmt
echo Formatting code...
go fmt .\...
echo Code formatted
goto end

:all
call :clean
call :proto
call :wire
call :build
echo All operations complete
goto end

:end
endlocal
