package main

import (
	"flag"
	"log"

	"auth_info/internal/config"
	"auth_info/internal/data"
	dataauth "auth_info/internal/data/auth"
	datadict "auth_info/internal/data/dict"
	"auth_info/internal/logger"
)

func main() {
	configPath := flag.String("config", "./config", "配置文件路径")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	logg := logger.GetLogger()
	db, err := data.NewDB(cfg, logg)
	if err != nil {
		log.Fatalf("Failed to connect mysql: %v", err)
	}

	// 新增模块时在此处追加对应持久化模型即可。
	if err := data.RunMigrations(db,
		&dataauth.User{},
		&datadict.DictType{},
		&datadict.DictItem{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}
