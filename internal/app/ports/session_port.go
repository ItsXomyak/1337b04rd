package ports

import (
	"context"

	"1337b04rd/internal/domain/session"
)

type SessionPort interface {
	GetSessionByID(ctx context.Context, id string) (*session.Session, error)
	CreateSession(ctx context.Context, s *session.Session) error
	DeleteExpired(ctx context.Context) error
	ListActiveSessions(ctx context.Context) ([]*session.Session, error)
	UpdateDisplayName(ctx context.Context, id string, name string) error
}
