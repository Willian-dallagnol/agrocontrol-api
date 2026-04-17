package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FarmHandler struct {
	Service *service.FarmService
}

func NewFarmHandler(service *service.FarmService) *FarmHandler {
	return &FarmHandler{Service: service}
}

func (h *FarmHandler) CreateFarm(c *gin.Context) {
	var req dto.CreateFarmRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	userID := c.GetUint("user_id")

	farm, err := h.Service.CreateFarm(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, farm)
}

func (h *FarmHandler) GetFarms(c *gin.Context) {
	farms, err := h.Service.GetFarms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, farms)
}

func (h *FarmHandler) GetFarmByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	farm, err := h.Service.GetFarmByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fazenda não encontrada"})
		return
	}

	c.JSON(http.StatusOK, farm)
}

func (h *FarmHandler) UpdateFarm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateFarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	farm, err := h.Service.UpdateFarm(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fazenda não encontrada"})
		return
	}

	c.JSON(http.StatusOK, farm)
}

func (h *FarmHandler) DeleteFarm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	err = h.Service.DeleteFarm(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fazenda não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "fazenda removida com sucesso"})
}
