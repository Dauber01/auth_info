package middleware

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/pkg/trace"
)

// TraceID carries the same identifier through response headers and business context.
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := trace.Normalize(c.GetHeader(trace.Header))
		c.Header(trace.Header, id)
		c.Request = c.Request.WithContext(trace.WithID(c.Request.Context(), id))
		c.Next()
	}
}
