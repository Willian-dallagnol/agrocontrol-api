package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 🌱 Handler responsável pelas rotas de Field (talhões)
type FieldHandler struct {
	Service *service.FieldService
	// 👉 referência ao service onde está a regra de negócio
}

// 🏗️ Construtor do handler
func NewFieldHandler(service *service.FieldService) *FieldHandler {
	return &FieldHandler{Service: service}
}

// 🚀 Criar novo talhão
func (h *FieldHandler) CreateField(c *gin.Context) {
	var req dto.CreateFieldRequest
	// 👉 recebe dados do body

	// 📥 valida JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		// ❌ erro de validação
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	// 🧠 chama o service
	field, err := h.Service.CreateField(req)
	if err != nil {
		// ❌ erro de regra (ex: Farm não existe)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusCreated, field)
}

// 📋 Listar todos os talhões
func (h *FieldHandler) GetFields(c *gin.Context) {

	// 🧠 busca no service
	fields, err := h.Service.GetFields()
	if err != nil {
		// ❌ erro interno
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ retorna lista
	c.JSON(http.StatusOK, fields)
}

// 🔍 Buscar talhão por ID
func (h *FieldHandler) GetFieldByID(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// ❌ ID inválido
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// 🧠 busca no service
	field, err := h.Service.GetFieldByID(uint(id))
	if err != nil {
		// ❌ não encontrado
		c.JSON(http.StatusNotFound, gin.H{"error": "talhão não encontrado"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, field)
}

// 🔄 Atualizar talhão
func (h *FieldHandler) UpdateField(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateFieldRequest
	// 👉 recebe dados do body

	// 📥 valida JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	// 🧠 chama service
	field, err := h.Service.UpdateField(uint(id), req)
	if err != nil {
		// ❌ não encontrado ou erro de regra
		c.JSON(http.StatusNotFound, gin.H{"error": "talhão não encontrado"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, field)
}

// 🗑️ Deletar talhão
func (h *FieldHandler) DeleteField(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// 🧠 chama service
	err = h.Service.DeleteField(uint(id))
	if err != nil {
		// ❌ não encontrado
		c.JSON(http.StatusNotFound, gin.H{"error": "talhão não encontrado"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, gin.H{"message": "talhão removido com sucesso"})
}

// 🔗 Listar talhões de uma fazenda específica
func (h *FieldHandler) GetFieldsByFarm(c *gin.Context) {

	// 📥 pega o ID da fazenda da URL
	farmIDParam := c.Param("id")

	// 🔄 converte string → uint
	farmID, err := strconv.ParseUint(farmIDParam, 10, 64)
	if err != nil {
		// ❌ ID inválido
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	// 🧠 busca todos os talhões da fazenda
	fields, err := h.Service.GetFieldsByFarmID(uint(farmID))
	if err != nil {
		// ❌ fazenda não encontrada ou erro
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	// ✅ sucesso
	c.JSON(200, fields)
}
