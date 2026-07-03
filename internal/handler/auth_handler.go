package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/service"
)

type AuthHandler struct {
	userSvc     *service.UserService
	masterToken string
}

func NewAuthHandler(userSvc *service.UserService, masterToken string) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, masterToken: masterToken}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	resp, err := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.OK("login successful", resp))
}

// POST /setup/bootstrap  (requires X-Master-Token header)
func (h *AuthHandler) Bootstrap(c *gin.Context) {
	if c.GetHeader("X-Master-Token") != h.masterToken {
		c.JSON(http.StatusUnauthorized, dto.Err("invalid master token"))
		return
	}

	var req dto.BootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}

	resp, err := h.userSvc.Bootstrap(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.OK("admin account created", resp))
}
