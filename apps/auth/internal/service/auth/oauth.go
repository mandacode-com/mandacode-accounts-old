package authsvc

import (
	"context"
	"errors"

	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	authdomain "mandacode.com/accounts/auth/internal/domain/service/auth"
)

type OAuthAuthService struct {
	oauthUser repodomain.OAuthUserRepository
}

func NewOAuthAuthService(oauthAuth repodomain.OAuthUserRepository) authdomain.OAuthAuthService {
	return &OAuthAuthService{
		oauthUser: oauthAuth,
	}
}

func (s *OAuthAuthService) LoginOAuthUser(ctx context.Context, provider oauthuser.Provider, providerID string) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.GetUserByProvider(provider, providerID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("user is not active")
	}
	if !user.IsVerified {
		return nil, errors.New("user is not verified")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}
