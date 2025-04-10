package ports

import (
	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/comment"
)

type CommentPort interface {
	CreateComment(c *comment.Comment) error
	GetCommentsByThreadID(threadID uuidHelper.UUID) ([]*comment.Comment, error)
}