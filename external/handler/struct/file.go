package s

type UploadImageRequest struct {
	BucketName string `json:"bucketName"`
	FileName   string `json:"fileName"`
	URL        string `json:"url"`
}
