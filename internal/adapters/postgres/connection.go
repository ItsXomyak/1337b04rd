package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"1337b04rd/config"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
dsn := fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	cfg.DB.Host,
	cfg.DB.Port,
	cfg.DB.User,
	cfg.DB.Password,
	cfg.DB.Name,
	cfg.DB.SSLMode,
)


	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.: %w", err)
	}

	return db, nil
}
