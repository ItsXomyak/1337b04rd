package postgres

import (
	"context"
	"database/sql"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/comment"
)

type CommentRepository struct {
	db     *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db:     db,
	}
}


func (r *CommentRepository) CreateComment(ctx context.Context, c *comment.Comment) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, CreateComment,
		c.ID,
		c.ThreadID,
		c.ParentCommentID,
		c.Content,
		c.ImageURL,
		c.SessionID,
		c.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepository) GetCommentsByThreadID(ctx context.Context, threadID utils.UUID) ([]*comment.Comment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, GetCommentsByThreadID, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			logger.Error("failed to scan comment", "error", err, "thread_id", threadID)
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		logger.Error("error in comment rows", "error", err, "thread_id", threadID)
		return nil, err
	}
	return comments, nil
}

func scanComment(scanner interface {
	Scan(dest ...interface{}) error
}) (*comment.Comment, error) {
	c := &comment.Comment{}
	var imageURL, parentID sql.NullString

	err := scanner.Scan(
		&c.ID,
		&c.ThreadID,
		&parentID,
		&c.Content,
		&imageURL,
		&c.SessionID,
		&c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if imageURL.Valid {
		c.ImageURL = &imageURL.String
	}
	if parentID.Valid {
		parsedID, err := utils.ParseUUID(parentID.String)
		if err != nil {
			return nil, err
		}
		c.ParentCommentID = &parsedID
	}

	return c, nil
}