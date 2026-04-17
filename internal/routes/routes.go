package routes

import (
	"agrocontrol-api/internal/handler"
	"agrocontrol-api/internal/middleware"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 🌐 Função responsável por registrar todas as rotas da aplicação
func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {

	// ✅ Rota simples para verificar se a API está online
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "AgroControl API rodando 🚜",
		})
	})

	// ===================== USER & AUTH =====================

	// 👤 USER
	// cria repository → service → handler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 🔐 AUTH
	// cria service e handler de autenticação
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// rotas públicas
	router.POST("/users", userHandler.CreateUser) // cadastro
	router.POST("/login", authHandler.Login)      // login

	// rotas protegidas de autenticação
	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.GET("/me", func(c *gin.Context) {
			// retorna dados do usuário autenticado a partir do token
			c.JSON(200, gin.H{
				"user_id": c.GetUint("user_id"),
				"email":   c.GetString("email"),
				"role":    c.GetString("role"),
			})
		})
	}

	// ===================== FARM =====================

	// 🚜 FARM
	farmRepo := repository.NewFarmRepository(db)
	farmService := service.NewFarmService(farmRepo)
	farmHandler := handler.NewFarmHandler(farmService)

	// ===================== FIELD =====================

	// 🌱 FIELD
	fieldRepo := repository.NewFieldRepository(db)
	fieldService := service.NewFieldService(fieldRepo, farmRepo)
	fieldHandler := handler.NewFieldHandler(fieldService)

	// ===================== CROP =====================

	// 🌾 CROP
	cropRepo := repository.NewCropRepository(db)
	cropService := service.NewCropService(cropRepo, fieldRepo)
	cropHandler := handler.NewCropHandler(cropService)

	// ===================== FARM ROUTES =====================

	farms := router.Group("/farms")
	farms.Use(middleware.AuthMiddleware(jwtSecret)) // protege todas as rotas do grupo
	{
		farms.POST("", farmHandler.CreateFarm)       // criar fazenda
		farms.GET("", farmHandler.GetFarms)          // listar fazendas
		farms.GET("/:id", farmHandler.GetFarmByID)   // buscar por ID
		farms.PUT("/:id", farmHandler.UpdateFarm)    // atualizar
		farms.DELETE("/:id", farmHandler.DeleteFarm) // deletar

		// 🔗 rota avançada: listar talhões de uma fazenda
		farms.GET("/:id/fields", fieldHandler.GetFieldsByFarm)
	}

	// ===================== FIELD ROUTES =====================

	fields := router.Group("/fields")
	fields.Use(middleware.AuthMiddleware(jwtSecret)) // protege todas as rotas do grupo
	{
		fields.POST("", fieldHandler.CreateField)       // criar talhão
		fields.GET("", fieldHandler.GetFields)          // listar todos
		fields.GET("/:id", fieldHandler.GetFieldByID)   // buscar por ID
		fields.PUT("/:id", fieldHandler.UpdateField)    // atualizar
		fields.DELETE("/:id", fieldHandler.DeleteField) // deletar
	}

	// ===================== CROP ROUTES =====================

	crops := router.Group("/crops")
	crops.Use(middleware.AuthMiddleware(jwtSecret)) // protege todas as rotas do grupo
	{
		crops.POST("", cropHandler.CreateCrop)       // criar cultura
		crops.GET("", cropHandler.GetCrops)          // listar todas
		crops.GET("/:id", cropHandler.GetCropByID)   // buscar por ID
		crops.PUT("/:id", cropHandler.UpdateCrop)    // atualizar
		crops.DELETE("/:id", cropHandler.DeleteCrop) // deletar
	}
}
