package middleware

import (
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"

	"auth_info/internal/apperr"
	bizauth "auth_info/internal/biz/auth"
)

const claimsKey = "claims"

// JWTAuth JWT 鉴权中间件，验证 Bearer Token 并将 Claims 写入上下文
func JWTAuth(uc *bizauth.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			_ = c.Error(apperr.New(apperr.CodeUnauthenticated, "missing or invalid Authorization header"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := uc.ParseToken(tokenStr)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
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
			_ = c.Error(apperr.New(apperr.CodeUnauthenticated, "unauthorized"))
			c.Abort()
			return
		}

		authClaims, ok := claims.(*bizauth.Claims)
		if !ok {
			_ = c.Error(apperr.New(apperr.CodeUnauthenticated, "unauthorized"))
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(authClaims.Role, c.FullPath(), c.Request.Method)
		if err != nil {
			_ = c.Error(apperr.Wrap(apperr.CodeInternal, "failed to enforce policy", err))
			c.Abort()
			return
		}
		if !allowed {
			_ = c.Error(apperr.New(apperr.CodePermissionDenied, "permission denied"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetClaims 从 gin.Context 中取出 JWT Claims（供 handler 使用）
func GetClaims(c *gin.Context) *bizauth.Claims {
	v, _ := c.Get(claimsKey)
	claims, _ := v.(*bizauth.Claims)
	return claims
}
