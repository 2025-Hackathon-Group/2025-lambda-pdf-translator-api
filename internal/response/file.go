package response

import "time"

type FileUploadResponse struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FileName  string    `json:"file_name" example:"example.pdf"`
	Bucket    string    `json:"bucket" example:"my-bucket"`
	Path      string    `json:"path" example:"path/to/file"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
}
