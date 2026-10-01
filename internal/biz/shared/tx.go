// Package shared contains cross-use-case contracts without storage dependencies.
package shared

import "context"

// TxScope executes participating repositories atomically. Nested scopes are rejected.
// The callback context must be passed to every repository; it must not escape the callback.
type TxScope interface {
	Do(context.Context, func(context.Context) error) error
}
