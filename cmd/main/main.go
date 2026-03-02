package main

// @title           Auth Info API
// @version         1.0
// @description     Auth Info Service API Documentation
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "auth_info/docs" // 注册 Swagger 文档
	"auth_info/internal/app"
	"auth_info/internal/config"
)

func main() {
	// 定义配置文件路径
	configPath := flag.String("config", "./config", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化应用（由 Wire 完成依赖注入）
	application, err := app.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// 创建通道用于优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 等待关闭信号
	<-sigChan
	if err := application.Stop(); err != nil {
		log.Fatalf("Server stop error: %v", err)
	}
}
