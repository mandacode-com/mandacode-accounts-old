package tokenmgrdomain

import (
	"context"

	tokenmodel "mandacode.com/accounts/web-auth/internal/domain/model/token"
)

type TokenManager interface {
	// VerifyToken verifies the provided token and returns the user ID if valid.
	//
	// Parameters:
	//   - ctx: The context for the verification operation.
	//   - accessToken: The access token to verify.
	//
	// Returns:
	//   - result: Verify Access Token Result
	//   - error: An error if the verification process fails, otherwise nil.
	VerifyAccessToken(ctx context.Context, accessToken string) (result *tokenmodel.VerifyAccessTokenResult, err error)

	// RefreshToken verifies the provided refresh token and returns the user ID if valid.
	//
	// Parameters:
	//	 - ctx: The context for the refresh operation.
	//   - refreshToken: The refresh token to verify.
	//
	// Returns:
	//   - result: Refresh Token Result
	//   - error: An error if the refresh token is invalid or verification fails, otherwise nil.
	RefreshToken(ctx context.Context, refreshToken string) (result *tokenmodel.RefresedToken, err error)
}
