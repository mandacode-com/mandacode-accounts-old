package grpcuc

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/app/profile"
	"mandacode.com/accounts/profile/internal/domain/dto"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
)

type InitializeProfileUsecase struct {
	ProfileService svcdomain.ProfileService
}

func NewInitializeProfileUsecase(profileService svcdomain.ProfileService) profile.InitializeProfileUsecase {
	return &InitializeProfileUsecase{
		ProfileService: profileService,
	}
}

func (u *InitializeProfileUsecase) InitializeProfile(userID uuid.UUID) (*dto.Profile, error) {
	profile, err := u.ProfileService.InitializeProfile(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
