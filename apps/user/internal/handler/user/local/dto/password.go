package localhandlerdto

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UpdatePasswordResponse struct {
	Message string `json:"message,omitempty"`
}
