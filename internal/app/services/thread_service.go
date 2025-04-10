package services

import (
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

func (s *ThreadService) CreateThread(title, content string, imageURL *string, sessionID uuidHelper.UUID) (*thread.Thread, error) {
	t, err := thread.NewThread(title, content, imageURL, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.threadRepo.CreateThread(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *ThreadService) GetThreadByID(id uuidHelper.UUID) (*thread.Thread, error) {
	t, err := s.threadRepo.GetThreadByID(id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *ThreadService) ListActiveThreads() ([]*thread.Thread, error) {
	threads, err := s.threadRepo.ListActiveThreads()
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

func (s *ThreadService) ListAllThreads() ([]*thread.Thread, error) {
	return s.threadRepo.ListAllThreads()
}

func (s *ThreadService) CleanupExpiredThreads() error {
	threads, err := s.threadRepo.ListActiveThreads()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, t := range threads {
		if t.ShouldDelete(now) {
			t.MarkAsDeleted()
			if err := s.threadRepo.UpdateThread(t); err != nil {
				return err
			}
		}
	}
	return nil
}