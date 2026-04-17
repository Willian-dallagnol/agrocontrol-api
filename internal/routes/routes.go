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

	farmRepo := repository.NewFarmRepository(db)
	farmService := service.NewFarmService(farmRepo)
	farmHandler := handler.NewFarmHandler(farmService)

	farms := router.Group("/farms")
	farms.Use(middleware.AuthMiddleware(jwtSecret))
	{
		farms.POST("", farmHandler.CreateFarm)
		farms.GET("", farmHandler.GetFarms)
		farms.GET("/:id", farmHandler.GetFarmByID)
		farms.PUT("/:id", farmHandler.UpdateFarm)
		farms.DELETE("/:id", farmHandler.DeleteFarm)
	}
}
