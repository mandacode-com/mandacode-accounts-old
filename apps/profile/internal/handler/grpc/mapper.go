package grpchandler

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/profile/internal/domain/dto"
	profilev1 "mandacode.com/accounts/proto/profile/v1"
)

func ToProtoProfile(p *dto.Profile) (*profilev1.Profile, error) {
	protoProfile := &profilev1.Profile{
		UserId:      p.UserID.String(),
		Email:       p.Email,
		DisplayName: p.DisplayName,
		Bio:         p.Bio,
		AvatarUrl:   p.AvatarURL,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
	err := protoProfile.ValidateAll()
	if err != nil {
		return nil, err
	}
	return protoProfile, nil
}
