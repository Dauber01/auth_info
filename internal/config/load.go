package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// DefaultPath selects the test environment unless a deployment explicitly chooses another file.
const DefaultPath = "./config/test.yaml"

// LoadConfig accepts a file or a directory containing test.yaml.
// Precedence: defaults < includes (in order) < entry file < APP_* variables.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat config: %w", err)
	}
	if info.IsDir() {
		path = filepath.Join(path, "test.yaml")
	}
	values, err := readConfig(path, make(map[string]bool))
	if err != nil {
		return nil, err
	}
	v := viper.New()
	setDefaults(v)
	if err := v.MergeConfigMap(values); err != nil {
		return nil, fmt.Errorf("merge config: %w", err)
	}
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)
	// Bind every leaf explicitly so environment-only fields are also unmarshaled.
	if err := bindEnv(v, reflect.TypeOf(Config{}), ""); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	if cfg.Server.GRPCPort == 0 {
		cfg.Server.GRPCPort = cfg.Server.Port + 1000
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func bindEnv(v *viper.Viper, typ reflect.Type, prefix string) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		key := prefix + field.Tag.Get("mapstructure")
		if field.Type.Kind() == reflect.Struct {
			if err := bindEnv(v, field.Type, key+"."); err != nil {
				return err
			}
		} else if err := v.BindEnv(key); err != nil {
			return fmt.Errorf("bind %s: %w", key, err)
		}
	}
	return nil
}

func readConfig(path string, visiting map[string]bool) (map[string]any, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if visiting[path] {
		return nil, fmt.Errorf("config includes cycle at %s", path)
	}
	visiting[path] = true
	defer delete(visiting, path)
	entry := viper.New()
	entry.SetConfigFile(path)
	if err := entry.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	merged := viper.New()
	for _, include := range entry.GetStringSlice("includes") {
		if !filepath.IsAbs(include) {
			include = filepath.Join(filepath.Dir(path), include)
		}
		values, err := readConfig(include, visiting)
		if err != nil {
			return nil, err
		}
		if err := merged.MergeConfigMap(values); err != nil {
			return nil, err
		}
	}
	values := entry.AllSettings()
	delete(values, "includes")
	if err := merged.MergeConfigMap(values); err != nil {
		return nil, err
	}
	return merged.AllSettings(), nil
}

func setDefaults(v *viper.Viper) {
	for key, value := range map[string]any{
		"server.port": 8080, "server.mode": "debug", "server.grpc_port": 0,
		"server.read_header_timeout": "5s", "server.read_timeout": "15s",
		"server.write_timeout": "35s", "server.idle_timeout": "120s",
		"server.request_timeout": "30s",
		"server.grpc_timeout":    "30s", "server.shutdown_timeout": "5s",
		"log.level": "info", "log.format": "json", "log.access": true,
		"log.file.path": "logs/app.log", "log.file.max_size_mb": 100,
		"log.file.max_backups": 10, "log.file.max_age_days": 30,
		"mysql.host": "localhost", "mysql.port": 3306, "mysql.charset": "utf8mb4",
		"mysql.pool.max_open_conns": 100, "mysql.pool.max_idle_conns": 10,
		"mysql.pool.conn_max_lifetime": "1h", "mysql.pool.conn_max_idle_time": "10m",
		"jwt.expire": 24, "casbin.model": "config/rbac_model.conf",
	} {
		v.SetDefault(key, value)
	}
}
