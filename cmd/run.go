package cmd

import (
	"1337b04rd/config"
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/common/logger"
)

func Run() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("Failed to connect to DB", "err", err)
		return
	}
	defer db.Close()

	logger.Info("Connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)
}
