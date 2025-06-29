package grpcuc

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/app/profile"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
)

type DeleteProfileUsecase struct {
	ProfileService svcdomain.ProfileService
}

func NewDeleteProfileUsecase(profileService svcdomain.ProfileService) profile.DeleteProfileUsecase {
	return &DeleteProfileUsecase{
		ProfileService: profileService,
	}
}

func (u *DeleteProfileUsecase) DeleteProfile(userID uuid.UUID) error {
	err := u.ProfileService.DeleteProfile(userID)
	if err != nil {
		return err
	}

	return nil
}
