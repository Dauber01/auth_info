package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestLoadConfig_MySQLPoolDuration(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := []byte(`
server:
  port: 8080
  mode: debug
log:
  level: debug
  format: json
mysql:
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: auth_info
  charset: utf8mb4
  pool:
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: "1h"
    conn_max_idle_time: "10m"
jwt:
  secret: test
  expire: 24
casbin:
  model: config/rbac_model.conf
`)

	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	viper.Reset()
	t.Cleanup(viper.Reset)

	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.MySQL.Pool.ConnMaxLifetime != time.Hour {
		t.Fatalf("unexpected conn max lifetime: %v", cfg.MySQL.Pool.ConnMaxLifetime)
	}
	if cfg.MySQL.Pool.ConnMaxIdleTime != 10*time.Minute {
		t.Fatalf("unexpected conn max idle time: %v", cfg.MySQL.Pool.ConnMaxIdleTime)
	}
}
