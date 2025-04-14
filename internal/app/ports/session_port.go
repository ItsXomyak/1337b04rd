package ports

import (
	"1337b04rd/internal/domain/session"
)

type SessionPort interface {
	CreateSession(s *session.Session) error
	DeleteExpired() error
	ListActiveSessions() ([]*session.Session, error)
	UpdateDisplayName(id string, name string) error
	GetSessionByID(id string) (*session.Session, error)
}

