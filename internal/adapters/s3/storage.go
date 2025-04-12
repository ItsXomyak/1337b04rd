package s3

// Я здесь хуй знает че творится, в тупую скатал, потом разберемся

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
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
	// Создаем HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Проверяем доступность MinIO и создаем бакет, если его нет
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("http://%s/%s", endpoint, bucket), nil)
	if err != nil {
		logger.Error("failed to create request", "error", err)
		return nil, err
	}
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to send request", "error", err)
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
	}, nil
}

// Функция для параллельной загрузки нескольких изображений
func (s *S3Client) UploadImagesParallel(files map[string]io.Reader) (map[string]string, error) {
	var wg sync.WaitGroup
	results := make(map[string]string)
	errors := make(chan error, len(files))

	// Загружаем каждый файл параллельно
	for fileName, file := range files {
		wg.Add(1)
		go func(fileName string, file io.Reader) {
			defer wg.Done()
			url, err := s.UploadImage(file, 0, "") // передаем загруженный файл и размер
			if err != nil {
				errors <- err
				return
			}
			results[fileName] = url
		}(fileName, file)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
	close(errors)

	// Проверяем, были ли ошибки
	if len(errors) > 0 {
		return nil, <-errors
	}

	return results, nil
}

func (s *S3Client) UploadImage(file io.Reader, size int64, contentType string) (string, error) {
	// Генерируем уникальное имя файла
	fileID, err := uuidHelper.NewUUID()
	if err != nil {
		logger.Error("failed to generate UUID", "error", err)
		return "", err
	}
	fileName := fmt.Sprintf("%s%s", fileID.String(), filepath.Ext(contentType)) // Например, "uuid.jpg"

	// Читаем файл в память
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read file", "error", err)
		return "", err
	}

	// Формируем URL для загрузки
	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)

	// Создаем PUT-запрос
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		logger.Error("failed to create upload request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.SetBasicAuth(s.accessKey, s.secretKey)

	// Выполняем запрос
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

// GetImageURL возвращает URL изображения по его имени (если нужно).
func (s *S3Client) GetImageURL(fileName string) string {
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)
}

// DeleteFile удаляет файл из MinIO по его имени.
func (s *S3Client) DeleteFile(fileName string) error {
	// Формируем URL для удаления файла
	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)

	// Создаем DELETE-запрос
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		logger.Error("failed to create delete request", "error", err)
		return err
	}

	// Добавляем аутентификацию
	req.SetBasicAuth(s.accessKey, s.secretKey)

	// Выполняем запрос
	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("failed to send delete request", "file", fileName, "error", err)
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusNoContent {
		logger.Error("unexpected response from S3 while deleting file", "file", fileName, "status", resp.Status)
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	logger.Info("successfully deleted file", "file", fileName)
	return nil
}
