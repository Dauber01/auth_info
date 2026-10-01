package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestTimeout requests cooperative cancellation; it never runs Gin in a background goroutine.
// ErrorHandler must wrap it so errors raised after the handler returns are rendered once.
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout <= 0 {
			c.Next()
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		if err := ctx.Err(); err != nil && !c.Writer.Written() && len(c.Errors) == 0 {
			_ = c.Error(err) // ErrorHandler owns logging and rendering.
			c.Abort()
		}
	}
}
