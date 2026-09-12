package hello

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizhello "auth_info/internal/biz/hello"
)

// Handler Hello HTTP 请求处理器
type Handler struct {
	uc *bizhello.UseCase
}

// NewHandler Wire Provider
func NewHandler(uc *bizhello.UseCase) *Handler {
	return &Handler{uc: uc}
}

// Hello 返回 Hello 消息
func (h *Handler) Hello(c *gin.Context) {
	req := apipb.HelloRequest{Name: strings.TrimSpace(c.Query("name"))}
	msg := h.uc.SayHello(c.Request.Context(), req.GetName())

	c.JSON(http.StatusOK, &apipb.HelloReply{
		Code:    http.StatusOK,
		Message: "success",
		Data: &apipb.HelloData{
			Message: msg,
		},
	})
}
