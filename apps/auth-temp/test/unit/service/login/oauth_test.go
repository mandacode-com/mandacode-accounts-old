package loginsvc_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
	logindomain "mandacode.com/accounts/auth/internal/domain/service/login"
	loginsvc "mandacode.com/accounts/auth/internal/infra/service/login"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
)

type MockOAuthAuthService struct {
	mockRepo *mock_repodomain.MockOAuthUserRepository
	svc      logindomain.OAuthLoginService
}

func (s *MockOAuthAuthService) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	s.mockRepo = mock_repodomain.NewMockOAuthUserRepository(ctrl)
	s.svc = loginsvc.NewOAuthLoginService(s.mockRepo)
}

func (s *MockOAuthAuthService) Teardown() {
	s.mockRepo = nil
	s.svc = nil
}

func TestLoginOAuthUser(t *testing.T) {
	mock := &MockOAuthAuthService{}
	mock.Setup(t)
	defer mock.Teardown()

	ctx := context.Background()
	id := uuid.New()
	email := "test@test.com"
	provider := oauthuser.ProviderGoogle
	providerID := "google-id-123"

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
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}
		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(entUser, nil)

		user, err := mock.svc.LoginOAuthUser(ctx, provider, providerID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(nil, errors.New("user not found"))

		user, err := mock.svc.LoginOAuthUser(ctx, provider, providerID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("User Not Active", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
			IsActive:   false,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(entUser, nil)

		user, err := mock.svc.LoginOAuthUser(ctx, provider, providerID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("User Not Verified", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
			IsActive:   true,
			IsVerified: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(entUser, nil)

		user, err := mock.svc.LoginOAuthUser(ctx, provider, providerID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}
