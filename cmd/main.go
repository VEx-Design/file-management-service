package main

import (
	"fmt"
	"log"

	handler "type-management-service/external/handler/router"
	repository "type-management-service/external/repository/adaptors/minio/controller"
	"type-management-service/internal/core/service"
	"type-management-service/pkg/minio"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MinIO
	minioClient := minio.ConnectToMinio()
	if minioClient == nil {
		log.Fatal("Failed to connect to MinIO")
	}
	fmt.Println("Successfully connected to MinIO")

	// Initialize dependencies
	client := minioClient.GetClient()
	fileRep := repository.NewFileRepositoryMinio(client)
	fileService := service.NewImageService(fileRep)
	uploadHandler := handler.NewUploadHandler(fileService)

	// Setup router
	r := gin.Default()
	r.POST("/upload", uploadHandler.UploadImage)
	r.POST("/upload-from-url", uploadHandler.UploadImageFromURL)
	// r.GET("/getImage", uploadHandler.GetImage)
	r.GET("/getbucket/:bucketName", uploadHandler.GetAllImages)
	r.GET("/getimg/:bucketName/:fileName", uploadHandler.GetImage)
	// Start the server
	r.Run(":8088")
}
