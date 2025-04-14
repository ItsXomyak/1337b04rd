package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	DB     struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}
	AvatarAPI struct {
		BaseURL string
	}
	Session struct {
		CookieName string
		Duration   time.Duration
	}
	S3 struct {
		Endpoint      string
		AccessKey     string
		SecretKey     string
		PostBucket    string
		CommentBucket string
		Region        string
		UseSSL        bool
	}
}

func Load() Config {
	godotenv.Load()

	cfg := Config{
		AppEnv: os.Getenv("APP_ENV"),
	}
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.SSLMode = os.Getenv("DB_SSLMODE")
	cfg.AvatarAPI.BaseURL = os.Getenv("AVATAR_API_BASE_URL")
	cfg.Session.CookieName = os.Getenv("SESSION_COOKIE_NAME")
	durationDays, _ := strconv.Atoi(os.Getenv("SESSION_DURATION_DAYS"))
	cfg.Session.Duration = time.Duration(durationDays) * 24 * time.Hour
	cfg.S3.Endpoint = os.Getenv("S3_ENDPOINT")
	cfg.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	cfg.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	cfg.S3.PostBucket = os.Getenv("S3_BUCKET_THREADS")
	cfg.S3.CommentBucket = os.Getenv("S3_BUCKET_POSTS")
	cfg.S3.Region = os.Getenv("S3_REGION")
	useSSL, _ := strconv.ParseBool(os.Getenv("S3_USE_SSL"))
	cfg.S3.UseSSL = useSSL

	return cfg
}