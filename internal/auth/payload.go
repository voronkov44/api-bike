package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"email@example.com"`
	Password string `json:"password" validate:"required" example:"secret"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required" example:"Ivan"`
	Email    string `json:"email" validate:"required,email" example:"email@example.com"`
	Password string `json:"password" validate:"required" example:"secret"`
}

type RegisterResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}
