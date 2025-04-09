package comment

import "time"

type Comment struct {
	ID              string
	ThreadID        string
	ParentCommentID *string
	Content         string
	ImageURL        *string
	SessionID       string
	CreatedAt       time.Time
}
