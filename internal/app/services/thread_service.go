package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
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
	files map[string]io.Reader,
	contentTypes map[string]string,
	sessionID uuidHelper.UUID,
) (*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in CreateThread", "error", err)
		return nil, err
	}

	var imageURLs []string
	if len(files) > 0 {
		urls, err := s.s3.UploadImages(files, contentTypes)
		if err != nil {
			logger.Error("failed to upload thread images", "error", err)
			return nil, err
		}

		// собираем в слайс
		for _, url := range urls {
			imageURLs = append(imageURLs, url)
		}
	}

	t, err := thread.NewThread(title, content, imageURLs, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.threadRepo.CreateThread(ctx, t); err != nil {
		logger.Error("failed to create new thread", "error", err)
		return nil, err
	}

	logger.Info("Successful! New thread created", "thread", t)
	return t, nil
}

func (s *ThreadService) PrepareFilesFromMultipart(form *multipart.Form) (map[string]io.Reader, map[string]string, error) {
	files := make(map[string]io.Reader)
	contentTypes := make(map[string]string)
	counter := 0

	if form == nil || form.File == nil {
		return files, contentTypes, nil
	}

	for _, fileHeaders := range form.File {
		for _, fh := range fileHeaders {
			file, err := fh.Open()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, file); err != nil {
				return nil, nil, fmt.Errorf("failed to buffer file: %w", err)
			}

			key := fmt.Sprintf("file_%d", counter)
			files[key] = bytes.NewReader(buf.Bytes())
			contentTypes[key] = fh.Header.Get("Content-Type")
			counter++
		}
	}

	return files, contentTypes, nil
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
		createdAgo := now.Sub(t.CreatedAt)

		switch {
		// 1. Удалить, если прошло >10 минут и комментариев не было
		case t.LastCommented == nil && createdAgo > 10*time.Minute:
			t.MarkAsDeleted()
			logger.Info("deleting thread with no comments after 10 minutes", "thread_id", t.ID)

		// 2. Удалить, если после последнего комментария прошло >15 минут
		case t.LastCommented != nil && now.Sub(*t.LastCommented) > 15*time.Minute:
			t.MarkAsDeleted()
			logger.Info("deleting thread after 15 minutes of inactivity", "thread_id", t.ID)

		default:
			continue
		}

		if err := s.threadRepo.UpdateThread(ctx, t); err != nil {
			logger.Error("failed to update (delete) thread", "error", err, "thread_id", t.ID)
			lastErr = err
		}
	}

	return lastErr
}
