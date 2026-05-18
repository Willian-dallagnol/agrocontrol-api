package routes

import (
	"agrocontrol-api/internal/handler"
	"agrocontrol-api/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth        *handler.AuthHandler
	User        *handler.UserHandler
	Farm        *handler.FarmHandler
	Field       *handler.FieldHandler
	Crop        *handler.CropHandler
	Season      *handler.SeasonHandler
	Planting    *handler.PlantingHandler
	Input       *handler.InputHandler
	Application *handler.ApplicationHandler
	Monitoring  *handler.MonitoringHandler
	Harvest     *handler.HarvestHandler
	Alert       *handler.AlertHandler
	Report      *handler.ReportHandler
}

func Setup(r *gin.Engine, h Handlers, jwtSecret string) {
	// ── Healthcheck público ───────────────────────────────────────────────
	r.GET("/health", healthHandler)

	// ── Autenticação pública ──────────────────────────────────────────────
	r.POST("/auth/login", h.Auth.Login)
	r.POST("/auth/refresh", h.Auth.RefreshToken)

	// ── Rotas protegidas por JWT ──────────────────────────────────────────
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtSecret))

	api.GET("/me", h.User.Me)
	api.GET("/dashboard", h.Report.GetDashboard)

	// Usuários (admin only)
	users := api.Group("/users", middleware.AdminOnly())
	{
		users.POST("", h.User.CreateUser)
		users.GET("", h.User.GetUsers)
		users.GET("/:id", h.User.GetUserByID)
	}

	// Fazendas
	farms := api.Group("/farms")
	{
		farms.POST("", middleware.ManagerOrAbove(), h.Farm.CreateFarm)
		farms.GET("", h.Farm.GetFarms)
		farms.GET("/:id", h.Farm.GetFarmByID)
		farms.PUT("/:id", middleware.ManagerOrAbove(), h.Farm.UpdateFarm)
		farms.DELETE("/:id", middleware.AdminOnly(), h.Farm.DeleteFarm)
		farms.GET("/:id/fields", h.Field.GetFieldsByFarm)
	}

	// Talhões
	fields := api.Group("/fields")
	{
		fields.POST("", middleware.ManagerOrAbove(), h.Field.CreateField)
		fields.GET("", h.Field.GetFields)
		fields.GET("/:id", h.Field.GetFieldByID)
		fields.PUT("/:id", middleware.ManagerOrAbove(), h.Field.UpdateField)
		fields.DELETE("/:id", middleware.AdminOnly(), h.Field.DeleteField)
		fields.GET("/:id/applications", h.Application.GetApplicationsByField)
		fields.GET("/:id/monitorings", h.Monitoring.GetMonitoringsByField)
	}

	// Culturas
	crops := api.Group("/crops")
	{
		crops.POST("", middleware.ManagerOrAbove(), h.Crop.CreateCrop)
		crops.GET("", h.Crop.GetCrops)
		crops.GET("/:id", h.Crop.GetCropByID)
		crops.PUT("/:id", middleware.ManagerOrAbove(), h.Crop.UpdateCrop)
		crops.DELETE("/:id", middleware.AdminOnly(), h.Crop.DeleteCrop)
	}

	// Safras
	seasons := api.Group("/seasons")
	{
		seasons.POST("", middleware.ManagerOrAbove(), h.Season.CreateSeason)
		seasons.GET("", h.Season.GetSeasons)
		seasons.GET("/:id", h.Season.GetSeasonByID)
		seasons.PUT("/:id", middleware.ManagerOrAbove(), h.Season.UpdateSeason)
		seasons.DELETE("/:id", middleware.AdminOnly(), h.Season.DeleteSeason)
	}

	// Plantios
	plantings := api.Group("/plantings")
	{
		plantings.POST("", middleware.ManagerOrAbove(), h.Planting.CreatePlanting)
		plantings.GET("", h.Planting.GetPlantings)
		plantings.GET("/:id", h.Planting.GetPlantingByID)
		plantings.PUT("/:id", middleware.ManagerOrAbove(), h.Planting.UpdatePlanting)
		plantings.DELETE("/:id", middleware.AdminOnly(), h.Planting.DeletePlanting)
	}

	// Insumos
	inputs := api.Group("/inputs")
	{
		inputs.POST("", middleware.ManagerOrAbove(), h.Input.CreateInput)
		inputs.GET("", h.Input.GetInputs)
		inputs.GET("/:id", h.Input.GetInputByID)
		inputs.PUT("/:id", middleware.ManagerOrAbove(), h.Input.UpdateInput)
		inputs.DELETE("/:id", middleware.AdminOnly(), h.Input.DeleteInput)
		inputs.POST("/:id/adjust-stock", middleware.ManagerOrAbove(), h.Input.AdjustStock)
	}

	// Aplicações
	applications := api.Group("/applications")
	{
		applications.POST("", h.Application.CreateApplication)
		applications.GET("", h.Application.GetApplications)
		applications.GET("/:id", h.Application.GetApplicationByID)
	}

	// Monitoramentos
	monitorings := api.Group("/monitorings")
	{
		monitorings.POST("", h.Monitoring.CreateMonitoring)
		monitorings.GET("", h.Monitoring.GetMonitorings)
		monitorings.GET("/:id", h.Monitoring.GetMonitoringByID)
	}

	// Colheitas
	harvests := api.Group("/harvests")
	{
		harvests.POST("", middleware.ManagerOrAbove(), h.Harvest.CreateHarvest)
		harvests.GET("", h.Harvest.GetHarvests)
		harvests.GET("/:id", h.Harvest.GetHarvestByID)
	}

	// Alertas
	alerts := api.Group("/alerts")
	{
		alerts.POST("", h.Alert.CreateAlert)
		alerts.GET("", h.Alert.GetAlerts)
		alerts.GET("/open", h.Alert.GetOpenAlerts)
		alerts.GET("/:id", h.Alert.GetAlertByID)
		alerts.PATCH("/:id/status", h.Alert.UpdateStatus)
	}

	// Relatórios
	reports := api.Group("/reports")
	{
		reports.GET("/productivity", h.Report.GetProductivityReport)
		reports.GET("/cost-per-field", h.Report.GetCostPerFieldReport)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "2.0.0",
	})
}
