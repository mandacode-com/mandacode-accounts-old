package login_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/internal/app/login"
	"mandacode.com/accounts/auth/internal/domain/dto"
	mock_logindomain "mandacode.com/accounts/auth/test/mock/domain/service/login"
	mock_tokendomain "mandacode.com/accounts/auth/test/mock/domain/service/token"
)

type MockLocalLoginApp struct {
	mockTokenService      *mock_tokendomain.MockTokenService
	mockLocalLoginService *mock_logindomain.MockLocalLoginService
	app                   *login.LocalLoginApp
}

func (m *MockLocalLoginApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.mockTokenService = mock_tokendomain.NewMockTokenService(ctrl)
	m.mockLocalLoginService = mock_logindomain.NewMockLocalLoginService(ctrl)
	m.app = login.NewLocalLoginApp(m.mockTokenService, m.mockLocalLoginService)
}

func (m *MockLocalLoginApp) Teardown() {
	m.mockTokenService = nil
	m.mockLocalLoginService = nil
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
		dtoUser := &dto.LocalUser{
			ID:         id,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mock.mockLocalLoginService.EXPECT().LoginLocalUser(ctx, email, password).Return(dtoUser, nil)
		mock.mockTokenService.EXPECT().GenerateAccessToken(ctx, id.String()).Return("access-token", time.Now().Unix(), nil)
		mock.mockTokenService.EXPECT().GenerateRefreshToken(ctx, id.String()).Return("refresh-token", time.Now().Unix(), nil)

		userID, accessToken, refreshToken, err := mock.app.LoginLocalUser(ctx, email, password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if *userID != id.String() {
			t.Errorf("expected user ID %s, got %s", id.String(), *userID)
		}
		if accessToken == nil || refreshToken == nil {
			t.Error("expected non-nil access and refresh tokens")
		}
	})

	t.Run("Login Failure", func(t *testing.T) {
		mock.mockLocalLoginService.EXPECT().LoginLocalUser(ctx, email, password).Return(nil, errors.New("login failed"))

		userID, accessToken, refreshToken, err := mock.app.LoginLocalUser(ctx, email, password)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if userID != nil || accessToken != nil || refreshToken != nil {
			t.Error("expected nil user ID and tokens on error")
		}
	})
}
