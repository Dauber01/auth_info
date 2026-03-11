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
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := validateRegisterRequest(&req); err != nil {
		badRequest(c, err)
		return
	}

	if err := h.uc.Register(req.GetUsername(), req.GetPassword()); err != nil {
		writeOperationReply(c, http.StatusConflict, err.Error())
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
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := validateLoginRequest(&req); err != nil {
		badRequest(c, err)
		return
	}

	token, err := h.uc.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		writeOperationReply(c, http.StatusUnauthorized, err.Error())
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
