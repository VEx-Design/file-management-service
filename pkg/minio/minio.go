package minio

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient struct holds the MinIO client instance
type MinioClient struct {
	client *minio.Client
}

// ConnectToMinio connects to the MinIO server
func ConnectToMinio() *MinioClient {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// MinIO connection details from environment variables
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Set true if using https
	})
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
		return nil
	}

	log.Println("Connected to MinIO")

	return &MinioClient{client: client}
}

func (minio *MinioClient) GetClient() *minio.Client {
	return minio.client
}
