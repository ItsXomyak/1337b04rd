package session

import "time"

type Session struct {
	ID          string
	AvatarURL   string
	DisplayName string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}
