package services

import (
	"context"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
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
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in CreateComment", "error", err)
		return nil, err
	}

	c, err := comment.NewComment(threadID, parentID, content, imageURL, sessionID)
	if err != nil {
		logger.Error("cannot to create new comment", "error", err)
		return nil, err
	}
	if err := s.commentRepo.CreateComment(ctx, c); err != nil {
		logger.Error("cannot to create comment", "error", err)
		return nil, err
	}

	t, err := s.threadRepo.GetThreadByID(ctx, threadID)
	if err != nil {
		logger.Error("cannot to get thread by ID", "error", err)
		return nil, err
	}
	if err := s.threadRepo.UpdateThread(ctx, t); err != nil {
		logger.Error("cannot to update thread", "error", err)
		return nil, err
	}

	logger.Info("a new comment has been created!", "comment", c)
	return c, nil
}

func (s *CommentService) GetCommentsByThreadID(ctx context.Context, threadID utils.UUID) ([]*comment.Comment, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in GetCommentsByThreadID", "error", err)
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByThreadID(ctx, threadID)
	if err != nil {
		logger.Error("failed to get comments", "error", err, "thread_id", threadID)
		return nil, err
	}

	logger.Info("comments retrieved", "thread_id", threadID, "count", len(comments))
	return comments, nil
}