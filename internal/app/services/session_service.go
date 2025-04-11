package services

import (
	"time"

	"1337b04rd/internal/app/ports"
	"1337b04rd/internal/domain/errors"
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
		return s.CreateNew()
	}

	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return s.CreateNew()
	}

	if sess.IsExpired() {
		return s.CreateNew()
	}

	return sess, nil
}

func (s *SessionService) CreateNew() (*session.Session, error) {
	avatar, err := s.avatarSvc.GetRandomAvatar()
	if err != nil {
		return nil, errors.ErrAvatarAssignment
	}

	newSess, err := session.NewSession(avatar.URL, avatar.DisplayName, s.sessionTTL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateSession(newSess); err != nil {
		return nil, err
	}

	return newSess, nil
}

func (s *SessionService) ListActiveSessions() ([]*session.Session, error) {
	all, err := s.repo.ListActiveSessions()
	if err != nil {
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
	return s.repo.DeleteExpired()
}
