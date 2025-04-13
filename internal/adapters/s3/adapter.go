package s3

import "io"

type Adapter struct {
	client *S3Client
}

func NewAdapter(client *S3Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) UploadImage(file io.Reader, size int64, contentType string) (string, error) {
	return a.client.UploadImage(file, size, contentType)
}

func (a *Adapter) DeleteFile(fileName string) error {
	return a.client.DeleteFile(fileName)
}
