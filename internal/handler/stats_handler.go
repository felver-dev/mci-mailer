package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
)

type StatsHandler struct {
	repo repository.StatsRepository
}

func NewStatsHandler(repo repository.StatsRepository) *StatsHandler {
	return &StatsHandler{repo: repo}
}

// GET /portal/stats
func (h *StatsHandler) Get(c *gin.Context) {
	overview, err := h.repo.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch stats"))
		return
	}

	byApp, err := h.repo.GetPerApp(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch per-app stats"))
		return
	}

	c.JSON(http.StatusOK, dto.OK("", dto.StatsResponse{
		Overview: *overview,
		ByApp:    byApp,
	}))
}
