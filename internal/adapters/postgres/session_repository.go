package postgres

import (
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/errors"
	"1337b04rd/internal/domain/session"
	"context"
	"database/sql"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, s *session.Session) error {
	_, err := r.db.ExecContext(ctx, CreateSession,
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

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*session.Session, error) {
	row := r.db.QueryRowContext(ctx, GetSessionByID, id)

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

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, DeleteExpired, time.Now())
	if err != nil {
		logger.Error("failed to delete expired sessions", "error", err)
	}
	return err
}

func (r *SessionRepository) ListActiveSessions(ctx context.Context) ([]*session.Session, error) {
	rows, err := r.db.QueryContext(ctx, ListActiveSessions)
	if err != nil {
		logger.Error("failed to query active sessions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []*session.Session
	for rows.Next() {
		var s session.Session
		var uuidStr string
		if err := rows.Scan(&uuidStr, &s.AvatarURL, &s.DisplayName, &s.CreatedAt, &s.ExpiresAt); err != nil {
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

func (r *SessionRepository) UpdateDisplayName(ctx context.Context, id string, name string) error {
	_, err := r.db.ExecContext(ctx, UpdateDisplayName, name, id)
	if err != nil {
		logger.Error("failed to update display name", "id", id, "error", err)
	}
	return err
}
