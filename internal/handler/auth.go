package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizauth "auth_info/internal/biz/auth"
)

// AuthHandler 登录/注册 HTTP 处理器
type AuthHandler struct {
	uc *bizauth.UseCase
}

// NewAuthHandler Wire Provider
func NewAuthHandler(uc *bizauth.UseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Register
// @Summary 注册
// @Tags Auth
// @Accept json
// @Produce json
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req apipb.RegisterRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	if err := h.uc.Register(c.Request.Context(), req.GetUsername(), req.GetPassword()); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "registered successfully")
}

// Login
// @Summary 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req apipb.LoginRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	token, err := h.uc.Login(c.Request.Context(), req.GetUsername(), req.GetPassword())
	if err != nil {
		writeError(c, err)
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
