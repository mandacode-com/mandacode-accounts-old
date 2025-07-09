package oauthauthdomain

import (
	"context"

	"github.com/google/uuid"
)

type LoginUsecase interface {
	// Login authenticates a user with the provided OAuth access token.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//   - accessToken: The OAuth access token to authenticate the user.
	Login(ctx context.Context, input LoginInput) (accessToken string, refreshToken string, err error)

	// IssueLoginCode issues a login code for the user identified by the email.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - input: The input containing the email of the user to issue a login code for.
	IssueLoginCode(ctx context.Context, input LoginInput) (code string, userID uuid.UUID, err error)

	// VerifyLoginCode verifies the provided login code for the user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to verify.
	//   - code: The login code to be verified.
	VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error)

	// GetLoginURL returns the URL for the OAuth login flow.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//   - loginType: The type of login URL to retrieve (e.g., authorization code, implicit).
	GetLoginURL(ctx context.Context, provider string) (loginURL string, err error)
}
