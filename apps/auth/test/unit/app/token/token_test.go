package token_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"mandacode.com/accounts/auth/internal/app/token"
	mock_tokendomain "mandacode.com/accounts/auth/test/mock/domain/token"
)

type MockTokenApp struct {
	mockTokenProvider *mock_tokendomain.MockTokenProvider
	app               token.TokenApp
}

func (m *MockTokenApp) Setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	m.mockTokenProvider = mock_tokendomain.NewMockTokenProvider(ctrl)
	m.app = token.NewTokenApp(m.mockTokenProvider)
}

func (m *MockTokenApp) Teardown() {
	m.mockTokenProvider = nil
	m.app = nil
}

func TestTokenApp_RefreshToken(t *testing.T) {
	mock := &MockTokenApp{}
	mock.Setup(t)
	defer mock.Teardown()

	ctx := context.Background()
	refreshToken := "valid-refresh-token"
	userID := uuid.New().String()

	t.Run("Successful Refresh Token", func(t *testing.T) {
		newAccessToken := "new-access-token"
		newRefreshToken := "new-refresh-token"
		mock.mockTokenProvider.EXPECT().
			VerifyRefreshToken(ctx, refreshToken).
			Return(true, &userID, nil)
		mock.mockTokenProvider.EXPECT().
			GenerateAccessToken(ctx, userID).
			Return(newAccessToken, int64(3600), nil)
		mock.mockTokenProvider.EXPECT().
			GenerateRefreshToken(ctx, userID).
			Return(newRefreshToken, int64(7200), nil)

		newToken, err := mock.app.RefreshToken(ctx, refreshToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if newToken.AccessToken != newAccessToken || newToken.RefreshToken != newRefreshToken {
			t.Fatalf("expected new tokens to match, got %v and %v", newToken.AccessToken, newToken.RefreshToken)
		}
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		mock.mockTokenProvider.EXPECT().
			VerifyRefreshToken(ctx, refreshToken).
			Return(false, nil, errors.New("invalid refresh token"))

		_, err := mock.app.RefreshToken(ctx, refreshToken)
		if err == nil {
			t.Fatal("expected error for invalid refresh token, got nil")
		}
	})

	t.Run("Token Generation Failure", func(t *testing.T) {
		mock.mockTokenProvider.EXPECT().
			VerifyRefreshToken(ctx, refreshToken).
			Return(true, &userID, nil)
		mock.mockTokenProvider.EXPECT().
			GenerateAccessToken(ctx, userID).
			Return("", int64(0), errors.New("token generation failed"))

		_, err := mock.app.RefreshToken(ctx, refreshToken)
		if err == nil {
			t.Fatal("expected error for token generation failure, got nil")
		}
	})
}

func TestTokenApp_VerifyToken(t *testing.T) {
	mock := &MockTokenApp{}
	mock.Setup(t)
	defer mock.Teardown()

	ctx := context.Background()
	accessToken := "valid-access-token"
	userID := uuid.New().String()

	t.Run("Successful Verify Token", func(t *testing.T) {
		mock.mockTokenProvider.EXPECT().
			VerifyAccessToken(ctx, accessToken).
			Return(true, &userID, nil)

		verifyResult, err := mock.app.VerifyToken(ctx, accessToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if verifyResult.UserID != userID || !verifyResult.Valid {
			t.Fatalf("expected valid token for user %s, got %v", userID, verifyResult)
		}
	})

	t.Run("Invalid Access Token", func(t *testing.T) {
		mock.mockTokenProvider.EXPECT().
			VerifyAccessToken(ctx, accessToken).
			Return(false, nil, errors.New("invalid access token"))

		_, err := mock.app.VerifyToken(ctx, accessToken)
		if err == nil {
			t.Fatal("expected error for invalid access token, got nil")
		}
	})
}
