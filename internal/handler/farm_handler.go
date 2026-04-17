package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 🚜 Handler responsável pelas rotas de Farm (fazendas)
type FarmHandler struct {
	Service *service.FarmService
	// 👉 referência ao service onde está a regra de negócio
}

// 🏗️ Construtor do handler
func NewFarmHandler(service *service.FarmService) *FarmHandler {
	return &FarmHandler{Service: service}
}

// 🚀 Criar nova fazenda
func (h *FarmHandler) CreateFarm(c *gin.Context) {
	var req dto.CreateFarmRequest
	// 👉 estrutura que recebe dados do body

	// 📥 valida e faz bind do JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		// ❌ erro de validação
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	// 🔐 pega o ID do usuário vindo do token (middleware JWT)
	userID := c.GetUint("user_id")

	// 🧠 chama o service para criar a fazenda
	farm, err := h.Service.CreateFarm(req, userID)
	if err != nil {
		// ❌ erro de regra de negócio
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusCreated, farm)
}

// 📋 Listar todas as fazendas
func (h *FarmHandler) GetFarms(c *gin.Context) {

	// 🧠 busca todas as fazendas no service
	farms, err := h.Service.GetFarms()
	if err != nil {
		// ❌ erro interno
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ retorna lista
	c.JSON(http.StatusOK, farms)
}

// 🔍 Buscar fazenda por ID
func (h *FarmHandler) GetFarmByID(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// ❌ ID inválido
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// 🧠 busca no service
	farm, err := h.Service.GetFarmByID(uint(id))
	if err != nil {
		// ❌ não encontrada
		c.JSON(http.StatusNotFound, gin.H{"error": "fazenda não encontrada"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, farm)
}

// 🔄 Atualizar fazenda
func (h *FarmHandler) UpdateFarm(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dto.UpdateFarmRequest
	// 👉 recebe dados do body

	// 📥 valida JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	// 🧠 chama service para atualizar
	farm, err := h.Service.UpdateFarm(uint(id), req)
	if err != nil {
		// ❌ não encontrada
		c.JSON(http.StatusNotFound, gin.H{"error": "fazenda não encontrada"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, farm)
}

// 🗑️ Deletar fazenda
func (h *FarmHandler) DeleteFarm(c *gin.Context) {

	// 📥 pega ID da URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// 🧠 chama service para deletar
	err = h.Service.DeleteFarm(uint(id))
	if err != nil {
		// ❌ não encontrada
		c.JSON(http.StatusNotFound, gin.H{"error": "fazenda não encontrada"})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusOK, gin.H{"message": "fazenda removida com sucesso"})
}
