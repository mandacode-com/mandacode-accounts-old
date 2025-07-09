package profile

import (
	"context"

	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/profile/v1"
	profiledomain "mandacode.com/accounts/user/internal/domain/port/profile"
)

type ProfileService struct {
	client profilev1.ProfileServiceClient
}

// CreateProfile implements profiledomain.ProfileService.
func (p *ProfileService) CreateProfile(ctx context.Context, userID uuid.UUID) (*profilev1.CreateProfileResponse, error) {
	resp, err := p.client.CreateProfile(ctx, &profilev1.CreateProfileRequest{UserId: userID.String()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteProfile implements profiledomain.ProfileService.
func (p *ProfileService) DeleteProfile(ctx context.Context, userID uuid.UUID) (*profilev1.DeleteProfileResponse, error) {
	resp, err := p.client.DeleteProfile(ctx, &profilev1.DeleteProfileRequest{UserId: userID.String()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewProfileService(client profilev1.ProfileServiceClient) profiledomain.ProfileService {
	return &ProfileService{client: client}
}
