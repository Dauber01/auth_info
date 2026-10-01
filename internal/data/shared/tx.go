// Package shared provides transaction context for repositories sharing a database pool.
package shared

import (
	"context"
	"errors"

	"gorm.io/gorm"

	bizshared "auth_info/internal/biz/shared"
)

// ErrNestedTransaction prevents accidental independent commits inside an outer unit of work.
var ErrNestedTransaction = errors.New("nested transaction scopes are not supported")

type transactionKey struct{}
type transaction struct {
	pool gorm.ConnPool
	db   *gorm.DB
}
type scope struct{ db *gorm.DB }

// NewTxScope constructs the implementation of the business-side transaction contract.
func NewTxScope(db *gorm.DB) bizshared.TxScope { return &scope{db: db} }

func (s *scope) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, ok := ctx.Value(transactionKey{}).(transaction); ok {
		return ErrNestedTransaction
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionKey{}, transaction{pool: s.db.ConnPool, db: tx})
		if err := fn(txCtx); err != nil {
			return err
		}
		return txCtx.Err()
	})
}

// DB selects the current transaction only for its owning pool, otherwise the normal connection.
// A transaction from another pool is rejected rather than silently committing unrelated writes.
func DB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey{}).(transaction); ok {
		if tx.pool != db.ConnPool {
			failed := db.Session(&gorm.Session{NewDB: true}).WithContext(ctx)
			failed.AddError(errors.New("repository uses a different transaction pool"))
			return failed
		}
		return tx.db.WithContext(ctx)
	}
	return db.WithContext(ctx)
}
