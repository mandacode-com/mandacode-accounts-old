package token_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"mandacode.com/accounts/token/internal/app/token"
)

func TestTokenService_GenerateEmailVerificationToken(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("GenerateToken", mock.MatchedBy(func(claims map[string]string) bool {
		return claims["sub"] == "u1" && claims["email"] == "a@b.c" && claims["code"] == "xyz"
	})).Return("emailtoken", int64(9999), nil)

	app := token.NewEmailVerificationTokenApp(mockGen)
	token, exp, err := app.GenerateToken("u1", "a@b.c", "xyz")
	require.NoError(t, err)
	require.Equal(t, "emailtoken", token)
	require.Equal(t, int64(9999), exp)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyEmailVerificationToken_Success(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "validtoken").Return(map[string]string{
		"sub":   "u1",
		"email": "a@b.c",
		"code":  "xyz",
	}, nil)

	app := token.NewEmailVerificationTokenApp(mockGen)
	uid, email, code, err := app.VerifyToken("validtoken")
	require.NoError(t, err)
	require.Equal(t, "u1", *uid)
	require.Equal(t, "a@b.c", *email)
	require.Equal(t, "xyz", *code)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyEmailVerificationToken_MissingClaims(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "invalidtoken").Return(map[string]string{}, nil)

	app := token.NewEmailVerificationTokenApp(mockGen)
	uid, email, code, err := app.VerifyToken("invalidtoken")
	require.Error(t, err)
	require.Nil(t, uid)
	require.Nil(t, email)
	require.Nil(t, code)

	mockGen.AssertExpectations(t)
}

func TestTokenService_VerifyEmailVerificationToken_Error(t *testing.T) {
	mockGen := new(MockTokenGenerator)
	mockGen.On("VerifyToken", "broken").Return(nil, errors.New("invalid token"))

	app := token.NewEmailVerificationTokenApp(mockGen)
	uid, email, code, err := app.VerifyToken("broken")
	require.Error(t, err)
	require.Nil(t, uid)
	require.Nil(t, email)
	require.Nil(t, code)

	mockGen.AssertExpectations(t)
}
