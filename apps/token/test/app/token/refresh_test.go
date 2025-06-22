package token_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"mandacode.com/accounts/token/internal/app/token"
)

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	mockUserID := "f7ba2494-bdad-415f-8750-257f510baecb" // UUID for testing
	mockGen := new(MockTokenGenerator)
	mockGen.On("GenerateToken", mock.MatchedBy(func(claims map[string]string) bool {
		return claims["sub"] == mockUserID
	})).Return("refreshtoken123", int64(1234567890), nil)

	svc := token.NewTokenService(mockGen, mockGen, mockGen)
	token, exp, err := svc.GenerateRefreshToken(mockUserID)
	require.NoError(t, err)
	require.Equal(t, "refreshtoken123", token)
	require.Equal(t, int64(1234567890), exp)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyRefreshToken_Success(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "validrefreshtoken").Return(map[string]string{"sub": "user123"}, nil)

	svc := token.NewTokenService(mockGen, mockGen, mockGen)
	userID, err := svc.VerifyRefreshToken("validrefreshtoken")
	require.NoError(t, err)
	require.Equal(t, "user123", *userID)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyRefreshToken_MissingClaim(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "invalidrefreshtoken").Return(map[string]string{}, nil)

	svc := token.NewTokenService(mockGen, mockGen, mockGen)
	userID, err := svc.VerifyRefreshToken("invalidrefreshtoken")
	require.Error(t, err)
	require.Nil(t, userID)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyRefreshToken_Error(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "brokenrefreshtoken").Return(nil, errors.New("invalid token"))

	svc := token.NewTokenService(mockGen, mockGen, mockGen)
	userID, err := svc.VerifyRefreshToken("brokenrefreshtoken")
	require.Error(t, err)
	require.Nil(t, userID)

	mockGen.AssertExpectations(t)
}
