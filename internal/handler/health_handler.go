package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	smtpClient "github.com/mcicare/mci-mailer/internal/smtp"
)

type HealthHandler struct {
	smtp *smtpClient.Client
}

func NewHealthHandler(smtp *smtpClient.Client) *HealthHandler {
	return &HealthHandler{smtp: smtp}
}

func (h *HealthHandler) Check(c *gin.Context) {
	smtpOK := true
	smtpErr := ""
	if err := h.smtp.Ping(); err != nil {
		smtpOK = false
		smtpErr = err.Error()
	}

	status := http.StatusOK
	if !smtpOK {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"status":    map[bool]string{true: "ok", false: "degraded"}[smtpOK],
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"smtp": gin.H{
			"ok":    smtpOK,
			"error": smtpErr,
		},
	})
}
