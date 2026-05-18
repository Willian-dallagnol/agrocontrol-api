package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ── PlantingHandler ───────────────────────────────────────────────────────────

type PlantingHandler struct{ service *service.PlantingService }

func NewPlantingHandler(s *service.PlantingService) *PlantingHandler {
	return &PlantingHandler{service: s}
}

// CreatePlanting godoc
// @Summary      Criar plantio
// @Description  Registra um novo plantio em um talhão ativo. Valida se o talhão está ativo e se o usuário tem acesso.
// @Tags         plantings
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePlantingRequest true "Dados do plantio"
// @Success      201 {object} dto.PlantingResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      422 {object} utils.ErrorResponse "Talhão inativo"
// @Security     BearerAuth
// @Router       /api/v1/plantings [post]
func (h *PlantingHandler) CreatePlanting(c *gin.Context) {
	var req dto.CreatePlantingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	p, err := h.service.CreatePlanting(req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

// GetPlantings godoc
// @Summary      Listar plantios
// @Description  Admin vê todos os plantios, outros veem apenas os seus
// @Tags         plantings
// @Produce      json
// @Success      200 {array} dto.PlantingResponse
// @Failure      401 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/plantings [get]
func (h *PlantingHandler) GetPlantings(c *gin.Context) {
	result, err := h.service.GetPlantings(c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetPlantingByID godoc
// @Summary      Buscar plantio por ID
// @Tags         plantings
// @Produce      json
// @Param        id path int true "ID do plantio"
// @Success      200 {object} dto.PlantingResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/plantings/{id} [get]
func (h *PlantingHandler) GetPlantingByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	p, err := h.service.GetPlantingByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

// UpdatePlanting godoc
// @Summary      Atualizar plantio
// @Tags         plantings
// @Accept       json
// @Produce      json
// @Param        id      path int                     true "ID do plantio"
// @Param        request body dto.UpdatePlantingRequest true "Dados atualizados"
// @Success      200 {object} dto.PlantingResponse
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/plantings/{id} [put]
func (h *PlantingHandler) UpdatePlanting(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdatePlantingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	p, err := h.service.UpdatePlanting(id, req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

// DeletePlanting godoc
// @Summary      Excluir plantio
// @Tags         plantings
// @Param        id path int true "ID do plantio"
// @Success      204
// @Failure      403 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/plantings/{id} [delete]
func (h *PlantingHandler) DeletePlanting(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	if err := h.service.DeletePlanting(id, c.GetUint("user_id"), c.GetString("role")); err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── InputHandler ──────────────────────────────────────────────────────────────

type InputHandler struct{ service *service.InputService }

func NewInputHandler(s *service.InputService) *InputHandler { return &InputHandler{service: s} }

// CreateInput godoc
// @Summary      Criar insumo
// @Description  Cadastra um novo insumo agrícola com controle de estoque
// @Tags         inputs
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateInputRequest true "Dados do insumo"
// @Success      201 {object} dto.InputResponse
// @Failure      400 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/inputs [post]
func (h *InputHandler) CreateInput(c *gin.Context) {
	var req dto.CreateInputRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	input, err := h.service.CreateInput(req, c.GetUint("user_id"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, input)
}

// GetInputs godoc
// @Summary      Listar insumos
// @Description  Retorna insumos paginados com filtro por categoria
// @Tags         inputs
// @Produce      json
// @Param        page     query int    false "Página"
// @Param        limit    query int    false "Itens por página"
// @Param        search   query string false "Busca por nome"
// @Param        category query string false "Filtro por categoria (fertilizer, herbicide, fungicide...)"
// @Success      200 {object} dto.PaginatedResponse
// @Security     BearerAuth
// @Router       /api/v1/inputs [get]
func (h *InputHandler) GetInputs(c *gin.Context) {
	var q dto.InputQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetInputsPaginated(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetInputByID godoc
// @Summary      Buscar insumo por ID
// @Tags         inputs
// @Produce      json
// @Param        id path int true "ID do insumo"
// @Success      200 {object} dto.InputResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/inputs/{id} [get]
func (h *InputHandler) GetInputByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	input, err := h.service.GetInputByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

// UpdateInput godoc
// @Summary      Atualizar insumo
// @Tags         inputs
// @Accept       json
// @Produce      json
// @Param        id      path int                  true "ID do insumo"
// @Param        request body dto.UpdateInputRequest true "Dados atualizados"
// @Success      200 {object} dto.InputResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/inputs/{id} [put]
func (h *InputHandler) UpdateInput(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdateInputRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	input, err := h.service.UpdateInput(id, req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

// DeleteInput godoc
// @Summary      Excluir insumo
// @Tags         inputs
// @Param        id path int true "ID do insumo"
// @Success      204
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/inputs/{id} [delete]
func (h *InputHandler) DeleteInput(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	if err := h.service.DeleteInput(id, c.GetUint("user_id"), c.GetString("role")); err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AdjustStock godoc
// @Summary      Ajustar estoque
// @Description  Ajusta o estoque do insumo. Positivo = entrada, negativo = saída manual.
// @Tags         inputs
// @Accept       json
// @Produce      json
// @Param        id      path int                   true "ID do insumo"
// @Param        request body dto.AdjustStockRequest true "Quantidade a ajustar"
// @Success      200 {object} dto.InputResponse
// @Failure      422 {object} utils.ErrorResponse "Estoque insuficiente"
// @Security     BearerAuth
// @Router       /api/v1/inputs/{id}/stock [patch]
func (h *InputHandler) AdjustStock(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	input, err := h.service.AdjustStock(id, req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

// ── ApplicationHandler ────────────────────────────────────────────────────────

type ApplicationHandler struct{ service *service.ApplicationService }

func NewApplicationHandler(s *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: s}
}

// CreateApplication godoc
// @Summary      Registrar aplicação de insumo
// @Description  Registra uma aplicação de insumo em um talhão. Debita o estoque automaticamente e cria alerta se estoque ficar baixo.
// @Tags         applications
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateApplicationRequest true "Dados da aplicação"
// @Success      201 {object} dto.ApplicationResponse
// @Failure      422 {object} utils.ErrorResponse "Estoque insuficiente ou talhão inativo"
// @Security     BearerAuth
// @Router       /api/v1/applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	app, err := h.service.CreateApplication(req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, app)
}

// GetApplications godoc
// @Summary      Listar aplicações
// @Tags         applications
// @Produce      json
// @Param        page     query int  false "Página"
// @Param        limit    query int  false "Itens por página"
// @Param        field_id query int  false "Filtrar por talhão"
// @Success      200 {object} dto.PaginatedResponse
// @Security     BearerAuth
// @Router       /api/v1/applications [get]
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	var q dto.ApplicationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetApplicationsPaginated(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetApplicationByID godoc
// @Summary      Buscar aplicação por ID
// @Tags         applications
// @Produce      json
// @Param        id path int true "ID da aplicação"
// @Success      200 {object} dto.ApplicationResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/applications/{id} [get]
func (h *ApplicationHandler) GetApplicationByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	app, err := h.service.GetApplicationByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, app)
}

// GetApplicationsByField godoc
// @Summary      Listar aplicações por talhão
// @Tags         applications
// @Produce      json
// @Param        id path int true "ID do talhão"
// @Success      200 {array} dto.ApplicationResponse
// @Security     BearerAuth
// @Router       /api/v1/fields/{id}/applications [get]
func (h *ApplicationHandler) GetApplicationsByField(c *gin.Context) {
	fieldID, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id do talhão inválido")
		return
	}
	apps, err := h.service.GetApplicationsByField(fieldID, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, apps)
}
