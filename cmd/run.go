package cmd

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"1337b04rd/config"
	httpadapter "1337b04rd/internal/adapters/http"
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/adapters/rickmorty"
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/services"
)

func Run() {
	// CLI flags
	port := flag.Int("port", 8080, "Port number")
	flag.Parse()

	// Load config and init logger
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	// Connect to PostgreSQL
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("failed to connect to DB", "err", err)
		return
	}
	defer db.Close()
	logger.Info("connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)

	// Repositories
	sessionRepo := postgres.NewSessionRepository(db)
	threadRepo := postgres.NewThreadRepository(db)
	commentRepo := postgres.NewCommentRepository(db)

	// External HTTP clients
	httpClient := &http.Client{}
	avatarClient := rickmorty.NewClient(cfg.AvatarAPI.BaseURL, httpClient)

	// S3 clients for threads and comments
	s3ThreadsClient, err := s3.NewS3Client(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.BucketThreads,
	)
	if err != nil {
		logger.Error("failed to create S3 client for threads", "error", err)
		return
	}

	s3CommentsClient, err := s3.NewS3Client(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.BucketComments,
	)
	if err != nil {
		logger.Error("failed to create S3 client for comments", "error", err)
		return
	}

	// Services
	avatarSvc := services.NewAvatarService(avatarClient)
	sessionSvc := services.NewSessionService(sessionRepo, avatarSvc, cfg.Session.Duration)

	threadS3Adapter := s3.NewAdapter(s3ThreadsClient)
	commentS3Adapter := s3.NewAdapter(s3CommentsClient)

	threadSvc := services.NewThreadService(threadRepo, threadS3Adapter)
	commentSvc := services.NewCommentService(commentRepo, threadRepo, commentS3Adapter, sessionRepo)

	// HTTP router
	router := httpadapter.NewRouter(sessionSvc, avatarSvc, threadSvc, commentSvc)
	corsRouter := withCORS(router)

	// запуск фонового удаления
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()
			if err := threadSvc.CleanupExpiredThreads(ctx); err != nil {
				logger.Error("thread cleanup failed", "error", err)
			}
		}
	}()

	addr := fmt.Sprintf(":%d", *port)
	logger.Info("starting server", "address", addr)

	if err := http.ListenAndServe(addr, corsRouter); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "http://localhost:5500" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
