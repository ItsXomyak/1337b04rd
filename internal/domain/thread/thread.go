package thread

import (
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
	. "1337b04rd/internal/domain/errors"
)

type Thread struct {
	ID            uuidHelper.UUID
	Title         string
	Content       string
	ImageURLs     []string
	SessionID     uuidHelper.UUID
	CreatedAt     time.Time
	LastCommented *time.Time
	IsDeleted     bool
}

func NewThread(title, content string, imageURLs []string, sessionID uuidHelper.UUID) (*Thread, error) {
	if title == "" {
		return nil, ErrEmptyTitle
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

	now := time.Now()
	return &Thread{
		ID:            id,
		Title:         title,
		Content:       content,
		ImageURLs:     imageURLs,
		SessionID:     sessionID,
		CreatedAt:     now,
		LastCommented: nil,
		IsDeleted:     false,
	}, nil
}

func (t *Thread) ShouldDelete(now time.Time) bool {
	if t.IsDeleted {
		return false
	}

	timeSinceCreation := now.Sub(t.CreatedAt)

	if t.LastCommented == nil && timeSinceCreation > 10*time.Minute {
		return true
	}

	if t.LastCommented != nil && now.Sub(*t.LastCommented) > 15*time.Minute {
		return true
	}

	return false
}

func (t *Thread) MarkAsDeleted() {
	t.IsDeleted = true
}
