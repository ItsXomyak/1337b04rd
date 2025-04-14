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
	"1337b04rd/internal/adapters/s3"
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

	// Repos
	sessionRepo := postgres.NewSessionRepository(db)
	threadRepo := postgres.NewThreadRepository(db)
	commentRepo := postgres.NewCommentRepository(db)

	// Clients
	httpClient := &http.Client{}
	avatarClient := rickmorty.NewClient(cfg.AvatarAPI.BaseURL, httpClient)

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

	avatarSvc := services.NewAvatarService(avatarClient)
	sessionSvc := services.NewSessionService(sessionRepo, avatarSvc, cfg.Session.Duration)
	threadSvc := services.NewThreadService(threadRepo, s3ThreadsClient)
	commentSvc := services.NewCommentService(commentRepo, threadRepo, s3CommentsClient)

	router := httpadapter.NewRouter(sessionSvc, avatarSvc, threadSvc, commentSvc)

	addr := fmt.Sprintf(":%d", *port)
	logger.Info("starting server", "address", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
