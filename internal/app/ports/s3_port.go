package ports

import "io"

type S3Port interface {
	UploadImages(files map[string]io.Reader, contentTypes map[string]string) (map[string]string, error)
	UploadImage(file io.Reader, size int64, contentType string) (string, error)
	DeleteFile(fileName string) error
}
