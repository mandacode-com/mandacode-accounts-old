package token

import (
	"context"

	tokendto "mandacode.com/accounts/auth/internal/app/token/dto"
)

type TokenApp interface {
	// VerifyToken checks the validity of a given token.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - accessToken: The token to be verified.
	//
	// Returns:
	//   - *tokendto.VerifyTokenResult: The result of the verification containing user ID and validity status.
	//   - error: An error if the verification fails, otherwise nil.
	VerifyToken(ctx context.Context, accessToken string) (*tokendto.VerifyTokenResult, error)

	// RefreshToken generates a new access token using a refresh token.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - refreshToken: The refresh token to generate a new access token.
	//
	// Returns:
	//   - *tokendto.NewToken: The new token containing access and refresh tokens.
	//   - error: An error if the refresh fails, otherwise nil.
	RefreshToken(ctx context.Context, refreshToken string) (*tokendto.NewToken, error)
}
