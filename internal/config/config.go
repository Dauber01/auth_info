// Package config loads an isolated, validated application configuration.
package config

import "time"

// Config contains deployment settings; business packages receive narrower options.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Casbin CasbinConfig `mapstructure:"casbin"`
}

// ServerConfig separates connection, ordinary request and long operation deadlines.
type ServerConfig struct {
	Port              int           `mapstructure:"port"`
	GRPCPort          int           `mapstructure:"grpc_port"`
	Mode              string        `mapstructure:"mode"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	RequestTimeout    time.Duration `mapstructure:"request_timeout"`
	GRPCTimeout       time.Duration `mapstructure:"grpc_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

// LogConfig supports console output and optional local rotation only.
type LogConfig struct {
	Level  string        `mapstructure:"level"`
	Format string        `mapstructure:"format"`
	Access bool          `mapstructure:"access"`
	File   LogFileConfig `mapstructure:"file"`
}

// LogFileConfig controls bounded local file retention.
type LogFileConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
}

// MySQLConfig describes a single database and its pool.
type MySQLConfig struct {
	Host     string          `mapstructure:"host"`
	Port     int             `mapstructure:"port"`
	User     string          `mapstructure:"user"`
	Password string          `mapstructure:"password"`
	DBName   string          `mapstructure:"dbname"`
	Charset  string          `mapstructure:"charset"`
	Pool     MySQLPoolConfig `mapstructure:"pool"`
}

// MySQLPoolConfig bounds connection usage.
type MySQLPoolConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// JWTConfig configures token signing and expiry in hours.
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

// CasbinConfig selects the RBAC model.
type CasbinConfig struct {
	Model string `mapstructure:"model"`
}
