package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfig_MySQLPoolDuration(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "test.yaml")
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

func writeConfig(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadConfigPrecedenceIsolationAndEnvironmentOnlyFields(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "base.yaml", "server:\n  port: 8100\nlog:\n  level: warn\njwt:\n  secret: fixture\n")
	entry := writeConfig(t, dir, "entry.yaml", "includes: [base.yaml]\nserver:\n  port: 8200\n")
	t.Setenv("APP_SERVER_PORT", "8300")
	t.Setenv("APP_DOCUMENT_FONT_PATH", "injected.ttf")
	cfg, err := LoadConfig(entry)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 8300 || cfg.Server.GRPCPort != 9300 || cfg.Log.Level != "warn" ||
		cfg.Document.FontPath != "injected.ttf" {
		t.Fatalf("unexpected precedence: %+v", cfg.Server)
	}
	t.Setenv("APP_SERVER_PORT", "8400")
	second := writeConfig(t, dir, "second.yaml", "jwt:\n  secret: different\n")
	cfg2, err := LoadConfig(second)
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.Log.Level != "info" || cfg.JWT.Secret != "fixture" || cfg.Server.Port != 8300 {
		t.Fatal("configuration leaked between loads")
	}
}

func TestLoadConfigIncludesAndValidationFailures(t *testing.T) {
	for name, body := range map[string]string{
		"cycle":            "includes: [test.yaml]\njwt:\n  secret: fixture\n",
		"missing":          "includes: [missing.yaml]\n",
		"invalid_port":     "server:\n  port: -1\njwt:\n  secret: fixture\n",
		"invalid_duration": "server:\n  request_timeout: -1s\njwt:\n  secret: fixture\n",
		"invalid_budget":   "server:\n  write_timeout: 1s\njwt:\n  secret: fixture\n",
		"empty_secret":     "jwt:\n  secret: ''\n",
		"invalid_log":      "log:\n  format: invalid\njwt:\n  secret: fixture\n",
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			writeConfig(t, dir, "test.yaml", body)
			if _, err := LoadConfig(dir); err == nil {
				t.Fatal("expected invalid configuration")
			}
		})
	}
	dir := t.TempDir()
	writeConfig(t, dir, "test.yaml", "server:\n  request_timeout: 0s\n  write_timeout: 0s\njwt:\n  secret: fixture\n")
	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.RequestTimeout != 0 || cfg.Server.WriteTimeout != 0 {
		t.Fatal("explicit disable lost")
	}
}

func TestIncludeOrderAndEntryOverride(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "first.yaml", "server:\n  port: 8100\nlog:\n  level: warn\njwt:\n  secret: fixture\n")
	writeConfig(t, dir, "second.yaml", "server:\n  port: 8200\nlog:\n  level: error\n")
	path := writeConfig(t, dir, "entry.yaml", "includes: [first.yaml, second.yaml]\nserver:\n  port: 8300\n")
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 8300 || cfg.Log.Level != "error" {
		t.Fatal("merge precedence changed")
	}
}

// Deployment environment variables must not influence the checked-in profile tests.
func clearConfigEnvironment(t *testing.T) {
	t.Helper()
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "APP_") {
			t.Setenv(key, "")
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestLoadConfigDefaultAndDirectorySelectTest(t *testing.T) {
	clearConfigEnvironment(t)
	root := t.TempDir()
	dir := filepath.Join(root, "config")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, dir, "test.yaml", "server:\n  mode: test\njwt:\n  secret: test-fixture\n")
	writeConfig(t, dir, "line.yaml", "server:\n  mode: release\njwt:\n  secret: line-fixture\n")
	writeConfig(t, dir, "config.yaml", "server:\n  mode: debug\njwt:\n  secret: legacy-fixture\n")
	t.Chdir(root)
	for _, path := range []string{"", "./config", DefaultPath} {
		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Server.Mode != "test" || cfg.JWT.Secret != "test-fixture" {
			t.Fatalf("%q selected the wrong environment", path)
		}
	}
	line, err := LoadConfig("./config/line.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if line.Server.Mode != "release" || line.JWT.Secret != "line-fixture" {
		t.Fatal("explicit environment ignored")
	}
	if err := os.Remove(filepath.Join(dir, "test.yaml")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"", "./config"} {
		if _, err := LoadConfig(path); err == nil {
			t.Fatalf("%q silently fell back to another environment", path)
		}
	}
}

func TestRepositoryEnvironmentProfiles(t *testing.T) {
	clearConfigEnvironment(t)
	files, err := filepath.Glob("../../config/*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || filepath.Base(files[0]) != "line.yaml" || filepath.Base(files[1]) != "test.yaml" {
		t.Fatalf("unexpected environment entries: %v", files)
	}
	testCfg, err := LoadConfig("../../config/test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if testCfg.Server.Mode != "test" || testCfg.Log.Level != "warn" || testCfg.JWT.Secret == "" {
		t.Fatal("test profile lost its settings")
	}
	if _, err := LoadConfig("../../config/line.yaml"); err == nil {
		t.Fatal("line must require its own signing secret")
	}
	t.Setenv("APP_JWT_SECRET", "line-fixture-secret")
	t.Setenv("APP_MYSQL_HOST", "mysql.example.invalid")
	t.Setenv("APP_MYSQL_USER", "line-fixture-user")
	t.Setenv("APP_MYSQL_PASSWORD", "line-fixture-password")
	t.Setenv("APP_MYSQL_DBNAME", "line_fixture")
	lineCfg, err := LoadConfig("../../config/line.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if lineCfg.Server.Mode != "release" || lineCfg.Log.Level != "info" ||
		lineCfg.JWT.Secret != "line-fixture-secret" || lineCfg.MySQL.Host != "mysql.example.invalid" ||
		lineCfg.MySQL.User != "line-fixture-user" || lineCfg.MySQL.Password != "line-fixture-password" ||
		lineCfg.MySQL.DBName != "line_fixture" {
		t.Fatal("line profile did not use deployment configuration")
	}
}
