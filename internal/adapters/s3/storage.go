package s3

import (
	"1337b04rd/internal/app/common/logger"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	uuidHelper "1337b04rd/internal/app/common/utils"
)

type S3Client struct {
	endpoint string
	bucket   string
}

func NewS3Client(endpoint, bucket string) *S3Client {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	return &S3Client{
		endpoint: endpoint,
		bucket:   bucket,
	}
}

func (s *S3Client) UploadImage(file io.Reader, _ int64, contentType string) (string, error) {
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

	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s\n%s", resp.Status, string(body))
	}

	logger.Info("image uploaded (raw)", "url", url)
	return url, nil
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

			url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, uniqueName)
			req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(data))
			if err != nil {
				errs <- fmt.Errorf("request error (%s): %w", name, err)
				return
			}
			req.Header.Set("Content-Type", ctype)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errs <- fmt.Errorf("http error (%s): %w", name, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				body, _ := io.ReadAll(resp.Body)
				errs <- fmt.Errorf("upload failed (%s): %s\n%s", name, resp.Status, string(body))
				return
			}

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
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)
}

func (s *S3Client) DeleteFile(fileName string) error {
	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, fileName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s\n%s", resp.Status, string(body))
	}

	return nil
}
