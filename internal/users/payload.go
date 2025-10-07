package users

type UserResponse struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Name      string `json:"name" example:"John Doe"`
	CreatedAt string `json:"created_at" example:"2025-10-07T12:00:00Z"`
}
