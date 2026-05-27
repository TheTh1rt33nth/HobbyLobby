package dto

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
