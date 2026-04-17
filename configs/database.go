package configs

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres" // driver do PostgreSQL para o GORM
	"gorm.io/gorm"            // ORM utilizado para acesso ao banco
)

// 🔗 Função responsável por conectar no banco de dados
func ConnectDatabase(cfg *Config) *gorm.DB {

	// 🧩 Monta a string de conexão (DSN) com base nas configurações
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, // host do banco
		cfg.DBUser, // usuário
		cfg.DBPass, // senha
		cfg.DBName, // nome do banco
		cfg.DBPort, // porta
	)

	// 🚀 Abre conexão com o banco usando GORM + PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// ❌ Se não conseguir conectar, encerra a aplicação
		log.Fatal("Erro ao conectar no banco:", err)
	}

	// ✅ Log indicando conexão bem sucedida
	log.Println("Banco conectado com sucesso")

	// 🔄 Retorna a conexão para ser usada no restante da aplicação
	return db
}
