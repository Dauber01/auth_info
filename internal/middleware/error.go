package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"auth_info/internal/apperr"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		log.Error(
			"request failed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("error", err.Error()),
		)

		if c.Writer.Written() {
			return
		}

		status := apperr.HTTPStatus(err)
		c.AbortWithStatusJSON(status, ErrorResponse{
			Code:    status,
			Message: apperr.Message(err),
		})
	}
}
