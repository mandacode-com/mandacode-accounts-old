package loginsvc

import (
	"context"
	"errors"

	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	logindomain "mandacode.com/accounts/auth/internal/domain/service/login"
)

type OAuthLoginService struct {
	oauthUser repodomain.OAuthUserRepository
}

func NewOAuthLoginService(oauthAuth repodomain.OAuthUserRepository) logindomain.OAuthLoginService {
	return &OAuthLoginService{
		oauthUser: oauthAuth,
	}
}

func (s *OAuthLoginService) LoginOAuthUser(ctx context.Context, provider oauthuser.Provider, providerID string) (*dto.OAuthUser, error) {
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
