package routes

import (
	"agrocontrol-api/internal/handler"
	"agrocontrol-api/internal/middleware"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "AgroControl API rodando 🚜",
		})
	})

	// ===================== USER & AUTH =====================

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	router.POST("/users", userHandler.CreateUser)
	router.POST("/login", authHandler.Login)

	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetUint("user_id"),
				"email":   c.GetString("email"),
				"role":    c.GetString("role"),
			})
		})
	}

	// ===================== FARM =====================

	farmRepo := repository.NewFarmRepository(db)
	farmService := service.NewFarmService(farmRepo)
	farmHandler := handler.NewFarmHandler(farmService)

	// ===================== FIELD =====================

	fieldRepo := repository.NewFieldRepository(db)
	fieldService := service.NewFieldService(fieldRepo, farmRepo)
	fieldHandler := handler.NewFieldHandler(fieldService)

	// ===================== CROP =====================

	cropRepo := repository.NewCropRepository(db)
	cropService := service.NewCropService(cropRepo, fieldRepo)
	cropHandler := handler.NewCropHandler(cropService)

	// ===================== FARM ROUTES =====================

	farms := router.Group("/farms")
	farms.Use(middleware.AuthMiddleware(jwtSecret))
	{
		farms.POST("", farmHandler.CreateFarm)
		farms.GET("", farmHandler.GetFarms)
		farms.GET("/:id", farmHandler.GetFarmByID)
		farms.PUT("/:id", farmHandler.UpdateFarm)
		farms.DELETE("/:id", farmHandler.DeleteFarm)
		farms.GET("/:id/fields", fieldHandler.GetFieldsByFarm)
	}

	// ===================== FIELD ROUTES =====================

	fields := router.Group("/fields")
	fields.Use(middleware.AuthMiddleware(jwtSecret))
	{
		fields.POST("", fieldHandler.CreateField)
		fields.GET("", fieldHandler.GetFields)
		fields.GET("/:id", fieldHandler.GetFieldByID)
		fields.PUT("/:id", fieldHandler.UpdateField)
		fields.DELETE("/:id", fieldHandler.DeleteField)
	}

	// ===================== CROP ROUTES =====================

	crops := router.Group("/crops")
	crops.Use(middleware.AuthMiddleware(jwtSecret))
	{
		crops.POST("", cropHandler.CreateCrop)
		crops.GET("", cropHandler.GetCrops)
		crops.GET("/:id", cropHandler.GetCropByID)
		crops.PUT("/:id", cropHandler.UpdateCrop)
		crops.DELETE("/:id", cropHandler.DeleteCrop)
	}
}
