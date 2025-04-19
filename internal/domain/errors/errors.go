package errors

import "errors"

var (
	ErrCommentNotFound    = errors.New("comment not found")
	ErrInvalidCommentID   = errors.New("invalid comment ID")
	ErrInvalidParentID    = errors.New("invalid parent comment ID")
	ErrInvalidDisplayName = errors.New("invalid display name")

	ErrThreadNotFound    = errors.New("thread not found")
	ErrInvalidThreadID   = errors.New("invalid thread ID")
	ErrEmptyTitle        = errors.New("thread title cannot be empty")
	ErrEmptyContent      = errors.New("thread content cannot be empty")
	ErrTooLongTitle      = errors.New("thread title is too long")
	ErrTooLongContent    = errors.New("thread content is too long")
	ErrImageUploadFailed = errors.New("failed to upload image for thread")

	ErrInvalidAvatar         = errors.New("avatar not found")
	ErrInvalidUserName       = errors.New("username is invalid")
	ErrFailedToFetchAvatar   = errors.New("failed to fetch avatar from external API")
	ErrNoAvailableAvatars    = errors.New("no available avatars left")
	ErrInvalidAvatarResponse = errors.New("invalid response format from avatar API")
	ErrAvatarAlreadyAssigned = errors.New("avatar already assigned to this session")
	ErrAfterMultipleAttempts = errors.New("failed to get valid avatar after multiple attempts")
	ErrInvalidAvatarURL      = errors.New("invalid avatar URL")

	ErrSessionRequired     = errors.New("session is required to create thread")
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidSessionID    = errors.New("invalid session ID")
	ErrSessionExpired      = errors.New("session expired")
	ErrAvatarAssignment    = errors.New("failed to assign avatar")
	ErrDisplayNameConflict = errors.New("display name already in use")
)
