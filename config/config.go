package config

import (
	"log"
	"os"
	"strconv"
	"strings"
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
		Endpoint       string
		AccessKey      string
		SecretKey      string
		BucketComments string
		BucketThreads  string
		Region         string
		UseSSL         bool
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
	loadDotEnv(".env")

	cfg := &Config{}

	// Port
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
	cfg.S3.BucketComments = mustGet("S3_BUCKET_COMMENTS")
	cfg.S3.Region = mustGet("S3_REGION")
	cfg.S3.UseSSL = getBool("S3_USE_SSL")

	if cfg.S3.Endpoint == "minio:9000" || cfg.S3.Endpoint == "http://minio:9000" {
		if _, err := os.Stat("/.dockerenv"); err != nil {
			cfg.S3.Endpoint = "http://localhost:9000"
		} else {
			cfg.S3.Endpoint = "http://minio:9000"
		}
	}

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

// === .env loader using stdlib ===

func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		_ = os.Setenv(key, val)
	}
}
