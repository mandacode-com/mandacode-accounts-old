package usersvc

import (
	"errors"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	userdomain "mandacode.com/accounts/auth/internal/domain/service/user"
)

type OAuthUserService struct {
	oauthUser repodomain.OAuthUserRepository
}

// NewOAuthUserService creates a new instance of oauthService
func NewOAuthUserService(oauthUser repodomain.OAuthUserRepository) userdomain.OAuthUserService {
	return &OAuthUserService{
		oauthUser: oauthUser,
	}
}

func (s *OAuthUserService) CreateOAuthUser(
	userID uuid.UUID, provider oauthuser.Provider, providerID string, email string, isActive *bool, isVerified *bool) (*dto.OAuthUser, error) {
	user, error := s.oauthUser.CreateOAuthUser(userID, provider, providerID, email, isActive, isVerified)
	if error != nil {
		return nil, error
	}

	return &dto.OAuthUser{
		ID:         user.ID,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}, nil
}

func (s *OAuthUserService) DeleteOAuthUser(userID uuid.UUID) error {
	return s.oauthUser.DeleteOAuthUser(userID)
}

func (s *OAuthUserService) DeleteOAuthUserByProvider(userID uuid.UUID, provider oauthuser.Provider) error {
	return s.oauthUser.DeleteOAuthUserByProvider(userID, provider)
}

func (s *OAuthUserService) UpdateOAuthUser(
	userID uuid.UUID, provider *oauthuser.Provider, providerID *string, email *string, isActive *bool, isVerified *bool) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.UpdateOAuthUser(userID, provider, providerID, email, isActive, isVerified)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if user.Provider.String() == "" || user.ProviderID == "" || user.Email == "" {
		return nil, errors.New("user data cannot be empty")
	}
	return &dto.OAuthUser{
		ID:         userID,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}, nil
}
