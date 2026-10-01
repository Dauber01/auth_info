package logger

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"auth_info/internal/config"
	"auth_info/internal/pkg/trace"
)

func TestLoggerInstancesAndLocalFileCleanup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "app.log")
	cfg := config.LogConfig{Level: "info", Format: "json", File: config.LogFileConfig{
		Enabled: true, Path: path, MaxSizeMB: 1, MaxBackups: 1, MaxAgeDays: 1,
	}}
	first, closeFirst, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = closeFirst() })
	second, closeSecond, err := New(config.LogConfig{Level: "error", Format: "console"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = closeSecond() })
	if first == second {
		t.Fatal("shared logger")
	}
	WithContext(first, trace.WithID(context.Background(), "fixture-trace")).Info("owned event")
	if err := closeFirst(); err != nil {
		t.Fatal(err)
	}
	if err := closeFirst(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &record); err != nil {
		t.Fatal(err)
	}
	if record["trace_id"] != "fixture-trace" || record["msg"] != "owned event" {
		t.Fatal(record)
	}
}
