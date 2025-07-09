package localhandlerdto

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}
