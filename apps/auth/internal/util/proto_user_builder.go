package util

import (
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/auth/internal/domain/dto"
	localuserv1 "mandacode.com/accounts/proto/auth/user/local/v1"
	oauthuserv1 "mandacode.com/accounts/proto/auth/user/oauth/v1"
)

func BuildProtoLocalUser(user *dto.LocalUser) (*localuserv1.LocalUser, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	protoLocalUser := &localuserv1.LocalUser{
		UserId:     user.ID.String(),
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}

	err := protoLocalUser.ValidateAll()
	if err != nil {
		return nil, errors.New("failed to validate local user proto: " + err.Error())
	}

	return protoLocalUser, nil
}

func BuildProtoOAuthUser(user *dto.OAuthUser) (*oauthuserv1.OAuthUser, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}
	providerProto, err := FromProviderToProto(user.Provider)
	if err != nil {
		return nil, err
	}

	protoOAuthUser := &oauthuserv1.OAuthUser{
		UserId:     user.ID.String(),
		Email:      user.Email,
		Provider:   providerProto,
		ProviderId: user.ProviderID,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}

	err = protoOAuthUser.ValidateAll()
	if err != nil {
		return nil, errors.New("failed to validate OAuth user proto: " + err.Error())
	}

	return protoOAuthUser, nil
}
