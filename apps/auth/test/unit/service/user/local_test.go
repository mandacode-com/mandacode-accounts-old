package usersvc_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/internal/domain/dto"
	userdomain "mandacode.com/accounts/auth/internal/domain/service/user"
	usersvc "mandacode.com/accounts/auth/internal/service/user"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
)

type MockLocalUserService struct {
	mockRepo *mock_repodomain.MockLocalUserRepository
	svc      userdomain.LocalUserService
}

func (s *MockLocalUserService) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	s.mockRepo = mock_repodomain.NewMockLocalUserRepository(ctrl)
	s.svc = usersvc.NewLocalUserService(s.mockRepo)
}

func (s *MockLocalUserService) Teardown() {
	s.mockRepo = nil
	s.svc = nil
}

func TestLocalCreateUser(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	t.Run("Successful User Creation", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().CreateUser(id, email, gomock.Any(), nil, nil).Return(entUser, nil)

		user, err := mock.svc.CreateUser(id, email, password, nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("Error on User Creation", func(t *testing.T) {
		mock.mockRepo.EXPECT().CreateUser(id, email, gomock.Any(), nil, nil).Return(nil, errors.New("error creating user"))

		user, err := mock.svc.CreateUser(id, email, password, nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on Password Hashing", func(t *testing.T) {
		// Simulate an error in password hashing by passing an invalid password
		mock.mockRepo.EXPECT().CreateUser(id, email, gomock.Any(), nil, nil).Return(nil, bcrypt.ErrHashTooShort)

		user, err := mock.svc.CreateUser(id, email, "short", nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("User Already Exists", func(t *testing.T) {
		mock.mockRepo.EXPECT().CreateUser(id, email, gomock.Any(), nil, nil).Return(nil, errors.New("user already exists"))

		user, err := mock.svc.CreateUser(id, email, password, nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		invalidEmail := "invalid-email"
		mock.mockRepo.EXPECT().CreateUser(id, invalidEmail, gomock.Any(), nil, nil).Return(nil, errors.New("invalid email format"))

		user, err := mock.svc.CreateUser(id, invalidEmail, password, nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestLocalGetUserByEmail(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	t.Run("Successful User Retrieval", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(entUser, nil)

		user, err := mock.svc.GetUserByEmail(email)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(nil, errors.New("user not found"))

		user, err := mock.svc.GetUserByEmail(email)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on User Retrieval", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(nil, errors.New("error retrieving user"))

		user, err := mock.svc.GetUserByEmail(email)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestLocalGetUserByID(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	t.Run("Successful User Retrieval", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().GetUserByID(id).Return(entUser, nil)

		user, err := mock.svc.GetUserByID(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByID(id).Return(nil, errors.New("user not found"))

		user, err := mock.svc.GetUserByID(id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on User Retrieval", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByID(id).Return(nil, errors.New("error retrieving user"))

		user, err := mock.svc.GetUserByID(id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestLocalDeleteUser(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()

	t.Run("Successful User Deletion", func(t *testing.T) {
		mock.mockRepo.EXPECT().DeleteUser(id).Return(nil)

		deletedUser, err := mock.svc.DeleteUser(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if deletedUser.ID != id {
			t.Errorf("expected deleted user ID %v, got %v", id, deletedUser.ID)
		}
	})

	t.Run("Error on User Deletion", func(t *testing.T) {
		mock.mockRepo.EXPECT().DeleteUser(id).Return(errors.New("error deleting user"))

		deletedUser, err := mock.svc.DeleteUser(id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if deletedUser != nil {
			t.Errorf("expected deleted user to be nil, got %v", deletedUser)
		}
	})
}

func TestLocalUpdateEmail(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	newEmail := "new@test.com"

	t.Run("Successful Email Update", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      newEmail,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      newEmail,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().UpdateUser(id, &newEmail, nil, nil, nil).Return(entUser, nil)

		user, err := mock.svc.UpdateEmail(id, newEmail)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, &newEmail, nil, nil, nil).Return(nil, errors.New("user not found"))

		user, err := mock.svc.UpdateEmail(id, newEmail)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on Email Update", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, &newEmail, nil, nil, nil).Return(nil, errors.New("error updating email"))

		user, err := mock.svc.UpdateEmail(id, newEmail)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestLocalUpdatePassword(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	currentPassword := "currentPassword123"
	newPassword := "newPassword123"
	hashedCurrentPassword, err := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash current password: %v", err)
	}
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash new password: %v", err)
	}

	t.Run("Successful Password Update", func(t *testing.T) {
		currentEntUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedCurrentPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		newEntUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedNewPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		newDtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  newEntUser.CreatedAt,
			UpdatedAt:  newEntUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().GetUserByID(id).Return(currentEntUser, nil)
		mock.mockRepo.EXPECT().UpdateUser(id, nil, gomock.Any(), nil, nil).Return(newEntUser, nil)
		user, err := mock.svc.UpdatePassword(id, currentPassword, newPassword)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, newDtoUser) {
			t.Errorf("expected user %v, got %v", newDtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByID(id).Return(nil, errors.New("user not found"))

		user, err := mock.svc.UpdatePassword(id, currentPassword, newPassword)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Invalid Current Password", func(t *testing.T) {
		currentEntUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedCurrentPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByID(id).Return(currentEntUser, nil)

		user, err := mock.svc.UpdatePassword(id, "wrongCurrentPassword", newPassword)
		if err == nil {
			t.Fatal("expected error for invalid current password, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on Password Update", func(t *testing.T) {
		currentEntUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedCurrentPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByID(id).Return(currentEntUser, nil)
		mock.mockRepo.EXPECT().UpdateUser(id, nil, gomock.Any(), nil, nil).Return(nil, errors.New("error updating password"))

		user, err := mock.svc.UpdatePassword(id, currentPassword, newPassword)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestLocalUpdateActiveStatus(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	newStatus := false

	t.Run("Successful Active Status Update", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   newStatus,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   newStatus,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, &newStatus, nil).Return(entUser, nil)

		user, err := mock.svc.UpdateActiveStatus(id, newStatus)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, &newStatus, nil).Return(nil, errors.New("user not found"))

		user, err := mock.svc.UpdateActiveStatus(id, newStatus)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on Active Status Update", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, &newStatus, nil).Return(nil, errors.New("error updating active status"))

		user, err := mock.svc.UpdateActiveStatus(id, newStatus)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestLocalUpdateVerifiedStatus(t *testing.T) {
	mock := &MockLocalUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	newStatus := true

	t.Run("Successful Verified Status Update", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: newStatus,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: newStatus,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, nil, &newStatus).Return(entUser, nil)

		user, err := mock.svc.UpdateVerifiedStatus(id, newStatus)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, nil, &newStatus).Return(nil, errors.New("user not found"))

		user, err := mock.svc.UpdateVerifiedStatus(id, newStatus)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on Verified Status Update", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, nil, &newStatus).Return(nil, errors.New("error updating verified status"))

		user, err := mock.svc.UpdateVerifiedStatus(id, newStatus)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}
