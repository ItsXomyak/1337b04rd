package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"1337b04rd/config"
	"1337b04rd/internal/adapters/postgres"
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
		logger.Error("Failed to connect to DB", "err", err)
		return
	}
	defer db.Close()

	logger.Info("Connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)

	threadRepo := postgres.NewThreadRepository(db)
	commentRepo := postgres.NewCommentRepository(db)
	threadSvc := services.NewThreadService(threadRepo)
	commentSvc := services.NewCommentService(commentRepo, threadRepo)

	handler := http.NewHandler(threadSvc, commentSvc)

	http.HandleFunc("/catalog", handler.CatalogHandler)
	http.HandleFunc("/archive", handler.ArchiveHandler)
	http.HandleFunc("/thread", handler.ThreadHandler)
	http.HandleFunc("/create-thread", handler.CreateThreadHandler)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := threadSvc.CleanupExpiredThreads(); err != nil {
				logger.Error("failed to cleanup expired threads", "error", err)
			}
		}
	}()

	addr := fmt.Sprintf(":%d", *port)
	logger.Info("starting server", "address", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
