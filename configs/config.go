package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv" // biblioteca para carregar variáveis do arquivo .env
)

// 📦 Estrutura que guarda todas as configurações da aplicação
type Config struct {
	Port      string // porta que a API vai rodar
	DBHost    string // host do banco (ex: localhost)
	DBPort    string // porta do banco
	DBUser    string // usuário do banco
	DBPass    string // senha do banco
	DBName    string // nome do banco
	JWTSecret string // chave usada para gerar/validar tokens JWT
}

// 🔧 Função responsável por carregar as configurações do sistema
func LoadConfig() *Config {

	// 📄 Tenta carregar o arquivo .env
	err := godotenv.Load()
	if err != nil {
		// ⚠️ Se não encontrar o .env, apenas avisa (não quebra o sistema)
		log.Println("No .env file found")
	}

	// 📥 Retorna todas as configurações, buscando do ambiente
	return &Config{
		Port:      getEnv("APP_PORT", "8080"),        // porta da aplicação
		DBHost:    getEnv("DB_HOST", "localhost"),    // host do banco
		DBPort:    getEnv("DB_PORT", "5432"),         // porta do banco
		DBUser:    getEnv("DB_USER", "postgres"),     // usuário padrão
		DBPass:    getEnv("DB_PASSWORD", "postgres"), // senha padrão
		DBName:    getEnv("DB_NAME", "agro_control"), // nome do banco
		JWTSecret: getEnv("JWT_SECRET", "secret"),    // chave JWT
	}
}

// 🔎 Função auxiliar para buscar variável de ambiente
func getEnv(key string, fallback string) string {

	// tenta pegar valor da variável de ambiente
	value := os.Getenv(key)

	// se não existir, usa valor padrão (fallback)
	if value == "" {
		return fallback
	}

	// retorna valor encontrado
	return value
}
