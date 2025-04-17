package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/ports"
	"1337b04rd/internal/domain/comment"
)

type CommentService struct {
	commentRepo ports.CommentPort
	threadRepo  ports.ThreadPort
	s3          ports.S3Port
	sessionRepo ports.SessionPort // Добавляем
}

func NewCommentService(
	commentRepo ports.CommentPort,
	threadRepo ports.ThreadPort,
	s3 ports.S3Port,
	sessionRepo ports.SessionPort, // Добавляем
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		threadRepo:  threadRepo,
		s3:          s3,
		sessionRepo: sessionRepo,
	}
}

func (s *CommentService) CreateComment(
	ctx context.Context,
	threadID utils.UUID,
	parentID *utils.UUID,
	content string,
	files map[string]io.Reader,
	contentTypes map[string]string,
	sessionID utils.UUID,
	displayName string,
	avatarURL string,
) (*comment.Comment, error) {
	if err := ctx.Err(); err != nil {
		logger.Warn("context canceled in CreateComment", "error", err)
		return nil, err
	}

	var imageURLs []string
	if len(files) > 0 {
		urls, err := s.s3.UploadImages(files, contentTypes)
		if err != nil {
			logger.Error("failed to upload comment images", "error", err)
			return nil, err
		}
		for _, url := range urls {
			updatedURL := strings.Replace(url, "http://minio:9000", "http://localhost:9000", 1)
			imageURLs = append(imageURLs, updatedURL)
		}
	}

	c, err := comment.NewComment(threadID, parentID, content, imageURLs, sessionID, displayName, avatarURL)
	if err != nil {
		logger.Error("cannot create new comment", "error", err)
		return nil, err
	}

	if err := s.commentRepo.CreateComment(ctx, c); err != nil {
		logger.Error("cannot save comment", "error", err)
		return nil, err
	}

	t, err := s.threadRepo.GetThreadByID(ctx, threadID)
	if err != nil {
		logger.Error("cannot fetch thread", "error", err)
		return nil, err
	}

	if err := s.threadRepo.UpdateThread(ctx, t); err != nil {
		logger.Error("cannot update thread", "error", err)
		return nil, err
	}

	logger.Info("comment created", "comment", c)
	return c, nil
}

func (s *CommentService) PrepareFilesFromMultipart(form *multipart.Form) (map[string]io.Reader, map[string]string, error) {
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

	for _, c := range comments {
		for i, url := range c.ImageURLs {
			c.ImageURLs[i] = strings.Replace(url, "http://minio:9000", "http://localhost:9000", 1)
		}

		if c.DisplayName == "" || c.AvatarURL == "" {
			session, err := s.sessionRepo.GetSessionByID(c.SessionID.String())
			if err != nil {
				logger.Warn("failed to get session for comment", "session_id", c.SessionID, "error", err)
				c.DisplayName = "Anonymous"
				c.AvatarURL = "http://localhost:9000/default-avatar.png"
			} else {
				c.DisplayName = session.DisplayName
				c.AvatarURL = session.AvatarURL
			}
		}
	}

	logger.Info("comments retrieved", "thread_id", threadID, "count", len(comments))
	return comments, nil
}
