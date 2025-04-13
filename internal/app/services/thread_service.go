package services

import (
	"context"
	"io"
	"time"

	"1337b04rd/internal/app/common/logger"
	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/ports"
	"1337b04rd/internal/domain/thread"
)

type ThreadService struct {
	threadRepo ports.ThreadPort
	s3         ports.S3Port
}

func NewThreadService(threadRepo ports.ThreadPort, s3 ports.S3Port) *ThreadService {
	return &ThreadService{threadRepo: threadRepo, s3: s3}
}

func (s *ThreadService) CreateThread(
	ctx context.Context,
	title, content string,
	image io.Reader,
	contentType string,
	sessionID uuidHelper.UUID,
) (*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in CreateThread", "error", err)
		return nil, err
	}

	var imageURL *string
	if image != nil {
		url, err := s.s3.UploadImage(image, 0, contentType)
		if err != nil {
			logger.Error("failed to upload image", "error", err)
			return nil, err
		}
		imageURL = &url
	}

	t, err := thread.NewThread(title, content, imageURL, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.threadRepo.CreateThread(ctx, t); err != nil {
		logger.Error("failed to create new thread", "error", err)
		return nil, err
	}

	logger.Info("Succesful! New thread created", "thread", t)
	return t, nil
}

func (s *ThreadService) GetThreadByID(ctx context.Context, id uuidHelper.UUID) (*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in GetThreadByID", "error", err)
		return nil, err
	}

	t, err := s.threadRepo.GetThreadByID(ctx, id)
	if err != nil {
		logger.Error("cannot to get thread by id", "error", err)
		return nil, err
	}
	logger.Info("received a thread by ID", "thread", t, "id", id)
	return t, nil
}

func (s *ThreadService) ListActiveThreads(ctx context.Context) ([]*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in ListActiveThreads", "error", err)
		return nil, err
	}

	threads, err := s.threadRepo.ListActiveThreads(ctx)
	if err != nil {
		logger.Error("couldn't get a list of active threads", "error", err)
		return nil, err
	}
	now := time.Now()
	var activeThreads []*thread.Thread
	for _, t := range threads {
		if !t.ShouldDelete(now) {
			activeThreads = append(activeThreads, t)
		}
	}
	logger.Info("the list of active threads has been received", "activeThreads", activeThreads)
	return activeThreads, nil
}

func (s *ThreadService) ListAllThreads(ctx context.Context) ([]*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in ListAllThreads", "error", err)
		return nil, err
	}

	threads, err := s.threadRepo.ListAllThreads(ctx)
	if err != nil {
		logger.Error("failed to list all threads", "error", err)
		return nil, err
	}

	logger.Info("list of all threads retrieved", "count", len(threads))
	return threads, nil
}

func (s *ThreadService) CleanupExpiredThreads(ctx context.Context) error {
	threads, err := s.threadRepo.ListActiveThreads(ctx)
	if err != nil {
		logger.Error("cannot get a list of active threads", "error", err)
		return err
	}

	now := time.Now()
	var lastErr error
	for _, t := range threads {
		if t.ShouldDelete(now) {
			t.MarkAsDeleted()
			if err := s.threadRepo.UpdateThread(ctx, t); err != nil {
				lastErr = err
				continue
			}
		}
	}
	return lastErr
}
