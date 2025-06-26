package login_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/app/login"
	"mandacode.com/accounts/auth/internal/domain/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	mock_logindomain "mandacode.com/accounts/auth/test/mock/domain/service/login"
	mock_oauthdomain "mandacode.com/accounts/auth/test/mock/domain/service/oauth"
	mock_tokendomain "mandacode.com/accounts/auth/test/mock/domain/service/token"
	providerv1 "mandacode.com/accounts/proto/common/provider/v1"
)

type MockOAuthLoginApp struct {
	mockProviders         map[oauthuser.Provider]oauthdomain.OAuthService
	mockTokenService      *mock_tokendomain.MockTokenService
	mockOAuthLoginService *mock_logindomain.MockOAuthLoginService
	app                   *login.OAuthLoginApp
}

func (m *MockOAuthLoginApp) Setup() {
	ctrl := gomock.NewController(nil)
	m.mockProviders = map[oauthuser.Provider]oauthdomain.OAuthService{
		oauthuser.ProviderGoogle: mock_oauthdomain.NewMockOAuthService(ctrl),
		oauthuser.ProviderNaver:  mock_oauthdomain.NewMockOAuthService(ctrl),
		oauthuser.ProviderKakao:  mock_oauthdomain.NewMockOAuthService(ctrl),
	}
	m.mockTokenService = mock_tokendomain.NewMockTokenService(ctrl)
	m.mockOAuthLoginService = mock_logindomain.NewMockOAuthLoginService(ctrl)
	m.app = login.NewOAuthLoginApp(&m.mockProviders, m.mockTokenService, m.mockOAuthLoginService)
}

func (m *MockOAuthLoginApp) Teardown() {
	m.mockTokenService = nil
	m.mockOAuthLoginService = nil
	m.app = nil
}

func (m *MockOAuthLoginApp) GetApp() *login.OAuthLoginApp {
	return m.app
}

func TestOAuthAuthApp_LoginOAuthUser(t *testing.T) {
	mock := &MockOAuthLoginApp{}
	mock.Setup()
	defer mock.Teardown()

	ctx := context.Background()
	id := uuid.New()
	provider := providerv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE
	providerEnum := oauthuser.ProviderGoogle
	email := "test@test.com"
	providerID := "google-id-123"
	oauthAccessToken := "valid-access-token"

	t.Run("Successful Login", func(t *testing.T) {
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Provider:   providerEnum,
			ProviderID: providerID,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		oauthUserInfo := &dto.OAuthUserInfo{
			ProviderID:    providerID,
			Email:         email,
			Name:          "Test User",
			EmailVerified: true,
		}

		mock.mockProviders[providerEnum].(*mock_oauthdomain.MockOAuthService).EXPECT().GetUserInfo(oauthAccessToken).Return(oauthUserInfo, nil)
		mock.mockOAuthLoginService.EXPECT().LoginOAuthUser(ctx, providerEnum, providerID).Return(dtoUser, nil)
		mock.mockTokenService.EXPECT().GenerateAccessToken(ctx, id.String()).Return("access-token", time.Now().Unix(), nil)
		mock.mockTokenService.EXPECT().GenerateRefreshToken(ctx, id.String()).Return("refresh-token", time.Now().Unix(), nil)

		userID, accessToken, refreshToken, err := mock.app.LoginOAuthUser(ctx, provider, oauthAccessToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if userID == nil || *userID != id.String() {
			t.Errorf("expected user ID %s, got %s", id.String(), *userID)
		}
		if accessToken == nil || *accessToken != "access-token" {
			t.Errorf("expected access token 'access-token', got %s", *accessToken)
		}
		if refreshToken == nil || *refreshToken != "refresh-token" {
			t.Errorf("expected refresh token 'refresh-token', got %s", *refreshToken)
		}
	})

	t.Run("Unsupported Provider", func(t *testing.T) {
		unsupportedProvider := providerv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED
		_, _, _, err := mock.app.LoginOAuthUser(ctx, unsupportedProvider, oauthAccessToken)
		if err == nil {
			t.Fatal("expected error for unsupported provider, got nil")
		}
	})

	t.Run("GetUserInfo Error", func(t *testing.T) {
		mock.mockProviders[providerEnum].(*mock_oauthdomain.MockOAuthService).EXPECT().GetUserInfo(oauthAccessToken).Return(nil, errors.New("failed to get user info"))

		_, _, _, err := mock.app.LoginOAuthUser(ctx, provider, oauthAccessToken)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
