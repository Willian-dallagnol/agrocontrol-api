package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FarmHandler gerencia as rotas de fazendas
type FarmHandler struct {
	service *service.FarmService
}

func NewFarmHandler(s *service.FarmService) *FarmHandler {
	return &FarmHandler{service: s}
}

// CreateFarm godoc
// @Summary      Criar fazenda
// @Description  Cria uma nova fazenda. Requer role manager ou admin.
// @Tags         farms
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateFarmRequest true "Dados da fazenda"
// @Success      201 {object} dto.FarmResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Failure      422 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/farms [post]
func (h *FarmHandler) CreateFarm(c *gin.Context) {
	var req dto.CreateFarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	farm, err := h.service.CreateFarm(req, c.GetUint("user_id"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, farm)
}

// GetFarms godoc
// @Summary      Listar fazendas
// @Description  Retorna fazendas paginadas. Admin vê todas, outros veem apenas as suas.
// @Tags         farms
// @Produce      json
// @Param        page    query int    false "Página (padrão 1)"
// @Param        limit   query int    false "Itens por página (padrão 10)"
// @Param        search  query string false "Busca por nome"
// @Success      200 {object} dto.PaginatedResponse
// @Failure      401 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/farms [get]
func (h *FarmHandler) GetFarms(c *gin.Context) {
	var q dto.FarmQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetFarmsPaginated(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetFarmByID godoc
// @Summary      Buscar fazenda por ID
// @Description  Retorna os dados de uma fazenda específica
// @Tags         farms
// @Produce      json
// @Param        id path int true "ID da fazenda"
// @Success      200 {object} dto.FarmResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/farms/{id} [get]
func (h *FarmHandler) GetFarmByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	farm, err := h.service.GetFarmByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, farm)
}

// UpdateFarm godoc
// @Summary      Atualizar fazenda
// @Description  Atualiza os dados de uma fazenda. Requer role manager ou admin.
// @Tags         farms
// @Accept       json
// @Produce      json
// @Param        id      path int                  true "ID da fazenda"
// @Param        request body dto.UpdateFarmRequest true "Dados atualizados"
// @Success      200 {object} dto.FarmResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/farms/{id} [put]
func (h *FarmHandler) UpdateFarm(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdateFarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	farm, err := h.service.UpdateFarm(id, req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, farm)
}

// DeleteFarm godoc
// @Summary      Excluir fazenda
// @Description  Exclui uma fazenda. Requer role admin.
// @Tags         farms
// @Produce      json
// @Param        id path int true "ID da fazenda"
// @Success      204
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/farms/{id} [delete]
func (h *FarmHandler) DeleteFarm(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	if err := h.service.DeleteFarm(id, c.GetUint("user_id"), c.GetString("role")); err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// parseID extrai e valida o parâmetro :id da URL
func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}
