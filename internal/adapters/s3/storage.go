package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"strings"
	"sync"
	"time"

	"1337b04rd/internal/app/common/logger"
	uuidHelper "1337b04rd/internal/app/common/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	client *minio.Client
	bucket string
}

func NewS3Client(endpoint, accessKey, secretKey, bucket string) (*S3Client, error) {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("failed to initialize minio client", "error", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
	if errBucketExists != nil {
		logger.Error("failed to check bucket existence", "error", errBucketExists)
		return nil, errBucketExists
	}
	if !exists {
		logger.Warn("bucket does not exist", "bucket", bucket)
	}

	return &S3Client{
		client: minioClient,
		bucket: bucket,
	}, nil
}

func (s *S3Client) UploadImage(file io.Reader, size int64, contentType string) (string, error) {
	fileID, err := uuidHelper.NewUUID()
	if err != nil {
		logger.Error("failed to generate UUID", "error", err)
		return "", err
	}

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
	reader := bytes.NewReader(data)

	uploadInfo, err := s.client.PutObject(context.Background(), s.bucket, fileName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		logger.Error("failed to upload image to minio", "file", fileName, "error", err)
		return "", err
	}

	logger.Info("image uploaded to minio", "location", uploadInfo.Location)
	return s.GetImageURL(fileName), nil
}

func (s *S3Client) UploadImagesParallel(files map[string]io.Reader, contentTypes map[string]string) (map[string]string, error) {
	var wg sync.WaitGroup
	mu := sync.Mutex{}

	results := make(map[string]string)
	errs := make(chan error, len(files))

	for fileName, file := range files {
		contentType := contentTypes[fileName]

		wg.Add(1)
		go func(name string, reader io.Reader, ctype string) {
			defer wg.Done()

			data, err := io.ReadAll(reader)
			if err != nil {
				errs <- fmt.Errorf("failed to read file %s: %w", name, err)
				return
			}
			buf := bytes.NewReader(data)

			exts, _ := mime.ExtensionsByType(ctype)
			ext := ""
			if len(exts) > 0 {
				ext = exts[0]
			}

			fileID, err := uuidHelper.NewUUID()
			if err != nil {
				errs <- fmt.Errorf("failed to generate UUID for %s: %w", name, err)
				return
			}
			uniqueName := fmt.Sprintf("%s%s", fileID.String(), ext)

			logger.Debug("uploading to S3 (parallel)", "originalName", name, "generatedName", uniqueName, "contentType", ctype, "size", len(data))

			_, err = s.client.PutObject(context.Background(), s.bucket, uniqueName, buf, int64(len(data)), minio.PutObjectOptions{
				ContentType: ctype,
			})
			if err != nil {
				errs <- fmt.Errorf("failed to upload %s: %w", name, err)
				return
			}

			url := s.GetImageURL(uniqueName)
			mu.Lock()
			results[name] = url
			mu.Unlock()
		}(fileName, file, contentType)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return nil, <-errs
	}

	return results, nil
}

func (s *S3Client) GetImageURL(fileName string) string {
	return fmt.Sprintf("http://%s/%s/%s", s.client.EndpointURL().Host, s.bucket, fileName)
}

func (s *S3Client) DeleteFile(fileName string) error {
	err := s.client.RemoveObject(context.Background(), s.bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		logger.Error("failed to delete file", "file", fileName, "error", err)
		return err
	}
	logger.Info("successfully deleted file", "file", fileName)
	return nil
}
