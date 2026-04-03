package main

import (
	"flag"
	"log"

	"auth_info/internal/config"
	"auth_info/internal/data"
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

	enforcer, err := data.NewEnforcer(db, cfg, logg)
	if err != nil {
		log.Fatalf("Failed to init casbin enforcer: %v", err)
	}

	if err := data.SeedDefaultPolicies(enforcer); err != nil {
		log.Fatalf("Seed policies failed: %v", err)
	}

	log.Println("Default policies seeded successfully")
}
