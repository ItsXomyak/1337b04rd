package thread

import "time"

type Thread struct {
	ID            string
	Title         string
	Content       string
	ImageURL      *string
	SessionID     string
	CreatedAt     time.Time
	LastCommented *time.Time
	IsDeleted     bool
}
