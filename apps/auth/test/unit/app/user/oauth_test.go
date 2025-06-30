package user_test

import (
	"errors"
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
	protoutil "mandacode.com/accounts/auth/internal/util/proto"
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

func TestOAuthUserApp_GetUser(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle

	t.Run("Successful User Retrieval", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: "google-id-123",
			Email:      "test@test.com",
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().
			GetUserByUserID(id, provider).
			Return(entUser, nil)
		user, err := mock.app.GetUser(id, provider)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			GetUserByUserID(id, provider).
			Return(nil, nil)
		user, err := mock.app.GetUser(id, provider)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil")
		}
	})
}

func TestOAuthUserApp_SyncUser(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	oauthAccessToken := "valid-access-token"
	t.Run("Successful User Sync", func(t *testing.T) {
		userInfo := &oauthdto.OAuthUserInfo{
			ProviderID:    "google-id-123",
			Email:         "test@test.com",
			Name:          "Test User",
			EmailVerified: true,
		}
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: userInfo.ProviderID,
			Email:      userInfo.Email,
			IsActive:   userInfo.EmailVerified,
			IsVerified: !userInfo.EmailVerified,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		(*mock.mockOAuthProviders)[provider].(*mock_oauthdomain.MockOAuthProvider).
			EXPECT().
			GetUserInfo(oauthAccessToken).
			Return(userInfo, nil)
		mock.mockRepo.EXPECT().
			UpdateUser(id, provider, &userInfo.ProviderID, &userInfo.Email, nil, &userInfo.EmailVerified).
			Return(entUser, nil)

		user, err := mock.app.SyncUser(id, provider, oauthAccessToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("Unsupported Provider", func(t *testing.T) {
		unsupportedProvider := oauthuser.Provider("unsupported")
		_, err := mock.app.SyncUser(id, unsupportedProvider, oauthAccessToken)
		if !errors.Is(err, protoutil.ErrUnsupportedProvider) {
			t.Fatalf("expected unsupported provider error, got %v", err)
		}
	})
}

func TestOAuthUserApp_UpdateActiveStatus(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	isActive := true

	t.Run("Successful Update Active Status", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: "google-id-123",
			Email:      "test@test.com",
			IsActive:   isActive,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().
			UpdateUser(id, provider, nil, nil, &isActive, nil).
			Return(entUser, nil)
		user, err := mock.app.UpdateActiveStatus(id, provider, isActive)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("Update Active Status Error", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			UpdateUser(id, provider, nil, nil, &isActive, nil).
			Return(nil, nil)
		user, err := mock.app.UpdateActiveStatus(id, provider, isActive)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil")
		}
	})
}

func TestOAuthUserApp_UpdateVerificationStatus(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	isVerified := true

	t.Run("Successful Update Verification Status", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: "google-id-123",
			Email:      "test@test.com",
			IsActive:   true,
			IsVerified: isVerified,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mock.mockRepo.EXPECT().
			UpdateUser(id, provider, nil, nil, nil, &isVerified).
			Return(entUser, nil)
		user, err := mock.app.UpdateVerificationStatus(id, provider, isVerified)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("Update Verification Status Error", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			UpdateUser(id, provider, nil, nil, nil, &isVerified).
			Return(nil, nil)
		user, err := mock.app.UpdateVerificationStatus(id, provider, isVerified)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil")
		}
	})
}

func TestOAuthUserApp_DeleteUser(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle

	t.Run("Successful User Deletion", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			DeleteUserByProvider(id, provider).
			Return(nil)
		err := mock.app.DeleteUser(id, provider)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
	t.Run("Delete User Error", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			DeleteUserByProvider(id, provider).
			Return(errors.New("delete error"))
		err := mock.app.DeleteUser(id, provider)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestOAuthUserApp_DeleteAllProviders(t *testing.T) {
	mock := &MockOAuthUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()

	t.Run("Successful User Deletion by Provider", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			DeleteUser(id).
			Return(nil)
		err := mock.app.DeleteAllProviders(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
	t.Run("Delete User by Provider Error", func(t *testing.T) {
		mock.mockRepo.EXPECT().
			DeleteUser(id).
			Return(errors.New("delete error"))
		err := mock.app.DeleteAllProviders(id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
