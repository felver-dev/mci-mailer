package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"golang.org/x/time/rate"
)

type keyLimiter struct {
	limiter *rate.Limiter
}

var (
	limiters sync.Map
)

func getLimiter(keyID string, rps float64, burst int) *rate.Limiter {
	v, _ := limiters.LoadOrStore(keyID, &keyLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	})
	return v.(*keyLimiter).limiter
}

func RateLimit(requestsPerSecond float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := GetApiKey(c)
		if apiKey == nil {
			c.Next()
			return
		}

		limiter := getLimiter(apiKey.ID.String(), requestsPerSecond, burst)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.Err("rate limit exceeded"))
			return
		}
		c.Next()
	}
}
