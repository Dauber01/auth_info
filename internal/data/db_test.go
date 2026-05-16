package data

import (
	"strings"
	"testing"
	"time"

	"auth_info/internal/config"
)

func TestNormalizeMySQLPoolConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   config.MySQLPoolConfig
		want    config.MySQLPoolConfig
		wantErr string
	}{
		{
			name:  "defaults",
			input: config.MySQLPoolConfig{},
			want: config.MySQLPoolConfig{
				MaxOpenConns:    defaultMySQLMaxOpenConns,
				MaxIdleConns:    defaultMySQLMaxIdleConns,
				ConnMaxLifetime: defaultMySQLConnMaxLifetime,
				ConnMaxIdleTime: defaultMySQLConnMaxIdleTime,
			},
		},
		{
			name: "custom",
			input: config.MySQLPoolConfig{
				MaxOpenConns:    50,
				MaxIdleConns:    5,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
			want: config.MySQLPoolConfig{
				MaxOpenConns:    50,
				MaxIdleConns:    5,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 5 * time.Minute,
			},
		},
		{
			name: "negative max open connections",
			input: config.MySQLPoolConfig{
				MaxOpenConns: -1,
			},
			wantErr: "max_open_conns",
		},
		{
			name: "idle connections exceed open connections",
			input: config.MySQLPoolConfig{
				MaxOpenConns: 5,
				MaxIdleConns: 10,
			},
			wantErr: "max_idle_conns cannot exceed max_open_conns",
		},
		{
			name: "negative max lifetime",
			input: config.MySQLPoolConfig{
				ConnMaxLifetime: -time.Second,
			},
			wantErr: "conn_max_lifetime",
		},
		{
			name: "negative max idle time",
			input: config.MySQLPoolConfig{
				ConnMaxIdleTime: -time.Second,
			},
			wantErr: "conn_max_idle_time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeMySQLPoolConfig(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected config: got %+v, want %+v", got, tt.want)
			}
		})
	}
}
