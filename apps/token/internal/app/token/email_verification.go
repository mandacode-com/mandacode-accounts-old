package token

import (
	"errors"

	tokendomain "mandacode.com/accounts/token/internal/domain/service/token"
)

type EmailVerificationTokenApp struct {
	tokenGenerator tokendomain.TokenGenerator
}

// NewEmailVerificationTokenApp creates a new instance of EmailVerificationTokenApp with the provided TokenGenerator.
//
// Parameters:
//   - tokenGenerator: an instance of TokenGenerator used for generating and verifying email verification tokens.
//
// Returns:
//   - EmailVerificationTokenApp: a new instance of EmailVerificationTokenApp.
func NewEmailVerificationTokenApp(
	tokenGenerator tokendomain.TokenGenerator,
) *EmailVerificationTokenApp {
	return &EmailVerificationTokenApp{
		tokenGenerator: tokenGenerator,
	}
}

// GenerateToken generates a signed email verification token for the specified user and email address.
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
func (s *EmailVerificationTokenApp) GenerateToken(userID string, email string, code string) (string, int64, error) {
	claims := map[string]string{
		"sub":   userID,
		"email": email,
		"code":  code,
	}
	return s.tokenGenerator.GenerateToken(
		claims,
	)

}

// VerifyToken verifies the provided email verification token and extracts the user ID, email, and verification code.
//
// Parameters:
//   - token: the JWT email verification token to verify.
//
// Returns:
//   - userID: the user ID extracted from the token claims if valid.
//   - email: the email address extracted from the token claims if valid.
//   - code: the verification code extracted from the token claims if valid.
//   - err: an error if verification fails or if the email or code claims are not found.
func (s *EmailVerificationTokenApp) VerifyToken(token string) (*string, *string, *string, error) {
	claims, err := s.tokenGenerator.VerifyToken(
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
