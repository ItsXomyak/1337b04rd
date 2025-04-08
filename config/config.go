package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port int

	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	S3 struct {
		Endpoint      string
		AccessKey     string
		SecretKey     string
		BucketPosts   string
		BucketThreads string
		Region        string
		UseSSL        bool
	}

	Session struct {
		CookieName string
		Duration   time.Duration
	}

	AvatarAPI struct {
		BaseURL string
	}

	AppEnv string
}

func Load() *Config {
	cfg := &Config{}

	cfg.Port = mustGetInt("PORT")

	// DB config
	cfg.DB.Host = mustGet("DB_HOST")
	cfg.DB.Port = mustGetInt("DB_PORT")
	cfg.DB.User = mustGet("DB_USER")
	cfg.DB.Password = mustGet("DB_PASSWORD")
	cfg.DB.Name = mustGet("DB_NAME")
	cfg.DB.SSLMode = getOrDefault("DB_SSLMODE", "disable")

	// S3
	cfg.S3.Endpoint = mustGet("S3_ENDPOINT")
	cfg.S3.AccessKey = mustGet("S3_ACCESS_KEY")
	cfg.S3.SecretKey = mustGet("S3_SECRET_KEY")
	cfg.S3.BucketThreads = mustGet("S3_BUCKET_THREADS")
	cfg.S3.BucketPosts = mustGet("S3_BUCKET_POSTS")
	cfg.S3.Region = mustGet("S3_REGION")
	cfg.S3.UseSSL = getBool("S3_USE_SSL")

	// Session
	cfg.Session.CookieName = getOrDefault("SESSION_COOKIE_NAME", "1337session")
	cfg.Session.Duration = time.Hour * 24 * time.Duration(mustGetInt("SESSION_DURATION_DAYS"))

	// Avatar API
	cfg.AvatarAPI.BaseURL = mustGet("AVATAR_API_BASE_URL")

	// App env
	cfg.AppEnv = getOrDefault("APP_ENV", "development")

	return cfg
}

// === helpers ===

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required env var: %s", key)
	}
	return val
}

func mustGetInt(key string) int {
	val := mustGet(key)
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid integer value for %s: %s", key, val)
	}
	return n
}

func getOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getBool(key string) bool {
	val := os.Getenv(key)
	return val == "true" || val == "1"
}
