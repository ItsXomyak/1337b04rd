package post

import "time"

type Post struct {
	ID        int64     `json:"id"`
	ThreadID  int64     `json:"thread_id"`
	AuthorID  int64     `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}