package user_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	localuser "mandacode.com/accounts/auth/internal/app/user/local"
	mock_repodomain "mandacode.com/accounts/auth/test/mock/domain/repository"
)

type MockLocalUserApp struct {
	mockRepo *mock_repodomain.MockLocalUserRepository
	validate *validator.Validate
	app      localuser.LocalUserApp
}

func (m *MockLocalUserApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.mockRepo = mock_repodomain.NewMockLocalUserRepository(ctrl)
	m.validate = validator.New()
	m.app = localuser.NewLocalUserApp(m.mockRepo)
}

func (m *MockLocalUserApp) Teardown() {
	m.mockRepo = nil
	m.validate = nil
	m.app = nil
}

func TestLocalUserApp_CreateUser(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	password := "password123"

	t.Run("Successful User Creation", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   password,
			IsActive:   true,
			IsVerified: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().CreateUser(id, email, gomock.Any(), nil, nil).Return(entUser, nil).Times(1)

		user, err := mock.app.CreateUser(id, email, password, nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
}

func TestLocalUserApp_DeleteUser(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()

	t.Run("Successful User Deletion", func(t *testing.T) {
		mock.mockRepo.EXPECT().DeleteUser(id).Return(nil).Times(1)

		err := mock.app.DeleteUser(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestLocalUserApp_GetUser(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	t.Run("Successful User Retrieval", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   "hashed-password",
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByID(id).Return(entUser, nil).Times(1)

		user, err := mock.app.GetUser(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("User Not Found", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByID(id).Return(nil, nil).Times(1)

		user, err := mock.app.GetUser(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil, got non-nil user")
		}
	})
}

func TestLocalUserApp_GetUserByEmail(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	email := "test@test.com"
	t.Run("Successful User Retrieval by Email", func(t *testing.T) {
		id := uuid.New()
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   "hashed-password",
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(entUser, nil).Times(1)

		user, err := mock.app.GetUserByEmail(email)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("User Not Found by Email", func(t *testing.T) {
		mock.mockRepo.EXPECT().GetUserByEmail(email).Return(nil, nil).Times(1)

		user, err := mock.app.GetUserByEmail(email)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil, got non-nil user")
		}
	})
}

func TestLocalUserApp_UpdateActiveStatus(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	email := "test@test.com"
	isActive := true

	t.Run("Successful Update Active Status", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      email,
			Password:   "hashed-password",
			IsActive:   isActive,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, &isActive, nil).Return(entUser, nil).Times(1)

		user, err := mock.app.UpdateActiveStatus(id, isActive)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("Update Active Status Error", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, nil, nil, &isActive, nil).Return(nil, nil).Times(1)

		user, err := mock.app.UpdateActiveStatus(id, isActive)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil, got non-nil user")
		}
	})
}

func TestLocalUserApp_UpdateEmail(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	newEmail := "new@test.com"
	t.Run("Successful Update Email", func(t *testing.T) {
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      newEmail,
			Password:   "hashed-password",
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().UpdateUser(id, &newEmail, nil, nil, nil).Return(entUser, nil).Times(1)

		user, err := mock.app.UpdateEmail(id, newEmail)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("Update Email Error", func(t *testing.T) {
		mock.mockRepo.EXPECT().UpdateUser(id, &newEmail, nil, nil, nil).Return(nil, nil).Times(1)

		user, err := mock.app.UpdateEmail(id, newEmail)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != nil {
			t.Fatal("expected user to be nil, got non-nil user")
		}
	})
}

func TestLocalUserApp_UpdatePassword(t *testing.T) {
	mock := &MockLocalUserApp{}
	mock.Setup(t)
	defer mock.Teardown()

	id := uuid.New()
	currentPassword := "current-password"
	newPassword := "new-password"

	t.Run("Successful Update Password", func(t *testing.T) {
		hashedCurrentPassword, err := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("failed to hash current password: %v", err)
		}
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      "test@test.com",
			Password:   string(hashedCurrentPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mock.mockRepo.EXPECT().GetUserByID(id).Return(entUser, nil).Times(1)
		mock.mockRepo.EXPECT().UpdateUser(id, nil, gomock.Any(), nil, nil).Return(entUser, nil).Times(1)
		user, err := mock.app.UpdatePassword(id, currentPassword, newPassword)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.validate.Struct(user); err != nil {
			t.Fatalf("validation failed: %v", err)
		}
	})
	t.Run("Update Password Error - Incorrect Current Password", func(t *testing.T) {
		hashedCurrentPassword, err := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("failed to hash current password: %v", err)
		}
		entUser := &ent.LocalUser{
			ID:         id,
			Email:      "test@test.com",
			Password:   string(hashedCurrentPassword),
			IsActive:   true,
			IsVerified: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mock.mockRepo.EXPECT().GetUserByID(id).Return(entUser, nil).Times(1)
		_, err = mock.app.UpdatePassword(id, "wrong-password", newPassword)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
