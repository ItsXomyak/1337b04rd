package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"1337b04rd/config"
	httpadapter "1337b04rd/internal/adapters/http"
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/adapters/rickmorty"
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/services"
)

func Run() {
	port := flag.Int("port", 8080, "Port number")
	flag.Parse()

	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("failed to connect to DB", "err", err)
		return
	}
	defer db.Close()
	logger.Info("connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)

	sessionRepo := postgres.NewSessionRepository(db)

	httpClient := &http.Client{}
	avatarClient := rickmorty.NewClient(cfg.AvatarAPI.BaseURL, httpClient)

	avatarSvc := services.NewAvatarService(avatarClient)
	sessionSvc := services.NewSessionService(sessionRepo, avatarSvc, cfg.Session.Duration)

	router := httpadapter.NewRouter(sessionSvc, avatarSvc)

	addr := fmt.Sprintf(":%d", *port)
	logger.Info("starting server", "address", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
