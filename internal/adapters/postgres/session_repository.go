package postgres

import (
	"database/sql"
	"time"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/errors"
	"1337b04rd/internal/domain/session"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(s *session.Session) error {
	query := `
		INSERT INTO sessions (id, avatar_url, display_name, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query,
		s.ID.String(),
		s.AvatarURL,
		s.DisplayName,
		s.CreatedAt,
		s.ExpiresAt,
	)
	if err != nil {
		logger.Error("failed to insert session", "error", err, "id", s.ID.String())
	}
	return err
}

func (r *SessionRepository) GetSessionByID(id string) (*session.Session, error) {
	query := `
		SELECT id, avatar_url, display_name, created_at, expires_at
		FROM sessions
		WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var s session.Session
	var uuidStr string
	err := row.Scan(&uuidStr, &s.AvatarURL, &s.DisplayName, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("session not found", "id", id)
			return nil, errors.ErrSessionNotFound
		}
		logger.Error("failed to scan session", "error", err, "id", id)
		return nil, err
	}

	uid, err := utils.ParseUUID(uuidStr)
	if err != nil {
		logger.Error("invalid UUID string from DB", "uuidStr", uuidStr, "error", err)
		return nil, err
	}
	s.ID = uid
	return &s, nil
}

func (r *SessionRepository) DeleteExpired() error {
	query := `
		DELETE FROM sessions
		WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		logger.Error("failed to delete expired sessions", "error", err)
	}
	return err
}

func (r *SessionRepository) ListActiveSessions() ([]*session.Session, error) {
	query := `
		SELECT id, avatar_url, display_name, created_at, expires_at
		FROM sessions`
	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error("failed to query active sessions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []*session.Session
	for rows.Next() {
		var s session.Session
		var uuidStr string
		err := rows.Scan(&uuidStr, &s.AvatarURL, &s.DisplayName, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			logger.Error("failed to scan session row", "error", err)
			return nil, err
		}
		uid, err := utils.ParseUUID(uuidStr)
		if err != nil {
			logger.Error("failed to parse UUID", "uuidStr", uuidStr, "error", err)
			return nil, err
		}
		s.ID = uid
		sessions = append(sessions, &s)
	}

	if err := rows.Err(); err != nil {
		logger.Error("rows iteration error", "error", err)
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) UpdateDisplayName(id string, name string) error {
	query := `UPDATE sessions SET display_name = $1 WHERE id = $2`
	_, err := r.db.Exec(query, name, id)
	if err != nil {
		logger.Error("failed to update display name", "id", id, "error", err)
	}
	return err
}
