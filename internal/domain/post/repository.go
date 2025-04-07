package post

type PostRepository interface {
	Create(post *Post) (int64, error)
	GetByID(id int64) (*Post, error)
	Update(post *Post) error
	Delete(id int64) error
	ListByThread(threadID int64) ([]*Post, error)
}
