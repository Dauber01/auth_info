package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"auth_info/internal/logger"
)

type HelloHandler struct {
}

// NewHelloHandler 创建 HelloHandler 实例
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// Hello
// @Summary Hello 示例接口
// @Description 返回 Hello 消息
// @Tags Hello
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /hello [get]
func (h *HelloHandler) Hello(c *gin.Context) {
	logger.GetLogger().Info("Hello handler called", zap.String("path", c.Request.URL.Path))

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"message": "Hello, World!",
		},
	})
}
