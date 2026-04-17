package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 👤 Handler responsável pelas operações de usuário
type UserHandler struct {
	Service *service.UserService
	// 👉 referência ao service onde está a lógica (ex: hash de senha)
}

// 🏗️ Construtor do handler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// 🚀 Criar novo usuário
func (h *UserHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest
	// 👉 estrutura que recebe os dados do body (nome, email, senha, role)

	// 📥 valida e faz bind do JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		// ❌ erro de validação (campos obrigatórios ou formato inválido)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	// 🧠 chama o service para criar o usuário
	user, err := h.Service.CreateUser(req)
	if err != nil {
		// ❌ erro de regra de negócio (ex: email já existe)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ sucesso
	c.JSON(http.StatusCreated, user)
	// 👉 retorna usuário sem expor senha
}
