package session

import "errors"

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidSessionID    = errors.New("invalid session ID")
	ErrSessionExpired      = errors.New("session expired")
	ErrAvatarAssignment    = errors.New("failed to assign avatar")
	ErrDisplayNameConflict = errors.New("display name already in use")
)
