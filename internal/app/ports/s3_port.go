package ports

import "io"

type S3Port interface {
	UploadImages(files map[string]io.Reader, contentTypes map[string]string) (map[string]string, error)
	DeleteFile(fileName string) error
}
