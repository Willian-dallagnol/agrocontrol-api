package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CropHandler struct {
	Service *service.CropService
}

func NewCropHandler(service *service.CropService) *CropHandler {
	return &CropHandler{Service: service}
}

func (h *CropHandler) CreateCrop(c *gin.Context) {
	var req dto.CreateCropRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	crop, err := h.Service.CreateCrop(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, crop)
}

func (h *CropHandler) GetCrops(c *gin.Context) {
	crops, err := h.Service.GetCrops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, crops)
}

func (h *CropHandler) GetCropByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	crop, err := h.Service.GetCropByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cultura não encontrada"})
		return
	}

	c.JSON(http.StatusOK, crop)
}

func (h *CropHandler) UpdateCrop(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateCropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	crop, err := h.Service.UpdateCrop(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, crop)
}

func (h *CropHandler) DeleteCrop(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	err = h.Service.DeleteCrop(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cultura não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cultura removida com sucesso"})
}
