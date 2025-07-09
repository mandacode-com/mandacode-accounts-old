package tokendomain

type TokenUsecase interface {
	// GenerateAccessToken generates an access token for a user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user for whom the access token is generated.
	GenerateAccessToken(userID string) (string, int64, error)

	// VerifyAccessToken verifies the provided access token and extracts the user ID.
	//
	// Parameters:
	//   - token: The JWT access token to verify.
	VerifyAccessToken(token string) (*string, error)

	// GenerateRefreshToken generates a refresh token for a user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user for whom the refresh token is generated.
	GenerateRefreshToken(userID string) (string, int64, error)

	// VerifyRefreshToken verifies the provided refresh token and extracts the user ID.
	//
	// Parameters:
	//   - token: The JWT refresh token to verify.
	VerifyRefreshToken(token string) (*string, error)

	// GenerateEmailVerificationToken generates an email verification token for a user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user for whom the email verification token is generated.
	//   - email: The email address to verify.
	//   - code: The verification code associated with the email.
	GenerateEmailVerificationToken(userID string, email string, code string) (string, int64, error)

	// VerifyEmailVerificationToken verifies the provided email verification token and extracts the user ID, email, and verification code.
	//
	// Parameters:
	//   - token: The JWT email verification token to verify.
	VerifyEmailVerificationToken(token string) (*string, *string, *string, error)
}
