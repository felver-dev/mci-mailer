package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		fmt.Printf("[MCM] %s | %3d | %12v | %s | %s %s\n",
			time.Now().Format("2006/01/02 15:04:05"),
			status,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
