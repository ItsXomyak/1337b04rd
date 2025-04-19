package s3

import (
	"1337b04rd/internal/app/common/logger"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"strings"
	"sync"
	"time"

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
	var (
		wg      sync.WaitGroup
		results sync.Map
		errs    = make(chan error, len(files))
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for fileName, file := range files {
		contentType := contentTypes[fileName]

		wg.Add(1)
		go func(name string, reader io.Reader, ctype string) {
			defer wg.Done()

			data, err := io.ReadAll(reader)
			if err != nil {
				errs <- fmt.Errorf("read error (%s): %w", name, err)
				return
			}

			buf := bytes.NewReader(data)

			ext := ""
			if exts, _ := mime.ExtensionsByType(ctype); len(exts) > 0 {
				ext = exts[0]
			}

			fileID, err := uuidHelper.NewUUID()
			if err != nil {
				errs <- fmt.Errorf("uuid error (%s): %w", name, err)
				return
			}
			uniqueName := fmt.Sprintf("%s%s", fileID.String(), ext)

			_, err = s.client.PutObject(ctx, s.bucket, uniqueName, buf, int64(len(data)), minio.PutObjectOptions{
				ContentType: ctype,
			})
			if err != nil {
				errs <- fmt.Errorf("upload error (%s): %w", name, err)
				return
			}

			url := s.GetImageURL(uniqueName)
			results.Store(name, url)
		}(fileName, file, contentType)
	}

	wg.Wait()
	close(errs)

	var allErrs []error
	for err := range errs {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) > 0 {
		for _, e := range allErrs {
			logger.Error("parallel upload error", "error", e)
		}
		return nil, fmt.Errorf("upload failed for %d file(s)", len(allErrs))
	}

	finalResults := make(map[string]string)
	results.Range(func(key, value any) bool {
		finalResults[key.(string)] = value.(string)
		return true
	})

	return finalResults, nil
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
	return nil
}
