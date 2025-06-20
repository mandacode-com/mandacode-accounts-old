package app

import (
	"errors"

	"mandacode.com/accounts/token/internal/domain"
)

// TokenService handles application logic for JWT generation
type TokenService struct {
	accessTokenGen            domain.TokenGenerator
	refreshTokenGen           domain.TokenGenerator
	emailVerificationTokenGen domain.TokenGenerator
}

// NewTokenService constructs the TokenService
func NewTokenService(
	accessTokenGen domain.TokenGenerator,
	refreshTokenGen domain.TokenGenerator,
	emailVerificationTokenGen domain.TokenGenerator,
) *TokenService {
	return &TokenService{
		accessTokenGen:            accessTokenGen,
		refreshTokenGen:           refreshTokenGen,
		emailVerificationTokenGen: emailVerificationTokenGen,
	}
}

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

// GenerateEmailVerificationToken generates a signed email verification token.
//
// Parameters:
//   - userID: the unique identifier of the user for whom the email verification is being generated.
//   - email: the email address to verify.
//   - code: the verification code associated with the email.
//
// Returns:
//   - token: the generated JWT email verification token as a string.
//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time.
//   - err: an error if token generation fails.
func (s *TokenService) GenerateEmailVerificationToken(userID string, email string, code string) (string, int64, error) {
	claims := map[string]string{
		"sub":   userID,
		"email": email,
		"code":  code,
	}
	return s.emailVerificationTokenGen.GenerateToken(
		claims,
	)
}

// VerifyEmailVerificationToken verifies the provided email verification token and extracts the email and code.
//
// Parameters:
//   - token: the JWT email verification token to verify.
//
// Returns:
//   - userID: the user ID extracted from the token claims if valid.
//   - email: the email address extracted from the token claims if valid.
//   - code: the verification code extracted from the token claims if valid.
//   - err: an error if verification fails or if the email or code claims are not found.
func (s *TokenService) VerifyEmailVerificationToken(token string) (*string, *string, *string, error) {
	claims, err := s.emailVerificationTokenGen.VerifyToken(
		token,
	)

	if err != nil {
		return nil, nil, nil, err
	}

	if claims == nil {
		return nil, nil, nil, errors.New("no claims found in email verification token")
	}

	sub, ok := claims["sub"]
	if !ok {
		return nil, nil, nil, errors.New("user_id claim not found in email verification token")
	}

	email, ok := claims["email"]
	if !ok {
		return nil, nil, nil, errors.New("email claim not found in email verification token")
	}

	code, ok := claims["code"]
	if !ok {
		return nil, nil, nil, errors.New("code claim not found in email verification token")
	}

	return &sub, &email, &code, nil
}
