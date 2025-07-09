package tokendomain

import "context"

type RefreshUsecase interface {
	// Refresh generates a new access token and refresh token using the provided refresh token.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - refreshToken: The refresh token to use for generating new tokens.
	Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error)
}
