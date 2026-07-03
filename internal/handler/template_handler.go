package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/service"
)

type TemplateHandler struct {
	svc *service.TemplateService
}

func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.OK("template created", resp))
}

func (h *TemplateHandler) List(c *gin.Context) {
	templates, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch templates"))
		return
	}
	c.JSON(http.StatusOK, dto.OK("", templates))
}

func (h *TemplateHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), name, req)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.OK("template updated", resp))
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusNotFound, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.OK("template deleted", nil))
}
