package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizhello "auth_info/internal/biz/hello"
)

// HelloHandler HTTP 请求处理器
type HelloHandler struct {
	uc *bizhello.UseCase
}

// NewHelloHandler Wire Provider
func NewHelloHandler(uc *bizhello.UseCase) *HelloHandler {
	return &HelloHandler{uc: uc}
}

// Hello
// @Summary Hello 示例接口
// @Description 返回 Hello 消息
// @Tags Hello
// @Accept json
// @Produce json
// @Param name query string false "名字"
// @Success 200 {object} apipb.HelloReply
// @Failure 401 {object} apipb.OperationReply
// @Security BearerAuth
// @Router /hello [get]
func (h *HelloHandler) Hello(c *gin.Context) {
	req := apipb.HelloRequest{Name: strings.TrimSpace(c.Query("name"))}
	msg := h.uc.SayHello(req.GetName())

	c.JSON(http.StatusOK, &apipb.HelloReply{
		Code:    http.StatusOK,
		Message: "success",
		Data: &apipb.HelloData{
			Message: msg,
		},
	})
}
