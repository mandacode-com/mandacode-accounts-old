package login_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
	oauthlogin "mandacode.com/accounts/auth/internal/app/login/oauth"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/oauth"
	oauthdto "mandacode.com/accounts/auth/internal/infra/oauth/dto"
	protoutil "mandacode.com/accounts/auth/internal/util/proto"
	mock_oauthdomain "mandacode.com/accounts/auth/test/mock/domain/oauth"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
	mock_tokendomain "mandacode.com/accounts/auth/test/mock/domain/token"
)

type MockOAuthLoginApp struct {
	mockTokenProvider  *mock_tokendomain.MockTokenProvider
	mockOAuthProviders *map[oauthuser.Provider]oauthdomain.OAuthProvider
	mockRepository     *mock_repodomain.MockOAuthUserRepository
	validate           *validator.Validate
	app                oauthlogin.OAuthLoginApp
}

func (m *MockOAuthLoginApp) Setup() {
	ctrl := gomock.NewController(nil)
	m.mockTokenProvider = mock_tokendomain.NewMockTokenProvider(ctrl)
	m.mockOAuthProviders = &map[oauthuser.Provider]oauthdomain.OAuthProvider{
		oauthuser.ProviderGoogle: mock_oauthdomain.NewMockOAuthProvider(ctrl),
		oauthuser.ProviderNaver:  mock_oauthdomain.NewMockOAuthProvider(ctrl),
		oauthuser.ProviderKakao:  mock_oauthdomain.NewMockOAuthProvider(ctrl),
	}
	m.mockRepository = mock_repodomain.NewMockOAuthUserRepository(ctrl)
	m.validate = validator.New()
	m.app = oauthlogin.NewOAuthLoginApp(
		m.mockTokenProvider,
		m.mockOAuthProviders,
		m.mockRepository,
	)
}

func (m *MockOAuthLoginApp) Teardown() {
	m.mockTokenProvider = nil
	m.mockOAuthProviders = nil
	m.mockRepository = nil
	m.validate = nil
	m.app = nil
}

func TestOAuthAuthApp_LoginOAuthUser(t *testing.T) {
	mock := &MockOAuthLoginApp{}
	mock.Setup()
	defer mock.Teardown()

	ctx := context.Background()
	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	email := "test@test.com"
	providerID := "google-id-123"
	oauthAccessToken := "valid-access-token"

	t.Run("Successful Login", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		oauthUserInfo := &oauthdto.OAuthUserInfo{
			ProviderID:    providerID,
			Email:         email,
			Name:          "Test User",
			EmailVerified: true,
		}

		(*mock.mockOAuthProviders)[provider].(*mock_oauthdomain.MockOAuthProvider).EXPECT().GetUserInfo(oauthAccessToken).Return(oauthUserInfo, nil)
		mock.mockRepository.EXPECT().GetUserByProvider(provider, providerID).Return(entUser, nil)
		mock.mockTokenProvider.EXPECT().GenerateAccessToken(ctx, id.String()).Return("access-token", time.Now().Unix(), nil)
		mock.mockTokenProvider.EXPECT().GenerateRefreshToken(ctx, id.String()).Return("refresh-token", time.Now().Unix(), nil)

		loginToken, err := mock.app.Login(ctx, provider, oauthAccessToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(loginToken); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})

	t.Run("Unsupported Provider", func(t *testing.T) {
		unsupportedProvider := oauthuser.Provider("unsupported")
		_, err := mock.app.Login(ctx, unsupportedProvider, oauthAccessToken)
		if !errors.Is(err, protoutil.ErrUnsupportedProvider) {
			t.Fatalf("expected unsupported provider error, got %v", err)
		}
	})

	t.Run("GetUserInfo Error", func(t *testing.T) {
		(*mock.mockOAuthProviders)[provider].(*mock_oauthdomain.MockOAuthProvider).EXPECT().GetUserInfo(oauthAccessToken).Return(nil, errors.New("failed to get user info"))

		_, err := mock.app.Login(ctx, provider, oauthAccessToken)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
