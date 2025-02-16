package logic

import "io"

type FileService interface {
	UploadImage(bucketName, fileName string, file io.Reader) error
	UploadImageFromURL(bucketName, fileName, url string) error
	GetAllImages(bucketName string) ([]string, error)
	GetImage(bucketName, fileName string) (io.ReadCloser, error)
}
