package app_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	token "mandacode.com/accounts/token/internal/app"
	mock_tokengendomain "mandacode.com/accounts/token/test/mock/domain/token"
)

type MockAccessTokenApp struct {
	crtl           *gomock.Controller
	tokenGenerator *mock_tokengendomain.MockTokenGenerator
	app            *token.AccessTokenApp
}

func (m *MockAccessTokenApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.crtl = ctrl
	m.tokenGenerator = mock_tokengendomain.NewMockTokenGenerator(ctrl)
	m.app = token.NewAccessTokenApp(m.tokenGenerator)
}

func (m *MockAccessTokenApp) Teardown() {
	m.crtl.Finish()
}

func TestAccessTokenApp_GenerateToken(t *testing.T) {
	mockApp := &MockAccessTokenApp{}
	mockApp.Setup(t)
	defer mockApp.Teardown()

	mockUserID := uuid.New().String()
	mockToken := "mock-token"
	mockExpiration := int64(3600) // 1 hour in seconds
	claims := map[string]string{"sub": mockUserID}

	t.Run("GenerateToken_Success", func(t *testing.T) {
		mockApp.tokenGenerator.EXPECT().GenerateToken(claims).Return(mockToken, mockExpiration, nil).Times(1)

		token, exp, err := mockApp.app.GenerateToken(mockUserID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token != mockToken {
			t.Errorf("expected token %s, got %s", mockToken, token)
		}
		if exp != mockExpiration {
			t.Errorf("expected expiration %d, got %d", mockExpiration, exp)
		}
	})

	t.Run("GenerateToken_Error", func(t *testing.T) {
		expectedError := errors.New("token generation error")
		mockApp.tokenGenerator.EXPECT().GenerateToken(claims).Return("", int64(0), expectedError).Times(1)

		token, exp, err := mockApp.app.GenerateToken(mockUserID)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if token != "" {
			t.Errorf("expected empty token, got %s", token)
		}
		if exp != 0 {
			t.Errorf("expected expiration 0, got %d", exp)
		}
	})
}

func TestAccessTokenApp_VerifyToken(t *testing.T) {
	mockApp := &MockAccessTokenApp{}
	mockApp.Setup(t)
	defer mockApp.Teardown()

	mockUserID := uuid.New().String()
	mockToken := "mock-token"
	claims := map[string]string{"sub": mockUserID}

	t.Run("VerifyToken_Success", func(t *testing.T) {
		mockApp.tokenGenerator.EXPECT().VerifyToken(mockToken).Return(claims, nil).Times(1)

		userID, err := mockApp.app.VerifyToken(mockToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if userID == nil || *userID != mockUserID {
			t.Errorf("expected user ID %s, got %v", mockUserID, userID)
		}
	})

	t.Run("VerifyToken_Error", func(t *testing.T) {
		expectedError := errors.New("token verification error")
		mockApp.tokenGenerator.EXPECT().VerifyToken(mockToken).Return(nil, expectedError).Times(1)

		userID, err := mockApp.app.VerifyToken(mockToken)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if userID != nil {
			t.Errorf("expected nil user ID, got %v", userID)
		}
	})

	t.Run("VerifyToken_NoClaims", func(t *testing.T) {
		expectedError := errors.New("no claims found in access token")
		mockApp.tokenGenerator.EXPECT().VerifyToken(mockToken).Return(nil, expectedError).Times(1)

		userID, err := mockApp.app.VerifyToken(mockToken)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if userID != nil {
			t.Errorf("expected nil user ID, got %v", userID)
		}
	})
}
