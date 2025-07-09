package login_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	logindto "mandacode.com/accounts/auth/internal/app/login/dto"
	locallogin "mandacode.com/accounts/auth/internal/app/login/local"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
	mock_tokendomain "mandacode.com/accounts/auth/test/mock/domain/token"
)

type MockLocalLoginApp struct {
	mockTokenProvider *mock_tokendomain.MockTokenProvider
	mockRepository    *mock_repodomain.MockLocalUserRepository
	validate          *validator.Validate
	app               locallogin.LocalLoginApp
}

func (m *MockLocalLoginApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.mockTokenProvider = mock_tokendomain.NewMockTokenProvider(ctrl)
	m.mockRepository = mock_repodomain.NewMockLocalUserRepository(ctrl)
	m.validate = validator.New()
	m.app = locallogin.NewLocalLoginApp(m.mockTokenProvider, m.mockRepository)
}

func (m *MockLocalLoginApp) Teardown() {
	m.mockTokenProvider = nil
	m.mockRepository = nil
	m.validate = nil
	m.app = nil
}

func TestLocalAuthApp_LoginLocalUser(t *testing.T) {
	mock := &MockLocalLoginApp{}
	mock.Setup(t)
	defer mock.Teardown()

	ctx := context.Background()
	id := uuid.New()
	email := "test@test.com"
	password := "password123"

	t.Run("Successful Login", func(t *testing.T) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mock.mockRepository.EXPECT().GetUserByEmail(email).Return(entUser, nil)
		mock.mockTokenProvider.EXPECT().GenerateAccessToken(ctx, id.String()).Return("access-token", int64(3600), nil)
		mock.mockTokenProvider.EXPECT().GenerateRefreshToken(ctx, id.String()).Return("refresh-token", int64(7200), nil)

		loginToken, err := mock.app.Login(ctx, email, password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(loginToken); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
		if loginToken.AccessToken == "" || loginToken.RefreshToken == "" {
			t.Fatal("expected non-empty access and refresh tokens")
		}
	})

	t.Run("Login Failure", func(t *testing.T) {
		mock.mockRepository.EXPECT().GetUserByEmail(email).Return(nil, errors.New("user not found"))

		_, err := mock.app.Login(ctx, email, password)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   "$2a$10$invalidhash", // Invalid hash
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mock.mockRepository.EXPECT().GetUserByEmail(email).Return(entUser, nil)

		_, err := mock.app.Login(ctx, email, password)
		if !errors.Is(err, logindto.ErrInvalidCredentials) {
			t.Fatalf("expected invalid credentials error, got %v", err)
		}
	})
}
