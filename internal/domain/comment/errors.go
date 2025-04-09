package comment

import "errors"

var (
	ErrCommentNotFound   = errors.New("comment not found")
	ErrInvalidCommentID  = errors.New("invalid comment ID")
	ErrEmptyContent      = errors.New("comment content cannot be empty")
	ErrTooLongContent    = errors.New("comment content is too long")
	ErrImageUploadFailed = errors.New("failed to upload image for comment")
	ErrInvalidThreadID   = errors.New("invalid thread ID for comment")
	ErrInvalidParentID   = errors.New("invalid parent comment ID")
	ErrSessionRequired   = errors.New("session is required to create comment")
)
