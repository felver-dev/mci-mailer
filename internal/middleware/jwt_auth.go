package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/service"
)

const ContextKeyJWT = "jwt_claims"

func JWTAuth(userSvc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("missing or invalid Authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := userSvc.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("invalid or expired token"))
			return
		}

		c.Set(ContextKeyJWT, claims)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetJWTClaims(c)
		if claims == nil || claims.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Err("admin role required"))
			return
		}
		c.Next()
	}
}

func GetJWTClaims(c *gin.Context) *service.JWTClaims {
	v, _ := c.Get(ContextKeyJWT)
	claims, _ := v.(*service.JWTClaims)
	return claims
}
