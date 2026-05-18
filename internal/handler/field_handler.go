package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FieldHandler struct{ service *service.FieldService }

func NewFieldHandler(s *service.FieldService) *FieldHandler { return &FieldHandler{service: s} }

// CreateField godoc
// @Summary      Criar talhão
// @Description  Cria um novo talhão em uma fazenda. Valida se o usuário tem acesso à fazenda.
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateFieldRequest true "Dados do talhão"
// @Success      201 {object} dto.FieldResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse "Fazenda não encontrada"
// @Security     BearerAuth
// @Router       /api/v1/fields [post]
func (h *FieldHandler) CreateField(c *gin.Context) {
	var req dto.CreateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	field, err := h.service.CreateField(req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, field)
}

// GetFields godoc
// @Summary      Listar talhões
// @Description  Retorna talhões paginados. Admin vê todos, outros veem apenas os seus.
// @Tags         fields
// @Produce      json
// @Param        page   query int    false "Página"
// @Param        limit  query int    false "Itens por página"
// @Param        search query string false "Busca por nome"
// @Success      200 {object} dto.PaginatedResponse
// @Security     BearerAuth
// @Router       /api/v1/fields [get]
func (h *FieldHandler) GetFields(c *gin.Context) {
	var q dto.FieldQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetFieldsPaginated(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetFieldByID godoc
// @Summary      Buscar talhão por ID
// @Tags         fields
// @Produce      json
// @Param        id path int true "ID do talhão"
// @Success      200 {object} dto.FieldResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/fields/{id} [get]
func (h *FieldHandler) GetFieldByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	field, err := h.service.GetFieldByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, field)
}

// GetFieldsByFarm godoc
// @Summary      Listar talhões de uma fazenda
// @Tags         fields
// @Produce      json
// @Param        id path int true "ID da fazenda"
// @Success      200 {array} dto.FieldResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/farms/{id}/fields [get]
func (h *FieldHandler) GetFieldsByFarm(c *gin.Context) {
	farmID, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id da fazenda inválido")
		return
	}
	fields, err := h.service.GetFieldsByFarmID(farmID, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, fields)
}

// UpdateField godoc
// @Summary      Atualizar talhão
// @Tags         fields
// @Accept       json
// @Produce      json
// @Param        id      path int                   true "ID do talhão"
// @Param        request body dto.UpdateFieldRequest true "Dados atualizados"
// @Success      200 {object} dto.FieldResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/fields/{id} [put]
func (h *FieldHandler) UpdateField(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	field, err := h.service.UpdateField(id, req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, field)
}

// DeleteField godoc
// @Summary      Excluir talhão
// @Tags         fields
// @Param        id path int true "ID do talhão"
// @Success      204
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/fields/{id} [delete]
func (h *FieldHandler) DeleteField(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	if err := h.service.DeleteField(id, c.GetUint("user_id"), c.GetString("role")); err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
