package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/middleware"
	"github.com/mcicare/mci-mailer/internal/repository"
)

type PortalAppHandler struct {
	repo repository.ApiKeyRepository
}

func NewPortalAppHandler(repo repository.ApiKeyRepository) *PortalAppHandler {
	return &PortalAppHandler{repo: repo}
}

// POST /portal/apps
func (h *PortalAppHandler) Create(c *gin.Context) {
	var req dto.CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	rawKey, keyHash, err := domain.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to generate API key"))
		return
	}

	claims := middleware.GetJWTClaims(c)
	userID, _ := uuid.Parse(claims.UserID)

	key := &domain.ApiKey{
		ID:              uuid.New(),
		Name:            req.Name,
		KeyHash:         keyHash,
		Scopes:          req.Scopes,
		IsActive:        true,
		CreatedAt:       time.Now(),
		CreatedByUserID: &userID,
	}

	if err := h.repo.Create(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to create app"))
		return
	}

	c.JSON(http.StatusCreated, dto.OK("app created — save this key, it will not be shown again", dto.CreateApiKeyResponse{
		ID:     key.ID.String(),
		Name:   key.Name,
		Key:    rawKey,
		Scopes: key.Scopes,
	}))
}

// GET /portal/apps
func (h *PortalAppHandler) List(c *gin.Context) {
	keys, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch apps"))
		return
	}

	type appItem struct {
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Scopes        []string   `json:"scopes"`
		IsActive      bool       `json:"is_active"`
		CreatedAt     string     `json:"created_at"`
		LastUsedAt    *string    `json:"last_used_at,omitempty"`
		CreatedBy     *string    `json:"created_by,omitempty"`
	}

	result := make([]appItem, len(keys))
	for i, k := range keys {
		item := appItem{
			ID:        k.ID.String(),
			Name:      k.Name,
			Scopes:    k.Scopes,
			IsActive:  k.IsActive,
			CreatedAt: k.CreatedAt.UTC().Format(time.RFC3339),
			CreatedBy: k.CreatedByUserName,
		}
		if k.LastUsedAt != nil {
			s := k.LastUsedAt.UTC().Format(time.RFC3339)
			item.LastUsedAt = &s
		}
		result[i] = item
	}

	c.JSON(http.StatusOK, dto.OK("", result))
}

// DELETE /portal/apps/:id
func (h *PortalAppHandler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err("invalid app id"))
		return
	}
	if err := h.repo.Revoke(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to revoke app"))
		return
	}
	c.JSON(http.StatusOK, dto.OK("app revoked", nil))
}
