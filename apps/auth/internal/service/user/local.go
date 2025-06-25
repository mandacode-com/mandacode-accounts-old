package usersvc

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/internal/domain/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	userdomain "mandacode.com/accounts/auth/internal/domain/service/user"
)

type LocalUserService struct {
	localUser repodomain.LocalUserRepository
}

// NewLocalUserService creates a new instance of localUserService
func NewLocalUserService(localUser repodomain.LocalUserRepository) userdomain.LocalUserService {
	return &LocalUserService{
		localUser: localUser,
	}
}

func (s *LocalUserService) CreateLocalUser(userID uuid.UUID, email string, password string, isActive *bool, isVerified *bool) (*dto.LocalUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.localUser.CreateLocalUser(userID, email, string(hashedPassword), isActive, isVerified)
	if err != nil {
		return nil, err
	}

	return &dto.LocalUser{
		ID:         user.ID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}, nil
}

func (s *LocalUserService) DeleteLocalUser(userID uuid.UUID) error {
	return s.localUser.DeleteLocalUser(userID)
}

func (s *LocalUserService) UpdateLocalUser(userID uuid.UUID, email *string, password *string, isActive *bool, isVerified *bool) (*dto.LocalUser, error) {
	var hashedPassword *string
	if password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPassword = new(string)
		*hashedPassword = string(hashed)
	}
	user, err := s.localUser.UpdateLocalUser(userID, email, hashedPassword, isActive, isVerified)
	if err != nil {
		return nil, err
	}
	return &dto.LocalUser{
		ID:         userID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}, nil
}
