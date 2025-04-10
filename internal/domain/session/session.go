package session

import (
	. "1337b04rd/internal/domain/errors"
	"time"
)

type Session struct {
	ID          string
	AvatarURL   string
	DisplayName string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}
