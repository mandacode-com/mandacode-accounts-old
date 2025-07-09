package token

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokenrepodomain "mandacode.com/accounts/auth/internal/domain/repository/token"
	tokendomain "mandacode.com/accounts/auth/internal/domain/usecase/token"
)

type RefreshUsecase struct {
	token tokenrepodomain.TokenRepository
}

// Refresh implements tokendomain.RefreshUsecase.
func (r *RefreshUsecase) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	// Validate the refresh token
	valid, userID, err := r.token.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.New("failed to verify refresh token", "Unauthorized", errcode.ErrUnauthorized)
	}
	if !valid || userID == nil {
		return "", "", errors.New("invalid refresh token", "Unauthorized", errcode.ErrUnauthorized)
	}

	// Generate new access and refresh tokens
	userUID, err := uuid.Parse(*userID)
	if err != nil {
		return "", "", errors.New("invalid user ID in refresh token", "InvalidUserID", errcode.ErrUnauthorized)
	}
	newAccessToken, _, err = r.token.GenerateAccessToken(ctx, userUID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate new access token")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}
	newRefreshToken, _, err = r.token.GenerateRefreshToken(ctx, userUID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate new refresh token")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}

	return newAccessToken, newRefreshToken, nil
}

func NewRefreshUsecase(token tokenrepodomain.TokenRepository) tokendomain.RefreshUsecase {
	return &RefreshUsecase{
		token: token,
	}
}
