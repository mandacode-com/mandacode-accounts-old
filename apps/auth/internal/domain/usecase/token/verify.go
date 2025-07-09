package tokendomain

import "context"

type VerifyUsecase interface {
	// Verify checks the validity of the provided access token.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - token: The access token to verify.
	Verify(ctx context.Context, token string) (valid bool, userID *string, err error)

	// VerifyRefresh checks the validity of the provided refresh token.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - token: The refresh token to verify.
	VerifyRefresh(ctx context.Context, token string) (valid bool, userID *string, err error)
}
