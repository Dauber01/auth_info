package shared_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	bizauth "auth_info/internal/biz/auth"
	bizdict "auth_info/internal/biz/dict"
	dataauth "auth_info/internal/data/auth"
	datadict "auth_info/internal/data/dict"
	"auth_info/internal/data/shared"
)

func transactionDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "fixture.db")),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Error(err)
		}
	})
	if err := db.AutoMigrate(&dataauth.User{}, &datadict.DictType{}, &datadict.DictItem{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCrossRepositoryCommitAndRollback(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{true: "rollback", false: "commit"}[fail], func(t *testing.T) {
			db := transactionDB(t)
			users := dataauth.NewUserRepository(db)
			dicts := datadict.NewDictRepository(db)
			expected := errors.New("second operation failed")
			err := shared.NewTxScope(db).Do(context.Background(), func(ctx context.Context) error {
				userToCreate := &bizauth.User{Username: "transaction-user", Password: "fixture", Role: "user"}
				if err := users.Create(ctx, userToCreate); err != nil {
					return err
				}
				if err := dicts.CreateDictType(ctx, &bizdict.DictType{Code: "transaction-type", Name: "fixture"}); err != nil {
					return err
				}
				// Reads must participate in the same transaction as both writes.
				user, err := users.GetByUsername(ctx, "transaction-user")
				if err != nil {
					return err
				}
				if user == nil {
					return errors.New("transaction read missed write")
				}
				if fail {
					return expected
				}
				return nil
			})
			if fail && !errors.Is(err, expected) || !fail && err != nil {
				t.Fatal(err)
			}
			user, err := users.GetByUsername(context.Background(), "transaction-user")
			if err != nil {
				t.Fatal(err)
			}
			dict, err := dicts.GetDictTypeByCode(context.Background(), "transaction-type")
			if err != nil {
				t.Fatal(err)
			}
			if (user != nil) == fail || (dict != nil) == fail {
				t.Fatal("transaction atomicity violated")
			}
		})
	}
}

func TestTransactionRejectsNestedForeignPoolAndCanceledContext(t *testing.T) {
	db := transactionDB(t)
	other := transactionDB(t)
	scope := shared.NewTxScope(db)
	err := scope.Do(context.Background(), func(ctx context.Context) error {
		if err := scope.Do(ctx, func(context.Context) error { return nil }); !errors.Is(err, shared.ErrNestedTransaction) {
			t.Fatal(err)
		}
		_, err := dataauth.NewUserRepository(other).GetByUsername(ctx, "any")
		if err == nil {
			t.Fatal("foreign pool accepted")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = scope.Do(ctx, func(context.Context) error {
		t.Fatal("canceled callback ran")
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestCancellationDuringTransactionRollsBack(t *testing.T) {
	db := transactionDB(t)
	users := dataauth.NewUserRepository(db)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := shared.NewTxScope(db).Do(ctx, func(ctx context.Context) error {
		if err := users.Create(ctx, &bizauth.User{Username: "cancelled", Password: "fixture", Role: "user"}); err != nil {
			return err
		}
		cancel()
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	user, err := users.GetByUsername(context.Background(), "cancelled")
	if err != nil {
		t.Fatal(err)
	}
	if user != nil {
		t.Fatal("canceled transaction committed its write")
	}
}
