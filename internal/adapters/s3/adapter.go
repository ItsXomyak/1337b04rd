package s3

import "io"

type Adapter struct {
	client *S3Client
}

func NewAdapter(client *S3Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) UploadImages(files map[string]io.Reader, contentTypes map[string]string) (map[string]string, error) {
	if len(files) == 1 {
		for name, file := range files {
			contentType := contentTypes[name]
			url, err := a.client.UploadImage(file, -1, contentType)
			if err != nil {
				return nil, err
			}
			return map[string]string{name: url}, nil
		}
	}
	return a.client.UploadImagesParallel(files, contentTypes)
}

func (a *Adapter) DeleteFile(fileName string) error {
	return a.client.DeleteFile(fileName)
}

func (a *Adapter) UploadImage(file io.Reader, size int64, contentType string) (string, error) {
	files := map[string]io.Reader{"image": file}
	contentTypes := map[string]string{"image": contentType}

	urls, err := a.client.UploadImagesParallel(files, contentTypes)
	if err != nil {
		return "", err
	}
	return urls["image"], nil
}
