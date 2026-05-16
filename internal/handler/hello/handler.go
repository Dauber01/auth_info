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

// Hello
// @Summary Hello 示例接口
// @Description 返回 Hello 消息
// @Tags Hello
// @Accept json
// @Produce json
// @Param name query string false "名字"
// @Success 200 {object} apipb.HelloReply "请求成功"
// @Failure 401 {object} apipb.OperationReply "未认证"
// @Failure 403 {object} apipb.OperationReply "无访问权限"
// @Failure 500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router /hello [get]
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
