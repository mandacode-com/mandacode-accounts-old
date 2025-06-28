package loginsvc

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/internal/domain/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	logindomain "mandacode.com/accounts/auth/internal/domain/service/login"
)

type LocalLoginService struct {
	localUser repodomain.LocalUserRepository
}

func NewLocalLoginService(localAuth repodomain.LocalUserRepository) logindomain.LocalLoginService {
	return &LocalLoginService{
		localUser: localAuth,
	}
}

func (s *LocalLoginService) LoginLocalUser(ctx context.Context, email string, password string) (*dto.LocalUser, error) {
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

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}
