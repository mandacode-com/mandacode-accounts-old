package httpuc

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/app/profile"
	"mandacode.com/accounts/profile/internal/domain/dto"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
)

type GetProfileUsecase struct {
	ProfileService svcdomain.ProfileService
}

func NewGetProfileUsecase(profileService svcdomain.ProfileService) profile.GetProfileUsecase {
	return &GetProfileUsecase{
		ProfileService: profileService,
	}
}

func (u *GetProfileUsecase) GetProfile(userID uuid.UUID) (*dto.Profile, error) {
	profile, err := u.ProfileService.GetProfileByID(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
