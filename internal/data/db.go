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
	"auth_info/internal/logger"
	"auth_info/internal/model"
)

// NewDB Wire Provider：初始化 MySQL 连接，自动迁移表结构
func NewDB(cfg *config.Config) (*gorm.DB, error) {
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

	if err = db.AutoMigrate(&model.User{}, &model.DictType{}, &model.DictItem{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	logger.GetLogger().Info("MySQL connected", zap.String("db", c.DBName))
	return db, nil
}

// NewEnforcer Wire Provider：初始化 Casbin + gorm-adapter，并写入默认策略
func NewEnforcer(db *gorm.DB, cfg *config.Config) (*casbin.Enforcer, error) {
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

	seedDefaultPolicies(enforcer)

	logger.GetLogger().Info("Casbin enforcer initialized")
	return enforcer, nil
}

// seedDefaultPolicies 写入默认 RBAC 策略（幂等）
func seedDefaultPolicies(e *casbin.Enforcer) {
	// admin 角色拥有所有权限
	policies := [][]string{
		{"admin", "/api/v1/*", "*"},
		{"user", "/api/v1/hello", "GET"},
		{"user", "/api/v1/dict/types", "GET"},
		{"user", "/api/v1/dict/items", "GET"},
	}
	for _, p := range policies {
		if ok, _ := e.HasPolicy(p[0], p[1], p[2]); !ok {
			_, _ = e.AddPolicy(p[0], p[1], p[2])
		}
	}
}
