package response

import (
	"time"

	"mandacode.com/accounts/profile/internal/domain/dto"
)

type ProfileResponse struct {
	UserID      string    `json:"user_id"`
	Email       *string   `json:"email,omitempty"`
	DisplayName *string   `json:"display_name,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProfileResponse(p *dto.Profile) *ProfileResponse {
	return &ProfileResponse{
		UserID:      p.UserID.String(),
		Email:       p.Email,
		DisplayName: p.DisplayName,
		Bio:         p.Bio,
		AvatarURL:   p.AvatarURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
