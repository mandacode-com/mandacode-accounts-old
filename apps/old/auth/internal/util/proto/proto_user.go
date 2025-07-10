package protoutil

import (
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"
	userdto "mandacode.com/accounts/auth/internal/app/user/dto"
	localuserv1 "github.com/mandacode-com/accounts-proto/auth/user/local/v1"
	oauthuserv1 "github.com/mandacode-com/accounts-proto/auth/user/oauth/v1"
)

func NewProtoLocalUser(user *userdto.LocalUser) (*localuserv1.LocalUser, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	return &localuserv1.LocalUser{
		UserId:     user.ID.String(),
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}

func NewProtoOAuthUser(user *userdto.OAuthUser) (*oauthuserv1.OAuthUser, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}
	entProvider, err := FromEntToProtoProvider(user.Provider)
	if err != nil {
		return nil, err
	}

	return &oauthuserv1.OAuthUser{
		UserId:     user.ID.String(),
		Provider:   entProvider,
		ProviderId: user.ProviderID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}
