package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
)

type S3Client struct {
	endpoint      string
	accessKey     string
	secretKey     string
	postBucket    string
	commentBucket string
	client        *http.Client
	logger        *slog.Logger
}

func NewS3Client(endpoint, accessKey, secretKey, postBucket, commentBucket string, logger *slog.Logger) (*S3Client, error) {
	// Удаляем схему из endpoint, если она есть
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, bucket := range []string{postBucket, commentBucket} {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := fmt.Sprintf("http://%s/%s", endpoint, bucket)
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
		if err != nil {
			logger.Error("failed to create bucket check request", "bucket", bucket, "error", err)
			return nil, err
		}
		req.SetBasicAuth(accessKey, secretKey)

		resp, err := client.Do(req)
		if err != nil {
			logger.Error("failed to check bucket", "bucket", bucket, "error", err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			logger.Warn("bucket does not exist, assuming it will be created manually", "bucket", bucket)
		} else if resp.StatusCode != http.StatusOK {
			logger.Error("unexpected response while checking bucket", "bucket", bucket, "status", resp.Status)
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}
	}

	return &S3Client{
		endpoint:      endpoint,
		accessKey:     accessKey,
		secretKey:     secretKey,
		postBucket:    postBucket,
		commentBucket: commentBucket,
		client:        client,
		logger:        logger,
	}, nil
}

func (s *S3Client) UploadPostImage(file io.Reader, fileName, contentType string) (string, error) {
	return s.uploadImage(file, fileName, contentType, s.postBucket)
}

func (s *S3Client) UploadCommentImage(file io.Reader, fileName, contentType string) (string, error) {
	return s.uploadImage(file, fileName, contentType, s.commentBucket)
}

func (s *S3Client) uploadImage(file io.Reader, fileName, contentType, bucket string) (string, error) {
	fileID, err := uuidHelper.NewUUID()
	if err != nil {
		s.logger.Error("failed to generate file ID", "error", err)
		return "", err
	}
	ext := filepath.Ext(fileName)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			s.logger.Error("unsupported content type", "contentType", contentType)
			return "", fmt.Errorf("unsupported content type: %s", contentType)
		}
	}
	newFileName := fmt.Sprintf("%s%s", fileID.String(), ext)

	data, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("failed to read file", "error", err)
		return "", err
	}

	if len(data) > 10<<20 {
		s.logger.Error("file too large", "size", len(data))
		return "", fmt.Errorf("file too large: %d bytes", len(data))
	}

	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, bucket, newFileName)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		s.logger.Error("failed to create upload request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.SetBasicAuth(s.accessKey, s.secretKey)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("failed to upload image", "file", newFileName, "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected response from S3", "file", newFileName, "status", resp.Status)
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	s.logger.Info("successfully uploaded image", "url", url)
	return url, nil
}

func (s *S3Client) GetImageURL(fileName, bucket string) string {
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, bucket, fileName)
}