package comment

import (
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
	. "1337b04rd/internal/domain/errors"
)

type Comment struct {
	ID              uuidHelper.UUID
	ThreadID        uuidHelper.UUID
	ParentCommentID *uuidHelper.UUID
	Content         string
	ImageURL        *string
	SessionID       uuidHelper.UUID
	CreatedAt       time.Time
	IsDeleted       bool
	DisplayName string 
	AvatarURL   string 
}

func NewComment(threadID uuidHelper.UUID, parentCommentID *uuidHelper.UUID, content string, imageURL *string, sessionID uuidHelper.UUID) (*Comment, error) {
	if threadID.IsZero() {
		return nil, ErrInvalidThreadID
	}
	if content == "" {
		return nil, ErrEmptyContent
	}
	if sessionID.IsZero() {
		return nil, ErrInvalidSessionID
	}

	id, err := uuidHelper.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Comment{
		ID:              id, 
		ThreadID:        threadID,
		ParentCommentID: parentCommentID,
		Content:         content,
		ImageURL:        imageURL,
		SessionID:       sessionID,
		CreatedAt:       time.Now(),
		IsDeleted:       false,
	}, nil
}

func (c *Comment) MarkAsDeleted() {
	c.IsDeleted = true
}

