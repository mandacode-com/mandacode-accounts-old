package token

import "errors"

// GenerateRefreshToken generates a signed refresh token for the specified user.
//
// Parameters:
//   - userID: the unique identifier of the user.
//
// Returns:
//   - token: the generated JWT refresh token as a string.
//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time.
//   - err: an error if token generation fails.
func (s *TokenService) GenerateRefreshToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return s.refreshTokenGen.GenerateToken(
		claims,
	)
}

// VerifyRefreshToken verifies the provided refresh token and extracts the user ID.

// Parameters:
//   - token: the JWT refresh token to verify.
//
// Returns:
//   - userID: the user ID extracted from the token claims if valid.
//   - err: an error if verification fails or if the user_id claim is not found.
func (s *TokenService) VerifyRefreshToken(token string) (*string, error) {
	claims, err := s.refreshTokenGen.VerifyToken(
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
