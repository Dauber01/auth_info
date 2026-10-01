package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"auth_info/internal/config"
	"auth_info/internal/data"
	dataauth "auth_info/internal/data/auth"
	datadict "auth_info/internal/data/dict"
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
	return data.RunMigrations(db, &dataauth.User{}, &datadict.DictType{}, &datadict.DictItem{})
}

func main() {
	path := flag.String("config", "./config", "配置文件或包含 config.yaml 的目录")
	flag.Parse()
	if err := run(*path); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
