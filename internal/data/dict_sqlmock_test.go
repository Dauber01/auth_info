package data

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newDictRepoWithSQLMock(t *testing.T) (*DictRepo, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		_ = sqlDB.Close()
		t.Fatalf("open gorm with sqlmock: %v", err)
	}

	cleanup := func() {
		mock.ExpectClose()
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("close sqlmock db: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
	}

	return NewDictRepository(gdb), mock, cleanup
}

func TestDictRepo_GetDictTypeByCode(t *testing.T) {
	now := time.Now()
	query := regexp.QuoteMeta(
		"SELECT * FROM `dict_types` WHERE code = ? AND `dict_types`.`deleted_at` IS NULL ORDER BY `dict_types`.`id` LIMIT ?",
	)
	dbErr := errors.New("db down")

	tests := []struct {
		name    string
		mockFn  func(sqlmock.Sqlmock)
		want    *DictType
		wantErr error
	}{
		{
			name: "success",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"code", "name", "description", "sort",
				}).AddRow(1, now, now, nil, "gender", "Gender", "gender dict", 1)
				mock.ExpectQuery(query).WithArgs("gender", 1).WillReturnRows(rows)
			},
			want: &DictType{Model: gorm.Model{ID: 1}, Code: "gender", Name: "Gender", Description: "gender dict", Sort: 1},
		},
		{
			name: "not_found",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WithArgs("gender", 1).WillReturnError(gorm.ErrRecordNotFound)
			},
			want: nil,
		},
		{
			name: "db_error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WithArgs("gender", 1).WillReturnError(dbErr)
			},
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, cleanup := newDictRepoWithSQLMock(t)
			defer cleanup()

			tt.mockFn(mock)

			got, err := repo.GetDictTypeByCode(context.Background(), "gender")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected dict type, got nil")
			}
			if got.ID != tt.want.ID || got.Code != tt.want.Code || got.Name != tt.want.Name || got.Description != tt.want.Description || got.Sort != tt.want.Sort {
				t.Fatalf("unexpected dict type: %+v", got)
			}
		})
	}
}
