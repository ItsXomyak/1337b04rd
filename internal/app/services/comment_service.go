package services

import (
	"context"

	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/ports"
	"1337b04rd/internal/domain/comment"
)

type CommentService struct {
	commentRepo ports.CommentPort
	threadRepo  ports.ThreadPort
}

func NewCommentService(commentRepo ports.CommentPort, threadRepo ports.ThreadPort) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		threadRepo:  threadRepo,
	}
}

func (s *CommentService) CreateComment(ctx context.Context,threadID uuidHelper.UUID, parentID *uuidHelper.UUID, content string, imageURL *string, sessionID uuidHelper.UUID) (*comment.Comment, error) {
	c, err := comment.NewComment(threadID, parentID, content, imageURL, sessionID)
	if err != nil {
		return nil, err
	}
	if err := s.commentRepo.CreateComment(ctx, c); err != nil {
		return nil, err
	}

	t, err := s.threadRepo.GetThreadByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if err := s.threadRepo.UpdateThread(ctx, t); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CommentService) GetCommentsByThreadID(ctx context.Context, threadID uuidHelper.UUID) ([]*comment.Comment, error) {
	return s.commentRepo.GetCommentsByThreadID(ctx, threadID)
}