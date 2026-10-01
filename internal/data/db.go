package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"auth_info/internal/config"
)

const (
	defaultMySQLMaxOpenConns    = 100
	defaultMySQLMaxIdleConns    = 10
	defaultMySQLConnMaxLifetime = time.Hour
	defaultMySQLConnMaxIdleTime = 10 * time.Minute
	mysqlPingTimeout            = 5 * time.Second
)

func NewDB(cfg *config.Config, log *zap.Logger) (db *gorm.DB, err error) {
	c := cfg.MySQL
	pool, err := normalizeMySQLPoolConfig(c.Pool)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql pool: %w", err)
	}
	// The provider owns the pool until initialization completes, including GORM version detection.
	defer func() {
		if err != nil {
			err = errors.Join(err, sqlDB.Close())
		}
	}()
	applyMySQLPoolConfig(sqlDB, pool)
	ctx, cancel := context.WithTimeout(context.Background(), mysqlPingTimeout)
	defer cancel()
	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("initialize mysql: %w", err)
	}

	log.Info("MySQL connected",
		zap.String("db", c.DBName),
		zap.Int("max_open_conns", pool.MaxOpenConns),
		zap.Int("max_idle_conns", pool.MaxIdleConns),
		zap.Duration("conn_max_lifetime", pool.ConnMaxLifetime),
		zap.Duration("conn_max_idle_time", pool.ConnMaxIdleTime),
	)
	return db, nil
}

func normalizeMySQLPoolConfig(pool config.MySQLPoolConfig) (config.MySQLPoolConfig, error) {
	if pool.MaxOpenConns < 0 {
		return pool, fmt.Errorf("mysql pool max_open_conns cannot be negative")
	}
	if pool.MaxIdleConns < 0 {
		return pool, fmt.Errorf("mysql pool max_idle_conns cannot be negative")
	}
	if pool.ConnMaxLifetime < 0 {
		return pool, fmt.Errorf("mysql pool conn_max_lifetime cannot be negative")
	}
	if pool.ConnMaxIdleTime < 0 {
		return pool, fmt.Errorf("mysql pool conn_max_idle_time cannot be negative")
	}

	if pool.MaxOpenConns == 0 {
		pool.MaxOpenConns = defaultMySQLMaxOpenConns
	}
	if pool.MaxIdleConns == 0 {
		pool.MaxIdleConns = defaultMySQLMaxIdleConns
	}
	if pool.ConnMaxLifetime == 0 {
		pool.ConnMaxLifetime = defaultMySQLConnMaxLifetime
	}
	if pool.ConnMaxIdleTime == 0 {
		pool.ConnMaxIdleTime = defaultMySQLConnMaxIdleTime
	}

	if pool.MaxIdleConns > pool.MaxOpenConns {
		return pool, fmt.Errorf("mysql pool max_idle_conns cannot exceed max_open_conns")
	}

	return pool, nil
}

func applyMySQLPoolConfig(sqlDB *sql.DB, pool config.MySQLPoolConfig) {
	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
}

// RunMigrations 由 cmd/migrate 显式传入各模块持久化模型，避免 data 反向依赖业务包。
func RunMigrations(db *gorm.DB, models ...any) error {
	if err := db.AutoMigrate(models...); err != nil {
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

// CloseDB closes the SQL pool owned by the caller.
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get mysql pool: %w", err)
	}
	return sqlDB.Close()
}
