package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"auth_info/internal/pkg/logger"
)

// Recovery reports panics once and keeps the existing numeric error envelope.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				logger.WithContext(log, c.Request.Context()).Error("request panicked", zap.ByteString("stack", debug.Stack()))
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
						Code: http.StatusInternalServerError, Message: "internal server error",
					})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
