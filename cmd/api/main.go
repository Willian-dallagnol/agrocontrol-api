package main

import (
	"agrocontrol-api/configs"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.LoadConfig()

	db := configs.ConnectDatabase(cfg)

	err := db.AutoMigrate(&entities.User{}, &entities.Farm{})
	if err != nil {
		log.Fatal("Erro ao executar migration:", err)
	}

	router := gin.Default()

	routes.RegisterRoutes(router, db, cfg.JWTSecret)

	log.Println("Servidor rodando na porta:", cfg.Port)
	router.Run(":" + cfg.Port)
}
