package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
)

type LogHandler struct {
	repo repository.EmailLogRepository
}

func NewLogHandler(repo repository.EmailLogRepository) *LogHandler {
	return &LogHandler{repo: repo}
}

func (h *LogHandler) List(c *gin.Context) {
	var filter dto.LogFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	logs, total, err := h.repo.FindAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch logs"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
		"data":      toLogResponses(logs),
	})
}

func toLogResponses(logs []domain.EmailLog) []gin.H {
	result := make([]gin.H, len(logs))
	for i, l := range logs {
		entry := gin.H{
			"id":            l.ID.String(),
			"from":          l.FromAddress,
			"to":            l.ToAddresses,
			"subject":       l.Subject,
			"status":        l.Status,
			"attempts":      l.Attempts,
			"created_at":    l.CreatedAt.Format(time.RFC3339),
		}
		if l.ApiKeyID != nil {
			entry["api_key_id"] = l.ApiKeyID.String()
		}
		if l.TemplateName != nil {
			entry["template"] = *l.TemplateName
		}
		if l.ErrorMsg != nil {
			entry["error"] = *l.ErrorMsg
		}
		if l.SentAt != nil {
			entry["sent_at"] = l.SentAt.Format(time.RFC3339)
		}
		result[i] = entry
	}
	return result
}
