package usersvc

import (
	"errors"
	"time"

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

func (s *OAuthUserService) CreateUser(
	userID uuid.UUID, provider oauthuser.Provider, providerID string, email string, isActive *bool, isVerified *bool) (*dto.OAuthUser, error) {
	user, error := s.oauthUser.CreateUser(userID, provider, providerID, email, isActive, isVerified)
	if error != nil {
		return nil, error
	}
	if user == nil {
		return nil, errors.New("failed to create user")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}

func (s *OAuthUserService) GetUserByProvider(provider oauthuser.Provider, providerID string) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.GetUserByProvider(provider, providerID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}

func (s *OAuthUserService) GetUserByUserID(userID uuid.UUID, provider oauthuser.Provider) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.GetUserByUserID(userID, provider)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}

func (s *OAuthUserService) DeleteUser(userID uuid.UUID) (*dto.OAuthDeletedUser, error) {
	err := s.oauthUser.DeleteUser(userID)
	if err != nil {
		return nil, err
	}

	return &dto.OAuthDeletedUser{
		ID:        userID,
		DeletedAt: time.Now(),
	}, nil

}

func (s *OAuthUserService) DeleteUserByProvider(userID uuid.UUID, provider oauthuser.Provider) (*dto.OAuthDeletedUser, error) {
	err := s.oauthUser.DeleteUserByProvider(userID, provider)
	if err != nil {
		return nil, err
	}

	return &dto.OAuthDeletedUser{
		ID:        userID,
		Provider:  &provider,
		DeletedAt: time.Now(),
	}, nil
}

func (s *OAuthUserService) UpdateUserBase(
	userID uuid.UUID, provider oauthuser.Provider, providerID string, email string, isVerified bool) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.UpdateUser(userID, provider, &providerID, &email, nil, &isVerified)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}

func (s *OAuthUserService) UpdateActiveStatus(userID uuid.UUID, provider oauthuser.Provider, isActive bool) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.UpdateUser(userID, provider, nil, nil, &isActive, nil)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}

func (s *OAuthUserService) UpdateVerifiedStatus(userID uuid.UUID, provider oauthuser.Provider, isVerified bool) (*dto.OAuthUser, error) {
	user, err := s.oauthUser.UpdateUser(userID, provider, nil, nil, nil, &isVerified)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	oauthUser := dto.NewOAuthUserFromEntity(user)
	return oauthUser, nil
}
