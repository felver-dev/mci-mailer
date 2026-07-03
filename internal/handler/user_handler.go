package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// POST /portal/users  (admin only)
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Err(err.Error()))
		return
	}
	user, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, dto.Err(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.OK("user created", user))
}

// GET /portal/users  (admin only)
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Err("failed to fetch users"))
		return
	}
	c.JSON(http.StatusOK, dto.OK("", users))
}
