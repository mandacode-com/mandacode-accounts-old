package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/internal/app/auth"
	"mandacode.com/accounts/auth/internal/domain/dto"
	mock_authdomain "mandacode.com/accounts/auth/test/mock/domain/service/auth"
	mock_tokendomain "mandacode.com/accounts/auth/test/mock/domain/service/token"
)

type MockLocalAuthApp struct {
	mockTokenService     *mock_tokendomain.MockTokenService
	mockLocalAuthService *mock_authdomain.MockLocalAuthService
	app                  *auth.LocalAuthApp
}

func (m *MockLocalAuthApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.mockTokenService = mock_tokendomain.NewMockTokenService(ctrl)
	m.mockLocalAuthService = mock_authdomain.NewMockLocalAuthService(ctrl)
	m.app = auth.NewLocalAuthApp(m.mockTokenService, m.mockLocalAuthService)
}

func (m *MockLocalAuthApp) Teardown() {
	m.mockTokenService = nil
	m.mockLocalAuthService = nil
	m.app = nil
}

func TestLocalAuthApp_LoginLocalUser(t *testing.T) {
	mock := &MockLocalAuthApp{}
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

		mock.mockLocalAuthService.EXPECT().LoginLocalUser(ctx, email, password).Return(dtoUser, nil)
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
		mock.mockLocalAuthService.EXPECT().LoginLocalUser(ctx, email, password).Return(nil, errors.New("login failed"))

		userID, accessToken, refreshToken, err := mock.app.LoginLocalUser(ctx, email, password)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if userID != nil || accessToken != nil || refreshToken != nil {
			t.Error("expected nil user ID and tokens on error")
		}
	})
}
