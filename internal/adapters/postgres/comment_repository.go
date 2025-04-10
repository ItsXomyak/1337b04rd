package postgres

import (
	"database/sql"

	"1337b04rd/internal/domain/comment"
	uuidHelper "1337b04rd/internal/app/common/utils"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(c *comment.Comment) error {
	query := `
		INSERT INTO comments (id, thread_id, parent_comment_id, content, image_url, session_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, c.ID, c.ThreadID, c.ParentCommentID, c.Content, c.ImageURL, c.SessionID, c.CreatedAt)
	return err
}

func (r *CommentRepository) GetCommentsByThreadID(threadID uuidHelper.UUID) ([]*comment.Comment, error) {
	query := `
		SELECT id, thread_id, parent_comment_id, content, image_url, session_id, created_at
		FROM comments
		WHERE thread_id = $1`
	rows, err := r.db.Query(query, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		c := &comment.Comment{}
		var imageURL, parentID sql.NullString
		err := rows.Scan(&c.ID, &c.ThreadID, &parentID, &c.Content, &imageURL, &c.SessionID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		if imageURL.Valid {
			c.ImageURL = &imageURL.String
		}
		if parentID.Valid {
			parsedID := uuidHelper.UUID{} 
			c.ParentCommentID = &parsedID
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}