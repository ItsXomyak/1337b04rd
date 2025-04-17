package services

import (
	"errors"
	"time"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/ports"
	"1337b04rd/internal/domain/session"
)

type SessionService struct {
	repo       ports.SessionPort
	avatarSvc  ports.AvatarPort
	sessionTTL time.Duration
}

func NewSessionService(repo ports.SessionPort, avatarSvc ports.AvatarPort, ttl time.Duration) *SessionService {
	return &SessionService{
		repo:       repo,
		avatarSvc:  avatarSvc,
		sessionTTL: ttl,
	}
}

func (s *SessionService) GetOrCreate(sessionID string) (*session.Session, error) {
	if sessionID == "" {
		logger.Warn("no session ID provided")
		return nil, errors.New("session ID not provided")
	}

	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		logger.Warn("session not found", "sessionID", sessionID, "error", err)
		return nil, err
	}

	if sess.IsExpired() {
		return nil, err
	}

	return sess, nil
}

func (s *SessionService) CreateNew() (*session.Session, error) {
	avatar, err := s.avatarSvc.GetRandomAvatar()
	if err != nil {
		logger.Error("failed to assign avatar", "error", err)
		return nil, err
	}

	newSess, err := session.NewSession(avatar.URL, avatar.DisplayName, s.sessionTTL)
	if err != nil {
		logger.Error("failed to create new session object", "error", err)
		return nil, err
	}

	if err := s.repo.CreateSession(newSess); err != nil {
		logger.Error("failed to persist session", "id", newSess.ID, "error", err)
		return nil, err
	}

	logger.Info("new session created", "id", newSess.ID, "display_name", newSess.DisplayName)
	return newSess, nil
}

func (s *SessionService) ListActiveSessions() ([]*session.Session, error) {
	all, err := s.repo.ListActiveSessions()
	if err != nil {
		logger.Error("failed to list sessions", "error", err)
		return nil, err
	}

	var result []*session.Session
	for _, sess := range all {
		if !sess.IsExpired() {
			result = append(result, sess)
		}
	}
	return result, nil
}

func (s *SessionService) DeleteExpired() error {
	err := s.repo.DeleteExpired()
	if err != nil {
		logger.Error("failed to delete expired sessions", "error", err)
	}
	return err
}

func (s *SessionService) UpdateDisplayName(id utils.UUID, newName string) error {
	err := s.repo.UpdateDisplayName(id.String(), newName)
	if err != nil {
		logger.Error("failed to update display name", "id", id.String(), "error", err)
	}
	return err
}
