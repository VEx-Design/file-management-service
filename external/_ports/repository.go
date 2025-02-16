package ports

import "io"

type FileRepository interface {
	UploadImage(bucketName, fileName string, file io.Reader) error
	UploadImageFromURL(bucketName, fileName, url string) error
	GetAllImages(bucketName string) ([]string, error)
	GetImage(bucketName, fileName string) (io.ReadCloser, error)
}
