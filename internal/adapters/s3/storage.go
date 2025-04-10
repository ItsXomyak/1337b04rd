package s3

// Я здесь хуй знает че творится, в тупую скатал, потом разберемся

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
)

type S3Client struct {
	endpoint  string // Например, "localhost:9000"
	accessKey string
	secretKey string
	bucket    string
	client    *http.Client
	logger    *slog.Logger
}

func NewS3Client(endpoint, accessKey, secretKey, bucket string, logger *slog.Logger) (*S3Client, error) {
	// Создаем HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Проверяем доступность MinIO и создаем бакет, если его нет
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("http://%s/%s", endpoint, bucket), nil)
	if err != nil {
		logger.Error("failed to create bucket check request", "error", err)
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
		// Создаем бакет (MinIO не имеет прямого аналога MakeBucket через HTTP без XML, используем PUT)
		// Для простоты предполагаем, что бакет уже существует или создается вручную
		logger.Warn("bucket does not exist, assuming it will be created manually", "bucket", bucket)
	} else if resp.StatusCode != http.StatusOK {
		logger.Error("unexpected response while checking bucket", "status", resp.Status)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return &S3Client{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		client:    client,
		logger:    logger,
	}, nil
}

// UploadImage загружает изображение в S3 и возвращает его URL.
func (s *S3Client) UploadImage(file io.Reader, size int64, contentType string) (string, error) {
	// Генерируем уникальное имя файла
	fileID, err := uuidHelper.NewUUID()
	if err != nil {
		s.logger.Error("failed to generate file ID", "error", err)
		return "", err
	}
	fileName := fmt.Sprintf("%s%s", fileID.String(), filepath.Ext(contentType)) // Например, "uuid.jpg"

	// Читаем файл в память (для больших файлов можно использовать multipart upload, но начнем с простого)
	data, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("failed to read file", "error", err)
		return "", err
	}

	// Формируем URL для загрузки
	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)

	// Создаем PUT-запрос
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		s.logger.Error("failed to create upload request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.SetBasicAuth(s.accessKey, s.secretKey)

	// Выполняем запрос
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("failed to upload image", "file", fileName, "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected response from S3", "file", fileName, "status", resp.Status)
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	s.logger.Info("successfully uploaded image", "url", url)
	return url, nil
}

// GetImageURL возвращает URL изображения по его имени (если нужно).
func (s *S3Client) GetImageURL(fileName string) string {
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)
}