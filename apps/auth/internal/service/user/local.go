package usersvc

import (
	"errors"
	"time"

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

func (s *LocalUserService) CreateUser(userID uuid.UUID, email string, password string, isActive *bool, isVerified *bool) (*dto.LocalUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.localUser.CreateUser(userID, email, string(hashedPassword), isActive, isVerified)
	if err != nil {
		return nil, err
	}

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}

func (s *LocalUserService) GetUserByEmail(email string) (*dto.LocalUser, error) {
	user, err := s.localUser.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}

func (s *LocalUserService) GetUserByID(userID uuid.UUID) (*dto.LocalUser, error) {
	user, err := s.localUser.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}

func (s *LocalUserService) DeleteUser(userID uuid.UUID) (*dto.LocalDeletedUser, error) {
	err := s.localUser.DeleteUser(userID)
	if err != nil {
		return nil, err
	}

	localDeletedUser := dto.NewLocalDeletedUser(
		userID,
		time.Now(),
	)
	return localDeletedUser, nil
}

func (s *LocalUserService) UpdateEmail(userID uuid.UUID, newEmail string) (*dto.LocalUser, error) {
	user, err := s.localUser.UpdateUser(userID, &newEmail, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}

func (s *LocalUserService) UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) (*dto.LocalUser, error) {
	user, err := s.localUser.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // User not found
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return nil, err // Invalid current password
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedNewPasswordStr := string(hashedNewPassword)

	updatedUser, err := s.localUser.UpdateUser(userID, nil, &hashedNewPasswordStr, nil, nil)
	if err != nil {
		return nil, err
	}

	if updatedUser == nil {
		return nil, errors.New("user not found")
	}

	localUser := dto.NewLocalUserFromEntity(updatedUser)
	return localUser, nil
}

func (s *LocalUserService) UpdateActiveStatus(userID uuid.UUID, isActive bool) (*dto.LocalUser, error) {
	user, err := s.localUser.UpdateUser(userID, nil, nil, &isActive, nil)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}

func (s *LocalUserService) UpdateVerifiedStatus(userID uuid.UUID, isVerified bool) (*dto.LocalUser, error) {
	user, err := s.localUser.UpdateUser(userID, nil, nil, nil, &isVerified)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	localUser := dto.NewLocalUserFromEntity(user)
	return localUser, nil
}
