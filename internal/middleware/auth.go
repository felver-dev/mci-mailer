package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
)

const ContextKeyApiKey = "api_key"

func Auth(repo repository.ApiKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("missing or invalid Authorization header"))
			return
		}

		rawKey := strings.TrimPrefix(header, "Bearer ")
		if !strings.HasPrefix(rawKey, "MCM.") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("invalid API key format"))
			return
		}

		keyHash := domain.HashAPIKey(rawKey)
		apiKey, err := repo.FindByHash(c.Request.Context(), keyHash)
		if err != nil || apiKey == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("invalid or revoked API key"))
			return
		}

		go func() {
			_ = repo.UpdateLastUsed(c.Request.Context(), apiKey.ID)
		}()

		c.Set(ContextKeyApiKey, apiKey)
		c.Next()
	}
}

func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, exists := c.Get(ContextKeyApiKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Err("unauthorized"))
			return
		}
		apiKey, ok := key.(*domain.ApiKey)
		if !ok || !apiKey.HasScope(scope) {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Err("insufficient scope: "+scope+" required"))
			return
		}
		c.Next()
	}
}

func GetApiKey(c *gin.Context) *domain.ApiKey {
	key, _ := c.Get(ContextKeyApiKey)
	apiKey, _ := key.(*domain.ApiKey)
	return apiKey
}
