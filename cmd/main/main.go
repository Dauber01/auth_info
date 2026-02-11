package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// 初始化应用
	application, err := app.NewApp(cfg)
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
