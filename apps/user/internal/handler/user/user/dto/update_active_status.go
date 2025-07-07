package userhandlerdto

type UpdateActiveStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required"`
}
