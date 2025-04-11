package postgres

import (
	"database/sql"
	"time"

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
			return nil, errors.ErrSessionNotFound
		}
		return nil, err
	}

	uid, err := utils.ParseUUID(uuidStr)
	if err != nil {
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
	return err
}

func (r *SessionRepository) ListActiveSessions() ([]*session.Session, error) {
	query := `
		SELECT id, avatar_url, display_name, created_at, expires_at
		FROM sessions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*session.Session
	for rows.Next() {
		var s session.Session
		var uuidStr string
		err := rows.Scan(&uuidStr, &s.AvatarURL, &s.DisplayName, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}
		uid, err := utils.ParseUUID(uuidStr)
		if err != nil {
			return nil, err
		}
		s.ID = uid
		sessions = append(sessions, &s)
	}

	return sessions, rows.Err()
}
