package configs

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	JWTSecret   string
	JWTExpHours int
	RedisAddr   string
}

func LoadConfig() *Config {
	// Em produção, variáveis vêm do ambiente; .env só para desenvolvimento local
	if err := godotenv.Load(); err != nil {
		log.Println("[CONFIG] .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := &Config{
		Port:        getEnv("APP_PORT", "8080"),
		Env:         getEnv("APP_ENV", "production"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPass:      getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "agro_control"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTExpHours: getEnvInt("JWT_EXP_HOURS", 24),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}

	validateConfig(cfg)
	return cfg
}

func validateConfig(cfg *Config) {
	var errs []string

	if cfg.JWTSecret == "" {
		errs = append(errs, "JWT_SECRET não pode ser vazio — defina no .env")
	}
	if len(cfg.JWTSecret) < 32 {
		errs = append(errs, "JWT_SECRET deve ter pelo menos 32 caracteres")
	}
	if cfg.DBPass == "" {
		errs = append(errs, "DB_PASSWORD não pode ser vazio")
	}

	if len(errs) > 0 {
		log.Fatal("[CONFIG] Configuração inválida:\n  - " + strings.Join(errs, "\n  - "))
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
