package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"auth_info/internal/pkg/logger"
)

// AccessLog emits request metadata only. Error details belong to ErrorHandler.
func AccessLog(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.WithContext(log, c.Request.Context()).Info("http request completed",
			zap.String("method", c.Request.Method), zap.String("route", c.FullPath()),
			zap.Int("status", c.Writer.Status()), zap.Duration("duration", time.Since(start)),
		)
	}
}
