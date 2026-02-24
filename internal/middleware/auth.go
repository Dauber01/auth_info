package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"

	"auth_info/internal/biz"
)

const claimsKey = "claims"

// JWTAuth JWT 鉴权中间件，验证 Bearer Token 并将 Claims 写入上下文
func JWTAuth(uc *biz.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "missing or invalid Authorization header",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := uc.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "invalid token",
			})
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}

// CasbinAuth Casbin 权限校验中间件，在 JWT 中间件之后使用
func CasbinAuth(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(claimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			return
		}

		authClaims, ok := claims.(*biz.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			return
		}

		allowed, err := enforcer.Enforce(authClaims.Role, c.FullPath(), c.Request.Method)
		if err != nil || !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "permission denied",
			})
			return
		}

		c.Next()
	}
}

// GetClaims 从 gin.Context 中取出 JWT Claims（供 handler 使用）
func GetClaims(c *gin.Context) *biz.Claims {
	v, _ := c.Get(claimsKey)
	claims, _ := v.(*biz.Claims)
	return claims
}
