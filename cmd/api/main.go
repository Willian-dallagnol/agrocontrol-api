package main

import (
	"agrocontrol-api/configs"                  // responsável por carregar configs (env, banco, etc)
	"agrocontrol-api/internal/domain/entities" // entidades que viram tabelas no banco
	"agrocontrol-api/internal/routes"          // onde ficam definidas as rotas da API
	"log"

	"github.com/gin-gonic/gin" // framework HTTP usado para criar a API
)

func main() {

	// 🔧 Carrega configurações do projeto (.env)
	cfg := configs.LoadConfig()

	// 🗄️ Conecta ao banco de dados (PostgreSQL)
	db := configs.ConnectDatabase(cfg)

	// 🔄 Cria/atualiza automaticamente as tabelas no banco
	err := db.AutoMigrate(
		&entities.User{},  // tabela de usuários
		&entities.Farm{},  // tabela de fazendas
		&entities.Field{}, // tabela de talhões
		&entities.Crop{},  // tabela de culturas
	)
	if err != nil {
		// ❌ Se der erro na criação das tabelas, encerra a aplicação
		log.Fatal("Erro ao executar migration:", err)
	}

	// 🌐 Inicializa o servidor HTTP com Gin
	router := gin.Default()

	// 🔗 Registra todas as rotas da aplicação
	// (users, auth, farms, fields, crops)
	routes.RegisterRoutes(router, db, cfg.JWTSecret)

	// 📢 Log indicando que o servidor subiu corretamente
	log.Println("Servidor rodando na porta:", cfg.Port)

	// 🚀 Inicia o servidor na porta definida
	router.Run(":" + cfg.Port)
}
