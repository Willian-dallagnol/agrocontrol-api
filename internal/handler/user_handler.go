package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	user, err := h.Service.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}
