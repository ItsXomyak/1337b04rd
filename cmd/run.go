package cmd

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
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

	db, err := postgres.NewPostgresDB(&cfg)
	if err != nil {
		logger.Error("failed to connect to DB", "err", err)
		return
	}
	defer db.Close()
	logger.Info("connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)

	s3Logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s3Client, err := s3.NewS3Client(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.PostBucket,    // Например, "posts"
		cfg.S3.CommentBucket, // Например, "comments"
		s3Logger,
	)
	if err != nil {
		logger.Error("failed to initialize S3 client", "error", err)
		os.Exit(1)
	}

	sessionRepo := postgres.NewSessionRepository(db)
	threadRepo := postgres.NewThreadRepository(db)
	commentRepo := postgres.NewCommentRepository(db)

	httpClient := &http.Client{}
	avatarClient := rickmorty.NewClient(cfg.AvatarAPI.BaseURL, httpClient)

	avatarSvc := services.NewAvatarService(avatarClient)
	sessionSvc := services.NewSessionService(sessionRepo, avatarSvc, cfg.Session.Duration)
	threadSvc := services.NewThreadService(threadRepo)
	commentSvc := services.NewCommentService(commentRepo, threadRepo)

	if err := httpadapter.InitTemplates(); err != nil {
		log.Fatal("Failed to initialize templates:", err)
		os.Exit(1)
	}

	router := httpadapter.NewRouter(sessionSvc, avatarSvc, threadSvc, commentSvc, s3Client)

	addr := fmt.Sprintf(":%d", *port)
	logger.Info("starting server", "address", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}