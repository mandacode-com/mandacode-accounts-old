package userhandlerdto

type UpdateVerifiedStatusRequest struct {
	IsVerified bool `json:"is_verified" binding:"required"`
}
