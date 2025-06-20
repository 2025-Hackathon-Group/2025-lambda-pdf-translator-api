package response

type UserBasicResponse struct {
	ID           string                    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string                    `json:"name" example:"John Doe"`
	Email        string                    `json:"email" example:"john.doe@example.com"`
	Organisation OrganisationBasicResponse `json:"organisation"`
}
