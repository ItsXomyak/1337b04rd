package comment

import (
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
	. "1337b04rd/internal/domain/errors"
)

type Comment struct {
    ID              uuidHelper.UUID  `json:"ID"`
    ThreadID        uuidHelper.UUID  `json:"ThreadID"`
    ParentCommentID *uuidHelper.UUID `json:"ParentCommentID"`
    Content         string           `json:"Content"`
    ImageURLs       []string         `json:"ImageURLs"`
    SessionID       uuidHelper.UUID  `json:"SessionID"`
    CreatedAt       time.Time        `json:"CreatedAt"`
    IsDeleted       bool             `json:"IsDeleted"`
    DisplayName     string           `json:"display_name"`
    AvatarURL       string           `json:"avatar_url"`  
}

func NewComment(threadID uuidHelper.UUID, parentCommentID *uuidHelper.UUID, content string, imageURLs []string, sessionID uuidHelper.UUID, DisplayName string, AvatarURL string) (*Comment, error) {
	if threadID.IsZero() {
		return nil, ErrInvalidThreadID
	}
	if content == "" {
		return nil, ErrEmptyContent
	}
	if sessionID.IsZero() {
		return nil, ErrInvalidSessionID
	}
	if DisplayName == "" {
		return nil, ErrInvalidDisplayName
	}
	if AvatarURL == "" {
		return nil, ErrInvalidAvatarURL
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
		ImageURLs:       imageURLs,
		SessionID:       sessionID,
		CreatedAt:       time.Now(),
		IsDeleted:       false,
		DisplayName:  	 DisplayName,
		AvatarURL:  	   AvatarURL,
	}, nil
}

func (c *Comment) MarkAsDeleted() {
	c.IsDeleted = true
}
