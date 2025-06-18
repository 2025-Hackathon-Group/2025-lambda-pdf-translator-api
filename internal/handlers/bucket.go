package handler

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/repository"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type BucketHandler struct {
	bucketRepository repository.MinioBucketRepository
	fileRepository   repository.FileRepository
}

func NewBucketHandler(bucketRepository repository.MinioBucketRepository, fileRepository repository.FileRepository) *BucketHandler {
	return &BucketHandler{bucketRepository: bucketRepository, fileRepository: fileRepository}
}

// UploadFile uploads a file to the bucket
// @Summary Upload a file to the bucket
// @Description Uploads a file to the bucket
// @Tags buckets
// @ID upload-file
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param bucket_name formData string true "Bucket name"
// @Param file_path formData string true "File path in bucket"
// @Param file formData file true "File to upload"
// @Success 200 {object} object{message=string,file=response.FileBasicResponse} "File uploaded successfully"
// @Failure 400 {object} object{error=string} "Invalid input"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /bucket/upload [post]
func (h *BucketHandler) UploadFile(c *gin.Context) {
	var input struct {
		File       *multipart.FileHeader `form:"file" binding:"required" json:"name" example:"John Doe"`
		BucketName string                `form:"bucket_name" binding:"required"`
		FilePath   string                `form:"file_path" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := models.FileUpload{
		FileName: input.File.Filename,
		S3Bucket: input.BucketName,
		Path:     input.FilePath,
	}

	if err := h.fileRepository.CreateFile(&file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.bucketRepository.UploadFileFromMultpart(input.BucketName, input.FilePath, input.File); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file": file})
}

// GetFileByID returns file metadata by file ID
// @Summary Get file metadata by file ID
// @Description Returns file metadata for a given file ID
// @Tags buckets
// @ID get-file-by-id
// @Accept json
// @Produce json
// @Security Bearer
// @Param file_id path string true "File ID" example("1234")
// @Success 200 {object} object{file=response.FileBasicResponse} "File metadata"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /bucket/file/{file_id} [get]
func (h *BucketHandler) GetFileByID(c *gin.Context) {
	fileId := c.Param("file_id")

	// Get file from bucket
	file, err := h.fileRepository.GetFile(fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

// GetFileByPath returns a file by its path and bucket name
// @Summary Download file by path and bucket name
// @Description Downloads a file from the bucket using its path and bucket name
// @Tags buckets
// @ID get-file-by-path
// @Accept json
// @Produce octet-stream
// @Security Bearer
// @Param bucket_name query string true "Bucket name" example("my-bucket")
// @Param file_path query string true "File path in bucket" example("docs/example.pdf")
// @Success 200 {file} file "File download"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /bucket/file [get]
func (h *BucketHandler) GetFileByPath(c *gin.Context) {
	// Query params
	filePath := c.Query("file_path")
	bucketName := c.Query("bucket_name")

	// Get file from bucket
	file, err := h.bucketRepository.GetObjectFromPath(bucketName, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Copy file to tmp
	tmpFile, err := os.CreateTemp("", "file-*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tmpFile.Close()
	io.Copy(tmpFile, file)

	// Serve file
	c.File(tmpFile.Name())
}

// GetObjectFromID downloads a file by its ID
// @Summary Download file by file ID
// @Description Downloads a file from the bucket using its file ID
// @Tags buckets
// @ID get-object-from-id
// @Accept json
// @Produce octet-stream
// @Security Bearer
// @Param file_id path string true "File ID" example("1234")
// @Success 200 {file} file "File download"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /bucket/object/{file_id} [get]
func (h *BucketHandler) GetObjectFromID(c *gin.Context) {
	fileId := c.Param("file_id")

	// Get file from bucket
	file, err := h.fileRepository.GetFile(fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get file from bucket
	s3Object, err := h.bucketRepository.GetObjectFromPath(file.S3Bucket, file.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Copy file to tmp
	tmpFile, err := os.CreateTemp("", "file-*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tmpFile.Close()
	io.Copy(tmpFile, s3Object)

	// Serve file
	c.File(tmpFile.Name())
}
