package services

import (
	"context"
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/ports"
	"1337b04rd/internal/domain/thread"
)

type ThreadService struct {
	threadRepo ports.ThreadPort
}

func NewThreadService(threadRepo ports.ThreadPort) *ThreadService {
	return &ThreadService{threadRepo: threadRepo}
}

func (s *ThreadService) CreateThread(ctx context.Context, title, content string, imageURL *string, sessionID uuidHelper.UUID) (*thread.Thread, error) {
	t, err := thread.NewThread(title, content, imageURL, sessionID)
	if err != nil {  // здесь ебануть слог
		return nil, err
	}

	if err := s.threadRepo.CreateThread(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *ThreadService) GetThreadByID(ctx context.Context, id uuidHelper.UUID) (*thread.Thread, error) {
	t, err := s.threadRepo.GetThreadByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *ThreadService) ListActiveThreads(ctx context.Context) ([]*thread.Thread, error) {
	threads, err := s.threadRepo.ListActiveThreads(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var activeThreads []*thread.Thread
	for _, t := range threads {
		if !t.ShouldDelete(now) {
			activeThreads = append(activeThreads, t)
		}
	}
	return activeThreads, nil
}

func (s *ThreadService) ListAllThreads(ctx context.Context) ([]*thread.Thread, error) {
	return s.threadRepo.ListAllThreads(ctx)
}

func (s *ThreadService) CleanupExpiredThreads(ctx context.Context) error {
    threads, err := s.threadRepo.ListActiveThreads(ctx)
    if err != nil {
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