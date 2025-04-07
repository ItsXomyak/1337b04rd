package comment

import "time"

type Comment struct {
	ID        int64 `json:"id"`
	PostID   int64 `json:"post_id"`
	AuthorID int64 `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorName string `json:"-"`
	IsLocked bool `json:"-"`
	IsDeleted bool `json:"-"`
	IsReported bool `json:"-"`
	IsHighlighted bool `json:"-"`
	IsPinned bool `json:"-"`
}
