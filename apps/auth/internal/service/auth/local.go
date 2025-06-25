package authsvc

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/internal/domain/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	authdomain "mandacode.com/accounts/auth/internal/domain/service/auth"
)

type LocalAuthService struct {
	localUser repodomain.LocalUserRepository
}

func NewLocalAuthService(localAuth repodomain.LocalUserRepository) authdomain.LocalAuthService {
	return &LocalAuthService{
		localUser: localAuth,
	}
}

func (s *LocalAuthService) LoginLocalUser(ctx context.Context, email string, password string) (*dto.LocalUser, error) {
	user, err := s.localUser.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}
	if !user.IsVerified {
		return nil, errors.New("user is not verified")
	}

	return &dto.LocalUser{
		ID:         user.ID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}, nil
}
