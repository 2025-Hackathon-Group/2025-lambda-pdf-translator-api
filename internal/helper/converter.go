package helper

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/response"
)

// User

func ToUserResponse(user models.User) response.UserBasicResponse {

	return response.UserBasicResponse{
		ID:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		Organisation: ToOrganisationResponse(user.Organisation),
	}
}

func ToUsersResponse(users []models.User) []response.UserBasicResponse {
	usersResponse := make([]response.UserBasicResponse, len(users))
	for i, user := range users {
		usersResponse[i] = ToUserResponse(user)
	}
	return usersResponse
}

// Organisation

func ToOrganisationResponse(org models.Organisation) response.OrganisationBasicResponse {
	return response.OrganisationBasicResponse{
		ID:   org.ID.String(),
		Name: org.Name,
	}
}

func ToOrganisationsResponse(orgs []models.Organisation) []response.OrganisationBasicResponse {
	orgsResponse := make([]response.OrganisationBasicResponse, len(orgs))
	for i, org := range orgs {
		orgsResponse[i] = ToOrganisationResponse(org)
	}
	return orgsResponse
}

// File

func ToFileResponse(file models.FileUpload) response.FileUploadResponse {
	return response.FileUploadResponse{
		ID:        file.ID.String(),
		FileName:  file.FileName,
		Bucket:    file.S3Bucket,
		Path:      file.Path,
		CreatedAt: file.CreatedAt,
	}
}

func ToFilesResponse(files []models.FileUpload) []response.FileUploadResponse {
	filesResponse := make([]response.FileUploadResponse, len(files))
	for i, file := range files {
		filesResponse[i] = ToFileResponse(file)
	}
	return filesResponse
}
