package ports

import "io"

type S3Port interface {
	UploadImage(file io.Reader, size int64, contentType string) (string, error)
	DeleteFile(fileName string) error
}
