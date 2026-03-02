package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"auth_info/internal/biz"
)

// HelloHandler HTTP 请求处理器
type HelloHandler struct {
	uc *biz.HelloUseCase
}

// NewHelloHandler Wire Provider
func NewHelloHandler(uc *biz.HelloUseCase) *HelloHandler {
	return &HelloHandler{uc: uc}
}

// Hello
// @Summary Hello 示例接口
// @Description 返回 Hello 消息
// @Tags Hello
// @Accept json
// @Produce json
// @Param name query string false "名字"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hello [get]
func (h *HelloHandler) Hello(c *gin.Context) {
	name := c.Query("name")
	msg := h.uc.SayHello(name)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"message": msg,
		},
	})
}
