# 环境配置

2026-10-01 按用户要求收敛为 test 和 line 两个环境入口。来源：
[配置目录](../config/)、[加载器](../internal/config/load.go)、[Makefile](../Makefile)。

## 文件职责

| 文件 | 用途 |
| --- | --- |
| [test.yaml](../config/test.yaml) | 测试环境与默认启动入口；mode=test、log.level=warn，保留原测试数据库/JWT 设置 |
| [line.yaml](../config/line.yaml) | 线上环境；mode=release、log.level=info，部署凭据由 APP_* 注入 |
| [includes/base.yaml](../config/includes/base.yaml) | 两个环境共用的端口、超时、日志参数，不是第三个环境 |
| [rbac_model.conf](../config/rbac_model.conf) | 共用 Casbin 权限模型，不随环境复制 |

已移除 dev.yaml、pre.yaml 和旧 config.yaml。test 与 line 都直接包含公共 base，彼此不继承。
配置优先级：代码默认值 < includes（按声明顺序）< 环境入口文件 < APP_*。
includes 路径相对声明文件解析；Casbin 模型等资源路径相对进程工作目录解析。

## 启动与命令

在仓库根目录运行：

```sh
make run                    # test
make run ENV=test
make run ENV=line           # 先注入线上环境参数
./bin/auth_info              # 默认 ./config/test.yaml
./bin/auth_info -config ./config/line.yaml
```

migrate、seed 同样默认 test，可通过 ENV=line 显式选择线上配置；执行会写数据库。
本轮仅验证命令选择，不执行这些数据库操作。make dev 是“生成后启动”的开发命令，默认也选 test。

Make 的 ENV 只接受 test/line；其他名称立即报错。保留 CONFIG_DIR 与 CONFIG_FILE 覆盖：

```sh
make config ENV=line
make config ENV=line CONFIG_DIR=/etc/auth_info
make config CONFIG_FILE=/etc/auth_info/line.yaml
```

CONFIG_FILE 显式给定时优先于 ENV/CONFIG_DIR 的路径拼接。直接命令的 -config 可接具体文件或目录；
目录固定读取 test.yaml，不根据目录中有哪些文件猜环境。空路径同样使用默认 test 文件。
若 test.yaml 缺失就报错，不回退到旧 config.yaml 或 line.yaml；旧目录形式 `-config ./config` 仍可用。

## 环境覆盖与参数分组

环境变量名为 APP_ 加字段路径大写，各层用下划线连接，例如 server.port → APP_SERVER_PORT。

| 分组 | 主要字段与职责 |
| --- | --- |
| server | HTTP/gRPC 端口、Gin mode、连接/普通请求/gRPC/关闭期限；grpc_port 为 0 或省略时使用 HTTP port + 1000 |
| log | level/format/access，以及可选本地滚动文件 file；不上传 ES |
| mysql | host/port/user/password/dbname/charset 与连接池 pool |
| jwt | secret 与按小时计的 expire |
| casbin | model 权限模型路径 |

line 的部署配置至少应提供 APP_MYSQL_HOST、APP_MYSQL_USER、APP_MYSQL_PASSWORD、APP_MYSQL_DBNAME、
APP_JWT_SECRET；端口等按需覆盖。未注入有效 JWT secret 时配置校验失败。
线上凭据不从 test.yaml 继承，也不在公共 base 中保存。测试环境现有设置本轮保持，只调整归属。
完整字段与校验规则见 [配置类型](../internal/config/config.go) 和 [校验](../internal/config/validate.go)。

验证记录见 [架构任务的配置收敛批次](tasks/20260930-structure-fusion-review/05-verification.md)。
复核触发：环境文件、默认路径、Make 选择规则、配置字段或资源路径语义改变。
