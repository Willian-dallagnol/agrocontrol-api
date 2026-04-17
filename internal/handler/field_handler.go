package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FieldHandler struct {
	Service *service.FieldService
}

func NewFieldHandler(service *service.FieldService) *FieldHandler {
	return &FieldHandler{Service: service}
}

func (h *FieldHandler) CreateField(c *gin.Context) {
	var req dto.CreateFieldRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	field, err := h.Service.CreateField(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, field)
}

func (h *FieldHandler) GetFields(c *gin.Context) {
	fields, err := h.Service.GetFields()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fields)
}

func (h *FieldHandler) GetFieldByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	field, err := h.Service.GetFieldByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "talhão não encontrado"})
		return
	}

	c.JSON(http.StatusOK, field)
}

func (h *FieldHandler) UpdateField(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	field, err := h.Service.UpdateField(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "talhão não encontrado"})
		return
	}

	c.JSON(http.StatusOK, field)
}

func (h *FieldHandler) DeleteField(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	err = h.Service.DeleteField(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "talhão não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "talhão removido com sucesso"})
}

func (h *FieldHandler) GetFieldsByFarm(c *gin.Context) {
	farmIDParam := c.Param("id")

	farmID, err := strconv.ParseUint(farmIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	fields, err := h.Service.GetFieldsByFarmID(uint(farmID))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, fields)
}
