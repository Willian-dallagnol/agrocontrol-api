package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 🔐 Handler responsável pelas rotas de autenticação
type AuthHandler struct {
	Service *service.AuthService
	// 👉 referência ao service onde está a lógica de login (JWT, validação, etc)
}

// 🏗️ Função construtora do AuthHandler
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// 🔑 Endpoint de login
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	// 👉 estrutura que vai receber os dados enviados no body (email e senha)

	// 📥 Faz o bind do JSON da requisição para o struct
	if err := c.ShouldBindJSON(&req); err != nil {
		// ❌ erro de validação (ex: campos obrigatórios ou formato inválido)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	// 🧠 Chama o service para processar o login
	response, err := h.Service.Login(req)
	if err != nil {
		// ❌ erro de autenticação (email/senha inválidos)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ login realizado com sucesso
	c.JSON(http.StatusOK, response)
	// 👉 retorna token JWT + dados do usuário
}
