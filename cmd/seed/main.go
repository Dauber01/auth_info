package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"auth_info/internal/config"
	"auth_info/internal/data"
	"auth_info/internal/pkg/logger"
)

func run(path string) (err error) {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return err
	}
	logg, closeLog, err := logger.New(cfg.Log)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, closeLog()) }()
	db, err := data.NewDB(cfg, logg)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, data.CloseDB(db)) }()
	enforcer, err := data.NewEnforcer(db, cfg, logg)
	if err != nil {
		return err
	}
	return data.SeedDefaultPolicies(enforcer)
}

func main() {
	path := flag.String("config", "./config", "配置文件或包含 config.yaml 的目录")
	flag.Parse()
	if err := run(*path); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
