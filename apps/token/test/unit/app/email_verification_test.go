package app_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	token "mandacode.com/accounts/token/internal/app"
	mock_svcdomain "mandacode.com/accounts/token/test/mock/domain/service"
)

type MockEmailVerificationTokenApp struct {
	crtl           *gomock.Controller
	tokenGenerator *mock_svcdomain.MockTokenGenerator
	app            *token.EmailVerificationTokenApp
}

func (m *MockEmailVerificationTokenApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.crtl = ctrl
	m.tokenGenerator = mock_svcdomain.NewMockTokenGenerator(ctrl)
	m.app = token.NewEmailVerificationTokenApp(m.tokenGenerator)
}

func (m *MockEmailVerificationTokenApp) Teardown() {
	m.crtl.Finish()
}

func TestEmailVerificationTokenApp_GenerateToken(t *testing.T) {
	mockApp := &MockEmailVerificationTokenApp{}
	mockApp.Setup(t)
	defer mockApp.Teardown()

	mockUserID := uuid.New().String()
	mockToken := "mock-token"
	mockEmail := "test@test.com"
	mockCode := "123456"
	mockExpiration := int64(3600) // 1 hour in seconds
	claims := map[string]string{"sub": mockUserID, "email": mockEmail, "code": mockCode}

	t.Run("GenerateToken_Success", func(t *testing.T) {
		mockApp.tokenGenerator.EXPECT().GenerateToken(claims).Return(mockToken, mockExpiration, nil).Times(1)

		token, exp, err := mockApp.app.GenerateToken(mockUserID, mockEmail, mockCode)
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

		token, exp, err := mockApp.app.GenerateToken(mockUserID, mockEmail, mockCode)
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

func TestEmailVerificationTokenApp_VerifyToken(t *testing.T) {
	mockApp := &MockEmailVerificationTokenApp{}
	mockApp.Setup(t)
	defer mockApp.Teardown()

	mockUserID := uuid.New().String()
	mockToken := "mock-token"
	mockEmail := "test@test.com"
	mockCode := "123456"
	claims := map[string]string{"sub": mockUserID, "email": mockEmail, "code": mockCode}

	t.Run("VerifyToken_Success", func(t *testing.T) {
		mockApp.tokenGenerator.EXPECT().VerifyToken(mockToken).Return(claims, nil).Times(1)

		userID, email, code, err := mockApp.app.VerifyToken(mockToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if userID == nil || *userID != mockUserID {
			t.Errorf("expected user ID %s, got %v", mockUserID, userID)
		}
		if email == nil || *email != mockEmail {
			t.Errorf("expected email %s, got %v", mockEmail, email)
		}
		if code == nil || *code != mockCode {
			t.Errorf("expected code %s, got %v", mockCode, code)
		}
	})

	t.Run("VerifyToken_Error", func(t *testing.T) {
		expectedError := errors.New("token verification error")
		mockApp.tokenGenerator.EXPECT().VerifyToken(mockToken).Return(nil, expectedError).Times(1)

		userID, email, code, err := mockApp.app.VerifyToken(mockToken)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if userID != nil {
			t.Errorf("expected nil user ID, got %v", userID)
		}
		if email != nil {
			t.Errorf("expected nil email, got %v", email)
		}
		if code != nil {
			t.Errorf("expected nil code, got %v", code)
		}
	})

	t.Run("VerifyToken_NoClaims", func(t *testing.T) {
		expectedError := errors.New("no claims found in access token")
		mockApp.tokenGenerator.EXPECT().VerifyToken(mockToken).Return(nil, expectedError).Times(1)

		userID, email, code, err := mockApp.app.VerifyToken(mockToken)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if userID != nil {
			t.Errorf("expected nil user ID, got %v", userID)
		}
		if email != nil {
			t.Errorf("expected nil email, got %v", email)
		}
		if code != nil {
			t.Errorf("expected nil code, got %v", code)
		}
	})
}
