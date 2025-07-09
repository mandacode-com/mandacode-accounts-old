package httphandlerdto

import "mandacode.com/accounts/profile/internal/domain/model"

type ProfileUpdateRequest struct {
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,min=3,max=20"`
	Bio         *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	AvatarURL   *string `json:"avatar_url,omitempty" binding:"omitempty,url,max=2048"`
}

type ProfileUpdateResponse struct {
	*model.Profile
}
