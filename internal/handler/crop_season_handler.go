package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ── CropHandler ──────────────────────────────────────────────────────────────

type CropHandler struct{ service *service.CropService }

func NewCropHandler(s *service.CropService) *CropHandler { return &CropHandler{service: s} }

// CreateCrop godoc
// @Summary      Criar cultura
// @Description  Cadastra uma nova cultura agrícola (soja, milho, trigo, etc.)
// @Tags         crops
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateCropRequest true "Dados da cultura"
// @Success      201 {object} dto.CropResponse
// @Failure      400 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/crops [post]
func (h *CropHandler) CreateCrop(c *gin.Context) {
	var req dto.CreateCropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	crop, err := h.service.CreateCrop(req, c.GetUint("user_id"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, crop)
}

// GetCrops godoc
// @Summary      Listar culturas
// @Tags         crops
// @Produce      json
// @Param        page   query int    false "Página"
// @Param        limit  query int    false "Itens por página"
// @Param        search query string false "Busca por nome"
// @Success      200 {object} dto.PaginatedResponse
// @Security     BearerAuth
// @Router       /api/v1/crops [get]
func (h *CropHandler) GetCrops(c *gin.Context) {
	var q dto.CropQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetCropsPaginated(q)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetCropByID godoc
// @Summary      Buscar cultura por ID
// @Tags         crops
// @Produce      json
// @Param        id path int true "ID da cultura"
// @Success      200 {object} dto.CropResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/crops/{id} [get]
func (h *CropHandler) GetCropByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	crop, err := h.service.GetCropByID(id)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, crop)
}

// UpdateCrop godoc
// @Summary      Atualizar cultura
// @Tags         crops
// @Accept       json
// @Produce      json
// @Param        id      path int                  true "ID da cultura"
// @Param        request body dto.UpdateCropRequest true "Dados atualizados"
// @Success      200 {object} dto.CropResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/crops/{id} [put]
func (h *CropHandler) UpdateCrop(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdateCropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	crop, err := h.service.UpdateCrop(id, req, c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, crop)
}

// DeleteCrop godoc
// @Summary      Excluir cultura
// @Tags         crops
// @Param        id path int true "ID da cultura"
// @Success      204
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/crops/{id} [delete]
func (h *CropHandler) DeleteCrop(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	if err := h.service.DeleteCrop(id); err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── SeasonHandler ─────────────────────────────────────────────────────────────

type SeasonHandler struct{ service *service.SeasonService }

func NewSeasonHandler(s *service.SeasonService) *SeasonHandler { return &SeasonHandler{service: s} }

// CreateSeason godoc
// @Summary      Criar safra
// @Description  Cadastra uma nova safra com período de início e fim. Valida que data_fim > data_inicio.
// @Tags         seasons
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSeasonRequest true "Dados da safra"
// @Success      201 {object} dto.SeasonResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      422 {object} utils.ErrorResponse "Data fim anterior à data início"
// @Security     BearerAuth
// @Router       /api/v1/seasons [post]
func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	var req dto.CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	season, err := h.service.CreateSeason(req, c.GetUint("user_id"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, season)
}

// GetSeasons godoc
// @Summary      Listar safras
// @Tags         seasons
// @Produce      json
// @Param        page   query int    false "Página"
// @Param        limit  query int    false "Itens por página"
// @Param        search query string false "Busca por nome"
// @Success      200 {object} dto.PaginatedResponse
// @Security     BearerAuth
// @Router       /api/v1/seasons [get]
func (h *SeasonHandler) GetSeasons(c *gin.Context) {
	var q dto.SeasonQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetSeasonsPaginated(q)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetSeasonByID godoc
// @Summary      Buscar safra por ID
// @Tags         seasons
// @Produce      json
// @Param        id path int true "ID da safra"
// @Success      200 {object} dto.SeasonResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/seasons/{id} [get]
func (h *SeasonHandler) GetSeasonByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	season, err := h.service.GetSeasonByID(id)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, season)
}

// UpdateSeason godoc
// @Summary      Atualizar safra
// @Tags         seasons
// @Accept       json
// @Produce      json
// @Param        id      path int                    true "ID da safra"
// @Param        request body dto.UpdateSeasonRequest true "Dados atualizados"
// @Success      200 {object} dto.SeasonResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/seasons/{id} [put]
func (h *SeasonHandler) UpdateSeason(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	season, err := h.service.UpdateSeason(id, req)
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, season)
}

// DeleteSeason godoc
// @Summary      Excluir safra
// @Tags         seasons
// @Param        id path int true "ID da safra"
// @Success      204
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/seasons/{id} [delete]
func (h *SeasonHandler) DeleteSeason(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	if err := h.service.DeleteSeason(id); err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
