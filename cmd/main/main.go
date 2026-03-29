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
	configPath := flag.String("config", "./config", "配置文件路径")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- application.Run()
	}()

	select {
	case runErr := <-runErrCh:
		if runErr != nil {
			log.Fatalf("Server error: %v", runErr)
		}
	case sig := <-sigChan:
		log.Printf("Received signal %s, shutting down...", sig)
		if err := application.Stop(); err != nil {
			log.Fatalf("Server stop error: %v", err)
		}
		if runErr := <-runErrCh; runErr != nil {
			log.Fatalf("Server error during shutdown: %v", runErr)
		}
	}
}
