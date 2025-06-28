package usersvc_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
	userdomain "mandacode.com/accounts/auth/internal/domain/service/user"
	usersvc "mandacode.com/accounts/auth/internal/infra/service/user"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
)

type MockOAuthUserService struct {
	mockRepo *mock_repodomain.MockOAuthUserRepository
	svc      userdomain.OAuthUserService
}

func (s *MockOAuthUserService) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	s.mockRepo = mock_repodomain.NewMockOAuthUserRepository(ctrl)
	s.svc = usersvc.NewOAuthUserService(s.mockRepo)
}

func (s *MockOAuthUserService) Teardown() {
	s.mockRepo = nil
	s.svc = nil
}

func TestOAuthCreateUser(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	provider := oauthuser.ProviderGoogle
	providerID := "google-id-123"

	t.Run("Successful User Creation", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().CreateUser(id, provider, providerID, email, nil, nil).Return(entUser, nil)

		user, err := mock.svc.CreateUser(id, provider, providerID, email, nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("Error on User Creation", func(t *testing.T) {
		mock.mockRepo.EXPECT().CreateUser(id, provider, providerID, email, nil, nil).Return(nil, errors.New("error creating user"))

		user, err := mock.svc.CreateUser(id, provider, providerID, email, nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on User Creation with Invalid Provider", func(t *testing.T) {
		invalidProvider := oauthuser.Provider("invalid")
		mock.mockRepo.EXPECT().CreateUser(id, invalidProvider, providerID, email, nil, nil).Return(nil, errors.New("invalid provider"))

		user, err := mock.svc.CreateUser(id, invalidProvider, providerID, email, nil, nil)
		if err == nil {
			t.Fatal("expected error for invalid provider, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestOAuthGetUserByProvider(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	provider := oauthuser.ProviderGoogle
	providerID := "google-id-123"

	t.Run("Successful User Retrieval", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(entUser, nil)

		user, err := mock.svc.GetUserByProvider(provider, providerID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(nil, errors.New("user not found"))

		user, err := mock.svc.GetUserByProvider(provider, providerID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on User Retrieval", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByProvider(provider, providerID).Return(nil, errors.New("error retrieving user"))

		user, err := mock.svc.GetUserByProvider(provider, providerID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestOAuthGetUserByUserID(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	provider := oauthuser.ProviderGoogle
	providerID := "google-id-123"

	t.Run("Successful User Retrieval by UserID", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().GetUserByUserID(id, provider).Return(entUser, nil)

		user, err := mock.svc.GetUserByUserID(id, provider)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("User Not Found by UserID", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByUserID(id, provider).Return(nil, errors.New("user not found"))

		user, err := mock.svc.GetUserByUserID(id, provider)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})

	t.Run("Error on User Retrieval by UserID", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByUserID(id, provider).Return(nil, errors.New("error retrieving user"))

		user, err := mock.svc.GetUserByUserID(id, provider)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestOAuthDeleteUser(t *testing.T) {
	mock := &MockOAuthUserService{}
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

func TestOAuthDeleteUserByProvider(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()

	t.Run("Successful User Deletion by Provider", func(t *testing.T) {
		mock.mockRepo.EXPECT().DeleteUser(id).Return(nil)

		deletedUser, err := mock.svc.DeleteUser(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if deletedUser.ID != id {
			t.Errorf("expected deleted user ID %v, got %v", id, deletedUser.ID)
		}
	})

	t.Run("Error on User Deletion by Provider", func(t *testing.T) {
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

func TestOAuthUpdateUserBase(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	providerID := "google-id-123"
	email := "test@test.com"
	isActive := true
	isVerified := true

	t.Run("Successful User Update", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   isActive,
			IsVerified: isVerified,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Email:      email,
			Provider:   provider,
			ProviderID: providerID,
			IsActive:   isActive,
			IsVerified: isVerified,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().UpdateUser(id, provider, &providerID, &email, nil, &isVerified).Return(entUser, nil)

		user, err := mock.svc.UpdateUserBase(id, provider, providerID, email, isVerified)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("Error on User Update", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, provider, &providerID, &email, nil, &isVerified).Return(nil, errors.New("error updating user"))

		user, err := mock.svc.UpdateUserBase(id, provider, providerID, email, isVerified)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestOAuthUpdateActiveStatus(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	isActive := true

	t.Run("Successful Active Status Update", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:        id,
			Provider:  provider,
			IsActive:  isActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		dtoUser := &dto.OAuthUser{
			ID:        id,
			Provider:  provider,
			IsActive:  isActive,
			CreatedAt: entUser.CreatedAt,
			UpdatedAt: entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().UpdateUser(id, provider, nil, nil, &isActive, nil).Return(entUser, nil)

		user, err := mock.svc.UpdateActiveStatus(id, provider, isActive)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("Error on Active Status Update", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, provider, nil, nil, &isActive, nil).Return(nil, errors.New("error updating active status"))

		user, err := mock.svc.UpdateActiveStatus(id, provider, isActive)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}

func TestOAuthUpdateVerifiedStatus(t *testing.T) {
	mock := &MockOAuthUserService{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	provider := oauthuser.ProviderGoogle
	isVerified := true

	t.Run("Successful Verified Status Update", func(t *testing.T) {
		entUser := &ent.OAuthUser{
			ID:         id,
			Provider:   provider,
			IsVerified: isVerified,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		dtoUser := &dto.OAuthUser{
			ID:         id,
			Provider:   provider,
			IsVerified: isVerified,
			CreatedAt:  entUser.CreatedAt,
			UpdatedAt:  entUser.UpdatedAt,
		}

		mock.mockRepo.EXPECT().UpdateUser(id, provider, nil, nil, nil, &isVerified).Return(entUser, nil)

		user, err := mock.svc.UpdateVerifiedStatus(id, provider, isVerified)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(user, dtoUser) {
			t.Errorf("expected user %v, got %v", dtoUser, user)
		}
	})

	t.Run("Error on Verified Status Update", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, provider, nil, nil, nil, &isVerified).Return(nil, errors.New("error updating verified status"))

		user, err := mock.svc.UpdateVerifiedStatus(id, provider, isVerified)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}
