package user_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
	oauthuserapp "mandacode.com/accounts/auth/internal/app/user/oauth"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/oauth"
	oauthdto "mandacode.com/accounts/auth/internal/infra/oauth/dto"
	mock_oauthdomain "mandacode.com/accounts/auth/test/mock/domain/oauth"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
)

type MockOAuthUserApp struct {
	mockRepo           *mock_repodomain.MockOAuthUserRepository
	mockOAuthProviders *map[oauthuser.Provider]oauthdomain.OAuthProvider
	validate           *validator.Validate
	app                oauthuserapp.OAuthUserApp
}

func (m *MockOAuthUserApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.mockRepo = mock_repodomain.NewMockOAuthUserRepository(ctrl)
	m.validate = validator.New()

	m.mockOAuthProviders = &map[oauthuser.Provider]oauthdomain.OAuthProvider{
		oauthuser.ProviderGoogle: mock_oauthdomain.NewMockOAuthProvider(ctrl),
		oauthuser.ProviderNaver:  mock_oauthdomain.NewMockOAuthProvider(ctrl),
		oauthuser.ProviderKakao:  mock_oauthdomain.NewMockOAuthProvider(ctrl),
	}

	m.app = oauthuserapp.NewOAuthUserApp(
		m.mockOAuthProviders,
		m.mockRepo,
	)
}

func (m *MockOAuthUserApp) Teardown() {
	m.mockRepo = nil
	m.validate = nil
	m.mockOAuthProviders = nil
	m.app = nil
}

func TestOAuthUserApp_CreateUser(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	provider := oauthuser.ProviderGoogle
	oauthAccessToken := "valid-access-token"

	t.Run("Successful User Creation", func(t *testing.T) {
		userInfo := &oauthdto.OAuthUserInfo{
			ProviderID:    "google-id-123",
			Email:         email,
			Name:          "Test User",
			EmailVerified: true,
		}
		entUser := &ent.OAuthUser{
			ID:         id,
			Email:      userInfo.Email,
			Provider:   oauthuser.ProviderGoogle,
			ProviderID: userInfo.ProviderID,
			IsActive:   userInfo.EmailVerified,
			IsVerified: userInfo.EmailVerified,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		(*mock.mockOAuthProviders)[provider].(*mock_oauthdomain.MockOAuthProvider).
			EXPECT().
			GetUserInfo(oauthAccessToken).
			Return(userInfo, nil)

		mock.mockRepo.EXPECT().
			CreateUser(id, provider, userInfo.ProviderID, userInfo.Email, nil, &userInfo.EmailVerified).
			Return(entUser, nil)

		user, err := mock.app.CreateUser(id, provider, oauthAccessToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
}
