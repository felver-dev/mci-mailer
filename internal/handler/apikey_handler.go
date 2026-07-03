package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/service"
)

type ApiKeyHandler struct {
	svc *service.ApiKeyService
}

func NewApiKeyHandler(svc *service.ApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{svc: svc}
}

func (h *ApiKeyHandler) Create(c *gin.Context) {
	var req dto.CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.OK("API key created — save this key, it will not be shown again", resp))
}

func (h *ApiKeyHandler) List(c *gin.Context) {
	keys, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch API keys"))
		return
	}
	c.JSON(http.StatusOK, dto.OK("", keys))
}

func (h *ApiKeyHandler) Revoke(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Revoke(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.OK("API key revoked", nil))
}
