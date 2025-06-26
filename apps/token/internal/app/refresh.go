package token

import (
	"errors"

	svcdomain "mandacode.com/accounts/token/internal/domain/service"
)

type RefreshTokenApp struct {
	tokenGenerator svcdomain.TokenGenerator
}

// NewRefreshTokenApp creates a new instance of RefreshTokenApp with the provided TokenGenerator.
//
// Parameters:
//   - tokenGenerator: an instance of TokenGenerator used for generating and verifying refresh tokens.
//
// Returns:
//   - RefreshTokenApp: a new instance of RefreshTokenApp.
func NewRefreshTokenApp(
	tokenGenerator svcdomain.TokenGenerator,
) *RefreshTokenApp {
	return &RefreshTokenApp{
		tokenGenerator: tokenGenerator,
	}
}

// GenerateToken generates a signed refresh token for the specified user.
//
// Parameters:
//   - userID: the unique identifier of the user.
//
// Returns:
//   - token: the generated JWT refresh token as a string.
//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time.
//   - err: an error if token generation fails.
func (s *RefreshTokenApp) GenerateToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return s.tokenGenerator.GenerateToken(
		claims,
	)
}

// VerifyToken verifies the provided refresh token and extracts the user ID.
// Parameters:
//   - token: the JWT refresh token to verify.
//
// Returns:
//   - userID: the user ID extracted from the token claims if valid.
//   - err: an error if verification fails or if the user_id claim is not found.
func (s *RefreshTokenApp) VerifyToken(token string) (*string, error) {
	claims, err := s.tokenGenerator.VerifyToken(
		token,
	)

	if err != nil {
		// return "", err
		return nil, err
	}

	if claims == nil {
		return nil, errors.New("no claims found in refresh token")
	}

	userID, ok := claims["sub"]
	if !ok {
		// return "", errors.New("user_id claim not found in refresh token")
		return nil, errors.New("user_id claim not found in refresh token")
	}

	// return userID, nil
	return &userID, nil
}
