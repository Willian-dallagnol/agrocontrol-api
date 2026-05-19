// @title           AgroControl API
// @version         2.0
// @description     API REST de gestão agrícola com autenticação JWT
// @host            agrocontrol-api-production.up.railway.app
// @schemes         https
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agrocontrol-api/configs"
	_ "agrocontrol-api/docs"
	"agrocontrol-api/internal/cache"
	"agrocontrol-api/internal/database"
	"agrocontrol-api/internal/handler"
	"agrocontrol-api/internal/middleware"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/routes"
	"agrocontrol-api/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := configs.LoadConfig()

	logLevel := slog.LevelInfo
	if cfg.Env == "development" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	db := configs.ConnectDatabase(cfg)
	redisClient := cache.NewClient(cfg.RedisAddr)

	if err := database.RunMigrations(db, "./migrations"); err != nil {
		log.Fatal("[MIGRATE] Falha: ", err)
	}

	// ── Repositórios concretos ────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	farmRepo := repository.NewFarmRepository(db)
	fieldRepo := repository.NewFieldRepository(db)
	cropRepo := repository.NewCropRepository(db)
	seasonRepo := repository.NewSeasonRepository(db)
	plantingRepo := repository.NewPlantingRepository(db)
	inputRepo := repository.NewInputRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	monitoringRepo := repository.NewMonitoringRepository(db)
	harvestRepo := repository.NewHarvestRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// ── TxRunner — único ponto de acoplamento ao GORM nos serviços ───────
	txRunner := repository.NewGormTxRunner(db)

	// ── Serviços com interfaces (ports) ───────────────────────────────────
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpHours)
	userService := service.NewUserService(userRepo)
	farmService := service.NewFarmService(farmRepo)
	fieldService := service.NewFieldService(fieldRepo, farmRepo)
	cropService := service.NewCropService(cropRepo)
	seasonService := service.NewSeasonService(seasonRepo)
	plantingService := service.NewPlantingService(plantingRepo, fieldRepo, seasonRepo, cropRepo)
	inputService := service.NewInputService(inputRepo, alertRepo, txRunner)
	applicationService := service.NewApplicationService(applicationRepo, fieldRepo, inputRepo, alertRepo, txRunner)
	monitoringService := service.NewMonitoringService(monitoringRepo, fieldRepo, alertRepo)
	harvestService := service.NewHarvestService(harvestRepo, plantingRepo, fieldRepo, txRunner)
	alertService := service.NewAlertService(alertRepo)
	reportService := service.NewReportService(reportRepo, redisClient)

	h := routes.Handlers{
		Auth:        handler.NewAuthHandler(authService),
		User:        handler.NewUserHandler(userService),
		Farm:        handler.NewFarmHandler(farmService),
		Field:       handler.NewFieldHandler(fieldService),
		Crop:        handler.NewCropHandler(cropService),
		Season:      handler.NewSeasonHandler(seasonService),
		Planting:    handler.NewPlantingHandler(plantingService),
		Input:       handler.NewInputHandler(inputService),
		Application: handler.NewApplicationHandler(applicationService),
		Monitoring:  handler.NewMonitoringHandler(monitoringService),
		Harvest:     handler.NewHarvestHandler(harvestService),
		Alert:       handler.NewAlertHandler(alertService),
		Report:      handler.NewReportHandler(reportService),
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID())       // UUID por request — rastreabilidade
	r.Use(middleware.Logger())          // log estruturado com request_id
	r.Use(middleware.SecurityHeaders()) // security headers em todas as respostas
	r.Use(gin.Recovery())
	r.Use(middleware.NewRateLimiter(10, 30).Middleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	routes.Setup(r, h, cfg.JWTSecret)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server: iniciando", "addr", srv.Addr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("[SERVER] Falha ao iniciar: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("server: recebido sinal de shutdown — aguardando requests em andamento...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server: erro no shutdown", "error", err)
	} else {
		slog.Info("server: encerrado com sucesso")
	}
}
