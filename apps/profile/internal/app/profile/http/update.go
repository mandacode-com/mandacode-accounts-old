package httpuc

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/app/profile"
	"mandacode.com/accounts/profile/internal/domain/dto"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
)

type UpdateProfileUsecase struct {
	ProfileService svcdomain.ProfileService
}

func NewUpdateProfileUsecase(profileService svcdomain.ProfileService) profile.UpdateProfileUsecase {
	return &UpdateProfileUsecase{
		ProfileService: profileService,
	}
}

func (u *UpdateProfileUsecase) UpdateProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string) (*dto.Profile, error) {
	profile, err := u.ProfileService.UpdateProfile(userID, email, displayName, bio, avatarURL)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
