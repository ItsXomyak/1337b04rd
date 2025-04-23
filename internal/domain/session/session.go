package session

import (
	"1337b04rd/internal/domain/errors"
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
)

type Session struct {
	ID          uuidHelper.UUID
	AvatarURL   string
	DisplayName string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

func NewSession(avatarURL, displayName string, duration time.Duration) (*Session, error) {
	if avatarURL == "" {
		return nil, errors.ErrInvalidAvatar
	}

	if displayName == "" {
		return nil, errors.ErrInvalidUserName
	}

	id, err := uuidHelper.NewUUID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Session{
		ID:          id,
		AvatarURL:   avatarURL,
		DisplayName: displayName,
		CreatedAt:   now,
		ExpiresAt:   now.Add(duration),
	}, nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
