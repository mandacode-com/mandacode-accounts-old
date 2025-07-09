package token

import (
	"context"

	tokenmodel "mandacode.com/accounts/web-auth/internal/domain/model/token"
)

type TokenApp interface {
	// VerifyAccessToken verifies the provided access token and returns the user ID if valid.
	//
	// Parameters:
	//   - ctx: The context for the verification operation.
	//   - accessToken: The access token to verify.
	//
	// Returns:
	//   - result: A pointer to a VerifyAccessTokenResult containing the verification result.
	//   - err: An error if the verification process fails, otherwise nil.
	VerifyAccessToken(ctx context.Context, accessToken string) (result *tokenmodel.VerifyAccessTokenResult, err error)

	// RefreshToken verifies the provided refresh token and returns new access and refresh tokens if valid.
	//
	// Parameters:
	//   - ctx: The context for the refresh operation.
	//   - refreshToken: The refresh token to verify.
	//
	// Returns:
	//   - result: A pointer to a RefresedToken containing the new access and refresh tokens.
	//   - err: An error if the refresh token is invalid or verification fails, otherwise nil.
	RefreshToken(ctx context.Context, refreshToken string) (result *tokenmodel.RefresedToken, err error)
}
