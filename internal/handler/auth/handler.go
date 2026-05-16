package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizauth "auth_info/internal/biz/auth"
	"auth_info/internal/handler/httpx"
)

// Handler 登录/注册 HTTP 处理器
type Handler struct {
	uc *bizauth.UseCase
}

// NewHandler Wire Provider
func NewHandler(uc *bizauth.UseCase) *Handler {
	return &Handler{uc: uc}
}

// Register
// @Summary 注册
// @Description 公开接口：使用用户名和密码创建新用户。
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body apipb.RegisterRequest true "注册参数"
// @Success 200 {object} apipb.OperationReply "注册成功"
// @Failure 400 {object} apipb.OperationReply "请求参数错误"
// @Failure 409 {object} apipb.OperationReply "用户名已存在"
// @Failure 500 {object} apipb.OperationReply "服务器内部错误"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req apipb.RegisterRequest
	if !httpx.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.uc.Register(c.Request.Context(), req.GetUsername(), req.GetPassword()); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "registered successfully")
}

// Login
// @Summary 登录
// @Description 公开接口：使用用户名和密码登录，成功后返回 JWT Token。
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body apipb.LoginRequest true "登录参数"
// @Success 200 {object} apipb.LoginReply "登录成功"
// @Failure 400 {object} apipb.OperationReply "请求参数错误"
// @Failure 401 {object} apipb.OperationReply "用户名或密码错误"
// @Failure 500 {object} apipb.OperationReply "服务器内部错误"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req apipb.LoginRequest
	if !httpx.BindAndValidateJSON(c, &req) {
		return
	}

	token, err := h.uc.Login(c.Request.Context(), req.GetUsername(), req.GetPassword())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, &apipb.LoginReply{
		Code:    http.StatusOK,
		Message: "success",
		Data: &apipb.LoginData{
			Token: token,
		},
	})
}
