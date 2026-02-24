package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"auth_info/internal/biz"
)

// AuthHandler 登录/注册 HTTP 处理器
type AuthHandler struct {
	uc *biz.AuthUseCase
}

// NewAuthHandler Wire Provider
func NewAuthHandler(uc *biz.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register
// @Summary 注册
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body registerRequest true "注册信息"
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if err := h.uc.Register(req.Username, req.Password); err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "registered successfully"})
}

// Login
// @Summary 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "登录信息"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	token, err := h.uc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{"token": token},
	})
}
