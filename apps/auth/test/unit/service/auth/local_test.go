package authsvc_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/internal/domain/dto"
	authdomain "mandacode.com/accounts/auth/internal/domain/service/auth"
	authsvc "mandacode.com/accounts/auth/internal/service/auth"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
)

type MockLocalAuthService struct {
	mockRepo *mock_repodomain.MockLocalUserRepository
	svc      authdomain.LocalAuthService
}

func (s *MockLocalAuthService) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	s.mockRepo = mock_repodomain.NewMockLocalUserRepository(ctrl)
	s.svc = authsvc.NewLocalAuthService(s.mockRepo)
}

func (s *MockLocalAuthService) Teardown() {
	s.mockRepo = nil
	s.svc = nil
}

func TestLoginLocalUser(t *testing.T) {
	mock := &MockLocalAuthService{}
	mock.Setup(t)
	defer mock.Teardown()

	ctx := context.Background()
	id := uuid.New()
	email := "test@test.com"
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	t.Run("Successful Login", func(t *testing.T) {
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

		user, err := mock.svc.LoginLocalUser(ctx, email, password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(nil, errors.New("user not found"))

		user, err := mock.svc.LoginLocalUser(ctx, email, password)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(entUser, nil)

		user, err := mock.svc.LoginLocalUser(ctx, email, "wrongpassword")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("User Not Active", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   false,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(entUser, nil)

		user, err := mock.svc.LoginLocalUser(ctx, email, password)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("User Not Verified", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(&ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: false,
		}, nil)

		user, err := mock.svc.LoginLocalUser(ctx, email, password)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}
