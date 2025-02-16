package handler

import (
	"io"
	"log"
	"net/http" // Replace with the correct import path
	ports "type-management-service/external/_ports"
	s "type-management-service/external/handler/struct"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	repository ports.FileRepository
}

func NewUploadHandler(repository ports.FileRepository) *UploadHandler {
	return &UploadHandler{repository: repository}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Get the bucket name and file name from the request (can be query params or form data)
	bucketName := c.DefaultPostForm("bucketName", "default-bucket") // Default to "default-bucket" if not provided
	fileName := c.DefaultPostForm("fileName", "uploaded_image.jpg") // Default to "uploaded_image.jpg" if not provided

	// Get the file from the request
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Call the repository's UploadImage method to upload the file to MinIO
	err = h.repository.UploadImage(bucketName, fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// func (h *UploadHandler) UploadImageFromURL(c *gin.Context) {
// 	bucketName := c.PostForm("bucketName")
// 	fileName := c.PostForm("fileName")
// 	url := c.PostForm("url")

// 	if url == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
// 		return
// 	}

// 	err := h.repository.UploadImageFromURL(bucketName, fileName, url)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image from URL"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully from URL"})
// }

func (h *UploadHandler) UploadImageFromURL(c *gin.Context) {
	var req s.UploadImageRequest // Using the struct from dto package

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if URL is provided
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// Check if the URL exists (status code 200)
	resp, err := http.Get(req.URL)
	if err != nil {
		log.Printf("Error while trying to access URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check URL existence"})
		return
	}
	defer resp.Body.Close()

	// If status code is not 200, URL doesn't exist
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Call the service function to upload the image
	err = h.repository.UploadImageFromURL(req.BucketName, req.FileName, req.URL)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image from URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully from URL"})
}

func (h *UploadHandler) GetAllImages(c *gin.Context) {
	bucketName := c.Param("bucketName") // Get bucket name from URL parameter

	images, err := h.repository.GetAllImages(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *UploadHandler) GetImage(c *gin.Context) {
	bucketName := c.Param("bucketName")
	fileName := c.Param("fileName")

	// Get the image from MinIO
	image, err := h.repository.GetImage(bucketName, fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	defer image.Close()

	// Set the response headers
	c.Header("Content-Type", "image/jpeg") // Adjust based on image type
	c.Header("Content-Disposition", "inline; filename="+fileName)

	// Stream the image to the response
	_, err = io.Copy(c.Writer, image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image"})
		return
	}
}
