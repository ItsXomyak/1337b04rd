package s3

// Я здесь хуй знает че творится, в тупую скатал, потом разберемся

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"sync"
	"time"

	"1337b04rd/internal/app/common/logger"
	uuidHelper "1337b04rd/internal/app/common/utils"
)

type S3Client struct {
	endpoint  string // Например, "localhost:9000"
	accessKey string
	secretKey string
	bucket    string
	client    *http.Client
}

func NewS3Client(endpoint, accessKey, secretKey, bucket string) (*S3Client, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &S3Client{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		client:    client,
	}, nil
}

func (s *S3Client) UploadImagesParallel(files map[string]io.Reader) (map[string]string, error) {
	var wg sync.WaitGroup
	results := make(map[string]string)
	errors := make(chan error, len(files))

	for fileName, file := range files {
		wg.Add(1)
		go func(fileName string, file io.Reader) {
			defer wg.Done()
			url, err := s.UploadImage(file, 0, "")
			if err != nil {
				errors <- err
				return
			}
			results[fileName] = url
		}(fileName, file)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return nil, <-errors
	}

	return results, nil
}

func (s *S3Client) UploadImage(file io.Reader, size int64, contentType string) (string, error) {
	fileID, err := uuidHelper.NewUUID()
	if err != nil {
		logger.Error("failed to generate UUID", "error", err)
		return "", err
	}

	// Получаем расширение из MIME-типа (например, ".jpg" из "image/jpeg")
	exts, _ := mime.ExtensionsByType(contentType)
	ext := ""
	if len(exts) > 0 {
		ext = exts[0]
	}

	fileName := fmt.Sprintf("%s%s", fileID.String(), ext)

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read file", "error", err)
		return "", err
	}

	logger.Debug("s3 upload debug", "contentType", contentType, "bucket", s.bucket, "fileName", fileName)

	url := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, fileName)

	logger.Debug("s3 upload debug",
		"contentType", contentType,
		"len(data)", len(data),
		"fileName", fileName,
		"bucket", s.bucket,
		"endpoint", s.endpoint,
	)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		logger.Error("failed to create upload request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(s.accessKey, s.secretKey)

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("failed to upload image", "file", fileName, "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("unexpected response from S3", "file", fileName, "status", resp.Status)
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	logger.Info("successfully uploaded image", "url", url)
	return url, nil
}

func (s *S3Client) GetImageURL(fileName string) string {
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, fileName)
}

func (s *S3Client) DeleteFile(fileName string) error {
	url := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, fileName)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		logger.Error("failed to create delete request", "error", err)
		return err
	}

	req.SetBasicAuth(s.accessKey, s.secretKey)

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("failed to send delete request", "file", fileName, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		logger.Error("unexpected response from S3 while deleting file", "file", fileName, "status", resp.Status)
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	logger.Info("successfully deleted file", "file", fileName)
	return nil
}
