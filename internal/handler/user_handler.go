package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler gerencia rotas de usuários (admin only)
type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser godoc
// @Summary      Criar usuário
// @Description  Cria um novo usuário. Apenas admin pode criar usuários. Roles disponíveis: admin, manager, operator.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserRequest true "Dados do usuário"
// @Success      201 {object} dto.UserResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      409 {object} utils.ErrorResponse "Email já cadastrado"
// @Security     BearerAuth
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	user, err := h.service.CreateUser(req)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary      Listar usuários
// @Description  Retorna todos os usuários do sistema. Apenas admin.
// @Tags         users
// @Produce      json
// @Success      200 {array} dto.UserResponse
// @Security     BearerAuth
// @Router       /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		utils.RespondInternalError(c, "erro ao buscar usuários")
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary      Buscar usuário por ID
// @Description  Retorna os dados de um usuário específico. Apenas admin.
// @Tags         users
// @Produce      json
// @Param        id path int true "ID do usuário"
// @Success      200 {object} dto.UserResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Me godoc
// @Summary      Meu perfil
// @Description  Retorna os dados do usuário autenticado
// @Tags         users
// @Produce      json
// @Success      200 {object} dto.UserResponse
// @Security     BearerAuth
// @Router       /api/v1/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	user, err := h.service.GetUserByID(c.GetUint("user_id"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
