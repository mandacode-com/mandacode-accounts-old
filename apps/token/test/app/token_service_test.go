package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"mandacode.com/accounts/token/internal/app"
)

// MockTokenGenerator implements domain.TokenGenerator for unit testing
type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(claims map[string]string) (string, int64, error) {
	args := m.Called(claims)
	if args.Get(0) == nil {
		return "", 0, args.Error(1)
	}
	return args.String(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockTokenGenerator) VerifyToken(token string) (map[string]string, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestTokenService_GenerateAccessToken(t *testing.T) {
	mockUserID := "f7ba2494-bdad-415f-8750-257f510baecb" // UUID for testing
	mockGen := new(MockTokenGenerator)
	mockGen.On("GenerateToken", mock.MatchedBy(func(claims map[string]string) bool {
		return claims["sub"] == mockUserID
	})).Return("token123", int64(1234567890), nil)

	svc := app.NewTokenService(mockGen, mockGen, mockGen)
	token, exp, err := svc.GenerateAccessToken(mockUserID)
	require.NoError(t, err)
	require.Equal(t, "token123", token)
	require.Equal(t, int64(1234567890), exp)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyAccessToken_Success(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "validtoken").Return(map[string]string{"sub": "user123"}, nil)

	svc := app.NewTokenService(mockGen, mockGen, mockGen)
	userID, err := svc.VerifyAccessToken("validtoken")
	require.NoError(t, err)
	require.Equal(t, "user123", *userID)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyAccessToken_MissingClaim(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "invalidtoken").Return(map[string]string{}, nil)

	svc := app.NewTokenService(mockGen, mockGen, mockGen)
	userID, err := svc.VerifyAccessToken("invalidtoken")
	require.Error(t, err)
	require.Nil(t, userID)

	mockGen.AssertExpectations(t)
}

func TestTokenService_GenerateEmailVerificationToken(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("GenerateToken", mock.MatchedBy(func(claims map[string]string) bool {
		return claims["sub"] == "u1" && claims["email"] == "a@b.c" && claims["code"] == "xyz"
	})).Return("emailtoken", int64(9999), nil)

	svc := app.NewTokenService(mockGen, mockGen, mockGen)
	token, exp, err := svc.GenerateEmailVerificationToken("u1", "a@b.c", "xyz")
	require.NoError(t, err)
	require.Equal(t, "emailtoken", token)
	require.Equal(t, int64(9999), exp)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyEmailVerificationToken_Failure(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "broken").Return(nil, errors.New("invalid token"))

	svc := app.NewTokenService(mockGen, mockGen, mockGen)
	uid, email, code, err := svc.VerifyEmailVerificationToken("broken")
	require.Error(t, err)
	require.Nil(t, uid)
	require.Nil(t, email)
	require.Nil(t, code)

	mockGen.AssertExpectations(t)
}
