package repository

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go"
)

type MinioBucketRepository struct {
	client    *minio.Client
	endpoint  string
	accessKey string
	secretKey string
}

// NewMinioBucketRepository creates a new MinioBucketRepository
func NewMinioBucketRepository(endpoint, accessKey, secretKey string) MinioBucketRepository {
	minioClient, err := minio.New(endpoint, accessKey, secretKey, false)
	if err != nil {
		log.Fatal(err)
	}

	return MinioBucketRepository{
		client:    minioClient,
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// UploadFileFromMultpart uploads a file to a bucket from a multipart file header
func (r *MinioBucketRepository) UploadFileFromMultpart(bucketName string, path string, file *multipart.FileHeader) (string, error) {
	tmpFile, err := os.CreateTemp("", "file-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	contents, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	if _, err := io.Copy(tmpFile, contents); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	if _, err := r.client.FPutObject(bucketName, path, tmpFile.Name(), minio.PutObjectOptions{}); err != nil {
		log.Println("failed to upload file: %w", err)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return path, nil
}

// GetObjectFromPath gets a file from a bucket
func (r *MinioBucketRepository) GetObjectFromPath(bucketName, objectName string) (io.Reader, error) {
	object, err := r.client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return object, nil
}

// TestConnection tests the connection to the bucket
func (r *MinioBucketRepository) TestConnection() error {
	_, err := r.client.ListBuckets()
	if err != nil {
		return fmt.Errorf("failed to test connection: %w", err)
	}

	return nil
}

// PrintInfo prints the information of the bucket
func (r *MinioBucketRepository) PrintInfo() {
	log.Println("Endpoint: ", r.endpoint)
	log.Println("AccessKey: ", r.accessKey)
	log.Println("SecretKey: ", r.secretKey)
}
