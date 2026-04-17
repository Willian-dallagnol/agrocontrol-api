package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 🌾 Handler responsável pelas rotas de Crop (culturas)
type CropHandler struct {
	Service *service.CropService
	// 👉 referência ao service onde está a regra de negócio
}

// 🏗️ Construtor do handler
func NewCropHandler(service *service.CropService) *CropHandler {
	return &CropHandler{Service: service}
}

// 🚀 Criar nova cultura
func (h *CropHandler) CreateCrop(c *gin.Context) {
	var req dto.CreateCropRequest
	// 👉 recebe os dados enviados no body

	// 📥 valida e faz bind do JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		// ❌ erro de validação (campos obrigatórios, formato)
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	// 🧠 chama o service para criar a cultura
	crop, err := h.Service.CreateCrop(req)
	if err != nil {
		// ❌ erro de regra de negócio (ex: Field não existe)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusCreated, crop)
}

// 📋 Listar todas as culturas
func (h *CropHandler) GetCrops(c *gin.Context) {

	// 🧠 busca todas as culturas no service
	crops, err := h.Service.GetCrops()
	if err != nil {
		// ❌ erro interno
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ retorna lista
	c.JSON(http.StatusOK, crops)
}

// 🔍 Buscar cultura por ID
func (h *CropHandler) GetCropByID(c *gin.Context) {

	// 📥 pega o ID da URL e converte para int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// ❌ ID inválido
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// 🧠 busca no service
	crop, err := h.Service.GetCropByID(uint(id))
	if err != nil {
		// ❌ não encontrado
		c.JSON(http.StatusNotFound, gin.H{"error": "cultura não encontrada"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, crop)
}

// 🔄 Atualizar cultura
func (h *CropHandler) UpdateCrop(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateCropRequest
	// 👉 recebe dados do body

	// 📥 valida JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	// 🧠 chama service para atualizar
	crop, err := h.Service.UpdateCrop(uint(id), req)
	if err != nil {
		// ❌ erro de regra (ex: Field inválido)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, crop)
}

// 🗑️ Deletar cultura
func (h *CropHandler) DeleteCrop(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// 🧠 chama service para deletar
	err = h.Service.DeleteCrop(uint(id))
	if err != nil {
		// ❌ não encontrado
		c.JSON(http.StatusNotFound, gin.H{"error": "cultura não encontrada"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, gin.H{"message": "cultura removida com sucesso"})
}
