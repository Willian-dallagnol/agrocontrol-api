package handler

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
	"agrocontrol-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

// GetDashboard godoc
// @Summary      Dashboard principal
// @Description  Retorna visão consolidada com métricas de fazendas, talhões, plantios, insumos e alertas. Resultado cacheado por 5 minutos.
// @Tags         dashboard
// @Produce      json
// @Success      200 {object} dto.DashboardOverviewResponse
// @Security     BearerAuth
// @Router       /api/v1/dashboard [get]
func (h *ReportHandler) GetDashboard(c *gin.Context) {
	result, err := h.service.GetDashboardOverview(c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetProductivityReport godoc
// @Summary      Relatório de produtividade
// @Description  Retorna produtividade por talhão em sacas/ha e kg/ha, filtrado por safra, fazenda e cultura
// @Tags         reports
// @Produce      json
// @Param        season_id query int false "Filtrar por safra"
// @Param        farm_id   query int false "Filtrar por fazenda"
// @Param        crop_id   query int false "Filtrar por cultura"
// @Param        page      query int false "Página"
// @Param        limit     query int false "Itens por página (máx 200)"
// @Success      200 {object} dto.ProductivityReportResponse
// @Security     BearerAuth
// @Router       /api/v1/reports/productivity [get]
func (h *ReportHandler) GetProductivityReport(c *gin.Context) {
	var q dto.ProductivityReportQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetProductivityReport(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetCostPerFieldReport godoc
// @Summary      Relatório de custo por talhão
// @Description  Retorna o custo total de insumos aplicados por talhão, agrupado por categoria
// @Tags         reports
// @Produce      json
// @Param        field_id  query int false "Filtrar por talhão"
// @Param        farm_id   query int false "Filtrar por fazenda"
// @Param        season_id query int false "Filtrar por safra"
// @Success      200 {object} dto.CostPerFieldResponse
// @Security     BearerAuth
// @Router       /api/v1/reports/cost-per-field [get]
func (h *ReportHandler) GetCostPerFieldReport(c *gin.Context) {
	var q dto.CostPerFieldQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.RespondBadRequest(c, "parâmetros inválidos: "+err.Error())
		return
	}
	result, err := h.service.GetCostPerFieldReport(q, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		utils.RespondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
