package token

import (
	"errors"

	tokendomain "mandacode.com/accounts/token/internal/domain/service/token"
)

type AccessTokenApp struct {
	tokenGenerator tokendomain.TokenGenerator
}

// NewAccessTokenApp creates a new instance of AccessTokenApp with the provided TokenGenerator.
//
// Parameters:
//   - tokenGenerator: an instance of TokenGenerator used for generating and verifying access tokens.
//
// Returns:
//   - AccessTokenApp: a new instance of AccessTokenApp.
func NewAccessTokenApp(
	tokenGenerator tokendomain.TokenGenerator,
) *AccessTokenApp {
	return &AccessTokenApp{
		tokenGenerator: tokenGenerator,
	}
}

// GenerateToken generates a signed access token for the specified user.
//
// Parameters:
//   - userID: the unique identifier of the user.
//
// Returns:
//   - token: the generated JWT access token as a string.
//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time.
//   - err: an error if token generation fails.
func (s *AccessTokenApp) GenerateToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return s.tokenGenerator.GenerateToken(
		claims,
	)
}

// VerifyToken verifies the provided access token and extracts the user ID.
//
// Parameters:
//   - token: the JWT access token to verify.
//
// Returns:
//   - userID: the user ID extracted from the token claims if valid.
//   - err: an error if verification fails or if the user_id claim is not found.
func (s *AccessTokenApp) VerifyToken(token string) (*string, error) {
	claims, err := s.tokenGenerator.VerifyToken(
		token,
	)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		return nil, errors.New("no claims found in access token")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("user_id claim not found in access token")
	}

	return &userID, nil
}
