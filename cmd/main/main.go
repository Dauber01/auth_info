package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth_info/internal/app"
	"auth_info/internal/config"
)

func run(ctx context.Context, path string) error {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return err
	}
	application, err := app.InitializeApp(cfg)
	if err != nil {
		return err
	}
	result := make(chan error, 1)
	go func() { result <- application.Run() }()
	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		stopErr := application.Stop()
		return errors.Join(stopErr, <-result)
	}
}

func main() {
	path := flag.String("config", config.DefaultPath, "配置文件或包含 test.yaml 的目录")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *path); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
