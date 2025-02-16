package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type fileRepositoryMinio struct {
	client *minio.Client
}

func NewFileRepositoryMinio(client *minio.Client) *fileRepositoryMinio {
	return &fileRepositoryMinio{
		client: client,
	}
}

func (r *fileRepositoryMinio) UploadImage(bucketName, fileName string, file io.Reader) error {
	ctx := context.Background()

	// Ensure the bucket exists
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Failed to check bucket: %v", err)
		return err
	}
	if !exists {
		err = r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Failed to create bucket: %v", err)
			return err
		}
	}

	// Upload file
	info, err := r.client.PutObject(ctx, bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return err
	}

	log.Printf("Successfully uploaded file: %s (size: %d)", info.Key, info.Size)
	return nil
}

func (r *fileRepositoryMinio) UploadImageFromURL(bucketName, fileName, url string) error {
	// Get the image from the URL
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch image from URL: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Ensure the bucket exists
	ctx := context.Background()
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Failed to check bucket: %v", err)
		return err
	}
	if !exists {
		err = r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Failed to create bucket: %v", err)
			return err
		}
	}

	// Upload the image from the URL to MinIO
	info, err := r.client.PutObject(ctx, bucketName, fileName, resp.Body, -1, minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		log.Printf("Failed to upload image from URL: %v", err)
		return err
	}

	log.Printf("Successfully uploaded image from URL: %s (size: %d)", info.Key, info.Size)
	return nil
}

func (r *fileRepositoryMinio) GetAllImages(bucketName string) ([]string, error) {
	ctx := context.Background()

	// Check if bucket exists
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Failed to check bucket: %v", err)
		return nil, err
	}
	if !exists {
		log.Printf("Bucket does not exist: %s", bucketName)
		return nil, fmt.Errorf("bucket does not exist")
	}

	var imageList []string

	// List all objects in the bucket
	objectCh := r.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})
	for object := range objectCh {
		if object.Err != nil {
			log.Printf("Error retrieving object: %v", object.Err)
			return nil, object.Err
		}
		imageList = append(imageList, object.Key)
	}

	return imageList, nil
}

func (r *fileRepositoryMinio) GetImage(bucketName, fileName string) (io.ReadCloser, error) {
	ctx := context.Background()

	// Check if bucket exists
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Failed to check bucket: %v", err)
		return nil, err
	}
	if !exists {
		log.Printf("Bucket does not exist: %s", bucketName)
		return nil, fmt.Errorf("bucket does not exist")
	}

	// Get the object from MinIO
	object, err := r.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Failed to retrieve file: %v", err)
		return nil, err
	}

	return object, nil
}
