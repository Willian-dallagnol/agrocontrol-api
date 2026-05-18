package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ── MonitoringHandler ────────────────────────────────────────────────────────

type MonitoringHandler struct{ service *service.MonitoringService }

func NewMonitoringHandler(s *service.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{service: s}
}

// CreateMonitoring godoc
// @Summary      Registrar monitoramento
// @Description  Registra um monitoramento de pragas, doenças ou condições do talhão. Cria alerta automático se urgente ou severidade crítica.
// @Tags         monitorings
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateMonitoringRequest true "Dados do monitoramento"
// @Success      201 {object} dto.MonitoringResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      403 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/monitorings [post]
func (h *MonitoringHandler) CreateMonitoring(c *gin.Context) {
	var req dto.CreateMonitoringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	m, err := h.service.CreateMonitoring(req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, m)
}

// GetMonitorings godoc
// @Summary      Listar monitoramentos
// @Description  Admin vê todos, outros veem apenas os seus
// @Tags         monitorings
// @Produce      json
// @Success      200 {array} dto.MonitoringResponse
// @Security     BearerAuth
// @Router       /api/v1/monitorings [get]
func (h *MonitoringHandler) GetMonitorings(c *gin.Context) {
	mons, err := h.service.GetMonitorings(c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, mons)
}

// GetMonitoringByID godoc
// @Summary      Buscar monitoramento por ID
// @Tags         monitorings
// @Produce      json
// @Param        id path int true "ID do monitoramento"
// @Success      200 {object} dto.MonitoringResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/monitorings/{id} [get]
func (h *MonitoringHandler) GetMonitoringByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	m, err := h.service.GetMonitoringByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

// GetMonitoringsByField godoc
// @Summary      Listar monitoramentos por talhão
// @Tags         monitorings
// @Produce      json
// @Param        id path int true "ID do talhão"
// @Success      200 {array} dto.MonitoringResponse
// @Security     BearerAuth
// @Router       /api/v1/fields/{id}/monitorings [get]
func (h *MonitoringHandler) GetMonitoringsByField(c *gin.Context) {
	fieldID, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id do talhão inválido")
		return
	}
	mons, err := h.service.GetMonitoringsByField(fieldID, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, mons)
}

// ── HarvestHandler ───────────────────────────────────────────────────────────

type HarvestHandler struct{ service *service.HarvestService }

func NewHarvestHandler(s *service.HarvestService) *HarvestHandler { return &HarvestHandler{service: s} }

// CreateHarvest godoc
// @Summary      Registrar colheita
// @Description  Registra a colheita de um plantio ativo. Atualiza o status do plantio para 'harvested'. Valida que o plantio não foi colhido antes.
// @Tags         harvests
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateHarvestRequest true "Dados da colheita"
// @Success      201 {object} dto.HarvestResponse
// @Failure      404 {object} utils.ErrorResponse "Plantio não encontrado"
// @Failure      422 {object} utils.ErrorResponse "Plantio já colhido ou inativo"
// @Security     BearerAuth
// @Router       /api/v1/harvests [post]
func (h *HarvestHandler) CreateHarvest(c *gin.Context) {
	var req dto.CreateHarvestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	harvest, err := h.service.CreateHarvest(req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, harvest)
}

// GetHarvests godoc
// @Summary      Listar colheitas
// @Tags         harvests
// @Produce      json
// @Param        page     query int false "Página"
// @Param        limit    query int false "Itens por página"
// @Param        field_id query int false "Filtrar por talhão"
// @Success      200 {object} dto.PaginatedResponse
// @Security     BearerAuth
// @Router       /api/v1/harvests [get]
func (h *HarvestHandler) GetHarvests(c *gin.Context) {
	var q dto.HarvestQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetHarvestsPaginated(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetHarvestByID godoc
// @Summary      Buscar colheita por ID
// @Tags         harvests
// @Produce      json
// @Param        id path int true "ID da colheita"
// @Success      200 {object} dto.HarvestResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/harvests/{id} [get]
func (h *HarvestHandler) GetHarvestByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	harvest, err := h.service.GetHarvestByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, harvest)
}

// ── AlertHandler ─────────────────────────────────────────────────────────────

type AlertHandler struct{ service *service.AlertService }

func NewAlertHandler(s *service.AlertService) *AlertHandler { return &AlertHandler{service: s} }

// CreateAlert godoc
// @Summary      Criar alerta
// @Description  Cria um alerta manual. Alertas automáticos são criados pelo sistema ao detectar estoque baixo ou monitoramento urgente.
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateAlertRequest true "Dados do alerta"
// @Success      201 {object} dto.AlertResponse
// @Security     BearerAuth
// @Router       /api/v1/alerts [post]
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var req dto.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	alert, err := h.service.CreateAlert(req, c.GetUint("user_id"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, alert)
}

// GetAlerts godoc
// @Summary      Listar todos os alertas
// @Tags         alerts
// @Produce      json
// @Success      200 {array} dto.AlertResponse
// @Security     BearerAuth
// @Router       /api/v1/alerts [get]
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	alerts, err := h.service.GetAlerts(c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, alerts)
}

// GetOpenAlerts godoc
// @Summary      Listar alertas abertos
// @Description  Retorna apenas alertas com status 'open', ordenados por prioridade. Limite de 20 alertas.
// @Tags         alerts
// @Produce      json
// @Success      200 {array} dto.AlertResponse
// @Security     BearerAuth
// @Router       /api/v1/alerts/open [get]
func (h *AlertHandler) GetOpenAlerts(c *gin.Context) {
	alerts, err := h.service.GetOpenAlerts(c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, alerts)
}

// GetAlertByID godoc
// @Summary      Buscar alerta por ID
// @Tags         alerts
// @Produce      json
// @Param        id path int true "ID do alerta"
// @Success      200 {object} dto.AlertResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/alerts/{id} [get]
func (h *AlertHandler) GetAlertByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	alert, err := h.service.GetAlertByID(id, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, alert)
}

// UpdateStatus godoc
// @Summary      Atualizar status do alerta
// @Description  Atualiza o status do alerta (open, resolved, ignored). Ao resolver, registra o timestamp de resolução.
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        id      path int                        true "ID do alerta"
// @Param        request body dto.UpdateAlertStatusRequest true "Novo status"
// @Success      200 {object} dto.AlertResponse
// @Failure      404 {object} utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/alerts/{id}/status [patch]
func (h *AlertHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.RespondBadRequest(c, "id inválido")
		return
	}
	var req dto.UpdateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "dados inválidos: "+err.Error())
		return
	}
	alert, err := h.service.UpdateStatus(id, req, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, alert)
}

// ── DashboardHandler ──────────────────────────────────────────────────────────

type DashboardHandler struct {
	farmRepo        interface{ Count() (int64, error) }
	fieldRepo       interface{ Count() (int64, error) }
	plantingRepo    interface{ CountActive() (int64, error) }
	alertRepo       interface{ CountOpen() (int64, error) }
	inputRepo       interface{ CountLowStock() (int64, error) }
	applicationRepo interface{ CountThisMonth() (int64, error) }
	alertService    *service.AlertService
}

func NewDashboardHandler(
	farmRepo interface{ Count() (int64, error) },
	fieldRepo interface{ Count() (int64, error) },
	plantingRepo interface{ CountActive() (int64, error) },
	alertRepo interface{ CountOpen() (int64, error) },
	inputRepo interface{ CountLowStock() (int64, error) },
	appRepo interface{ CountThisMonth() (int64, error) },
	alertService *service.AlertService,
) *DashboardHandler {
	return &DashboardHandler{
		farmRepo: farmRepo, fieldRepo: fieldRepo, plantingRepo: plantingRepo,
		alertRepo: alertRepo, inputRepo: inputRepo, applicationRepo: appRepo,
		alertService: alertService,
	}
}

// GetDashboard godoc
// @Summary      Dashboard consolidado
// @Description  Retorna métricas consolidadas: total de fazendas, talhões, plantios ativos, alertas abertos, insumos com estoque baixo e aplicações do mês.
// @Tags         dashboard
// @Produce      json
// @Success      200 {object} dto.DashboardResponse
// @Security     BearerAuth
// @Router       /api/v1/dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, role := c.GetUint("user_id"), c.GetString("role")

	totalFarms, _ := h.farmRepo.Count()
	totalFields, _ := h.fieldRepo.Count()
	activePlantings, _ := h.plantingRepo.CountActive()
	openAlerts, _ := h.alertRepo.CountOpen()
	lowStockInputs, _ := h.inputRepo.CountLowStock()
	appsThisMonth, _ := h.applicationRepo.CountThisMonth()
	lastAlerts, _ := h.alertService.GetOpenAlerts(userID, role)

	c.JSON(http.StatusOK, dto.DashboardResponse{
		TotalFarms:            totalFarms,
		TotalFields:           totalFields,
		ActivePlantings:       activePlantings,
		OpenAlerts:            openAlerts,
		LowStockInputs:        lowStockInputs,
		ApplicationsThisMonth: appsThisMonth,
		LastAlerts:            lastAlerts,
	})
}
