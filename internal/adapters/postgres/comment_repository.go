package postgres

import (
	"context"
	"database/sql"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/comment"

	"github.com/lib/pq"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, c *comment.Comment) error {
	if err := ctx.Err(); err != nil {
		logger.Error("context error while creating comment", "error", err, "comment_id", c.ID)
		return err
	}

	_, err := r.db.ExecContext(ctx, CreateComment,
		c.ID.String(),
		c.ThreadID.String(),
		nilIfNilUUID(c.ParentCommentID),
		c.Content,
		pq.Array(c.ImageURLs),
		c.SessionID.String(),
		c.CreatedAt,
	)
	if err != nil {
		logger.Error("failed to create comment", "error", err, "comment_id", c.ID)
		return err
	}
	return nil
}

func (r *CommentRepository) GetCommentsByThreadID(ctx context.Context, threadID utils.UUID) ([]*comment.Comment, error) {
	if err := ctx.Err(); err != nil {
		logger.Error("context error while getting comments", "error", err, "thread_id", threadID.String())
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, GetCommentsByThreadID, threadID.String())
	if err != nil {
		logger.Error("failed to query comments by thread id", "error", err, "thread_id", threadID.String())
		return nil, err
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			logger.Error("failed to scan comment", "error", err, "thread_id", threadID.String())
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
	var idStr, threadIDStr, sessionIDStr string
	var parentID sql.NullString
	var imageURLs pq.StringArray

	err := scanner.Scan(
		&idStr,
		&threadIDStr,
		&parentID,
		&c.Content,
		&imageURLs,
		&sessionIDStr,
		&c.CreatedAt,
	)
	if err != nil {
		logger.Error("failed to scan comment row", "error", err)
		return nil, err
	}

	c.ID, err = utils.ParseUUID(idStr)
	if err != nil {
		logger.Error("failed to parse comment id", "error", err)
		return nil, err
	}

	c.ThreadID, err = utils.ParseUUID(threadIDStr)
	if err != nil {
		logger.Error("failed to parse thread id", "error", err)
		return nil, err
	}

	c.SessionID, err = utils.ParseUUID(sessionIDStr)
	if err != nil {
		logger.Error("failed to parse session id", "error", err)
		return nil, err
	}

	if parentID.Valid {
		parsedID, err := utils.ParseUUID(parentID.String)
		if err != nil {
			logger.Error("failed to parse parent comment ID", "error", err, "raw_parent_id", parentID.String)
			return nil, err
		}
		c.ParentCommentID = &parsedID
	}

	c.ImageURLs = []string(imageURLs)
	return c, nil
}

func nilIfNilUUID(u *utils.UUID) interface{} {
	if u == nil {
		return nil
	}
	return u.String()
}
