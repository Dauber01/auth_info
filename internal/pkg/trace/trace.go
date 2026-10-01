// Package trace carries a bounded request identifier without transport dependencies.
package trace

import (
	"context"
	"crypto/rand"
)

// Header is the HTTP header and gRPC metadata key used for correlation.
const Header = "X-Trace-ID"

type contextKey struct{}

// Normalize preserves safe incoming identifiers or creates a fresh one.
func Normalize(id string) string {
	if len(id) > 0 && len(id) <= 128 {
		valid := true
		for _, c := range id {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
				c == '-' || c == '_' || c == '.') {
				valid = false
				break
			}
		}
		if valid {
			return id
		}
	}
	return rand.Text()
}

// WithID returns a context containing the request identifier.
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// ID retrieves the identifier, or an empty string outside a request.
func ID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
