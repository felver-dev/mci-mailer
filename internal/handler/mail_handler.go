package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/middleware"
	"github.com/mcicare/mci-mailer/internal/service"
)

type MailHandler struct {
	mailer *service.MailerService
}

func NewMailHandler(mailer *service.MailerService) *MailHandler {
	return &MailHandler{mailer: mailer}
}

func (h *MailHandler) Send(c *gin.Context) {
	var req dto.SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	apiKey := middleware.GetApiKey(c)
	resp, err := h.mailer.Send(c.Request.Context(), req, &apiKey.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"log_id":  resp.LogID,
			"status":  resp.Status,
			"message": resp.Message,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, dto.OK("email queued", resp))
}
