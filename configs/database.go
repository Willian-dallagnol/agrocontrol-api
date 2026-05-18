package configs

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=America/Sao_Paulo",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)

	gormCfg := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		log.Fatal("[DB] Falha ao conectar ao PostgreSQL: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("[DB] Falha ao obter *sql.DB: ", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("db: conectado", "driver", "postgres", "host", cfg.DBHost, "db", cfg.DBName)
	return db
}
