package thread

import "errors"

var (
	ErrThreadNotFound    = errors.New("thread not found")
	ErrInvalidThreadID   = errors.New("invalid thread ID")
	ErrEmptyTitle        = errors.New("thread title cannot be empty")
	ErrEmptyContent      = errors.New("thread content cannot be empty")
	ErrTooLongTitle      = errors.New("thread title is too long")
	ErrTooLongContent    = errors.New("thread content is too long")
	ErrImageUploadFailed = errors.New("failed to upload image for thread")
	ErrSessionRequired   = errors.New("session is required to create thread")
)
