package request

type ProfileUpdateRequest struct {
	UserID      string  `json:"user_id" binding:"required,uuid"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,min=3,max=20"`
	Bio         *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	AvatarURL   *string `json:"avatar_url,omitempty" binding:"omitempty,url,max=2048"`
}
