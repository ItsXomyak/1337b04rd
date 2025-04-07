package comment

type ComRepository interface {
	Create(comment *Comment) (int64, error)
	GetByID(id int64) (*Comment, error)
	ListByPost(postID int64) ([]*Comment, error)
	Update(comment *Comment) error
	Delete(id int64) error
}
