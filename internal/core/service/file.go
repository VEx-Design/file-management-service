// image_service.go
package service

import (
	"io"
	ports "type-management-service/external/_ports"
	"type-management-service/internal/core/logic"
)

type FileService struct {
	uploader ports.FileRepository
}

func NewImageService(uploader ports.FileRepository) logic.FileService {
	return &FileService{uploader: uploader}
}

func (s *FileService) UploadImage(bucketName, fileName string, file io.Reader) error {
	return s.uploader.UploadImage(bucketName, fileName, file)
}

func (s *FileService) UploadImageFromURL(bucketName, fileName, url string) error {
	return s.uploader.UploadImageFromURL(bucketName, fileName, url)
}

func (s *FileService) GetAllImages(bucketName string) ([]string, error) {
	return s.uploader.GetAllImages(bucketName)
}

func (s *FileService) GetImage(bucketName, fileName string) (io.ReadCloser, error) {
	return s.uploader.GetImage(bucketName, fileName)
}
