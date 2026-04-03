package data

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"auth_info/internal/config"
)

func NewDB(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	c := cfg.MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	log.Info("MySQL connected", zap.String("db", c.DBName))
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &DictType{}, &DictItem{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

func NewEnforcer(db *gorm.DB, cfg *config.Config, log *zap.Logger) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("casbin adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(cfg.Casbin.Model, adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin enforcer: %w", err)
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load policy: %w", err)
	}

	log.Info("Casbin enforcer initialized")
	return enforcer, nil
}

func SeedDefaultPolicies(e *casbin.Enforcer) error {
	policies := [][]string{
		{"admin", "/api/v1/*", "*"},
		{"user", "/api/v1/hello", "GET"},
		{"user", "/api/v1/dict/types", "GET"},
		{"user", "/api/v1/dict/items", "GET"},
	}
	for _, p := range policies {
		hasPolicy, err := e.HasPolicy(p[0], p[1], p[2])
		if err != nil {
			return fmt.Errorf("check policy role=%s path=%s act=%s: %w", p[0], p[1], p[2], err)
		}
		if hasPolicy {
			continue
		}

		if _, err = e.AddPolicy(p[0], p[1], p[2]); err != nil {
			return fmt.Errorf("add policy role=%s path=%s act=%s: %w", p[0], p[1], p[2], err)
		}
	}
	return nil
}
