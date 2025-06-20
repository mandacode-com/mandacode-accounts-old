package token

import "errors"

// GenerateAccessToken generates a signed access token for the specified user.
//
// Parameters:
//   - userID: the unique identifier of the user.
//
// Returns:
//   - token: the generated JWT access token as a string.
//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time.
//   - err: an error if token generation fails.
func (s *TokenService) GenerateAccessToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return s.accessTokenGen.GenerateToken(
		claims,
	)
}

// VerifyAccessToken verifies the provided access token and extracts the user ID.
//
// Parameters:
//   - token: the JWT access token to verify.
//
// Returns:
//   - userID: the user ID extracted from the token claims if valid.
//   - err: an error if verification fails or if the user_id claim is not found.
func (s *TokenService) VerifyAccessToken(token string) (*string, error) {
	claims, err := s.accessTokenGen.VerifyToken(
		token,
	)
	if err != nil {
		// return "", err
		return nil, err
	}
	if claims == nil {
		// return "", nil
		return nil, errors.New("no claims found in access token")
	}

	userID, ok := claims["sub"]
	if !ok {
		// return "", errors.New("user_id claim not found in access token")
		return nil, errors.New("user_id claim not found in access token")
	}

	// return userID, nil
	return &userID, nil
}
