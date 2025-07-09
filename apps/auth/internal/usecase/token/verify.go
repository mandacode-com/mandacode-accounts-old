package token

import (
	"context"

	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokenrepodomain "mandacode.com/accounts/auth/internal/domain/repository/token"
	tokendomain "mandacode.com/accounts/auth/internal/domain/usecase/token"
)

type VerifyUsecase struct {
	token tokenrepodomain.TokenRepository
}

// Verify implements tokendomain.VerifyUsecase.
func (v *VerifyUsecase) Verify(ctx context.Context, token string) (valid bool, userID *string, err error) {
	valid, userID, err = v.token.VerifyAccessToken(ctx, token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify access token")
		return false, nil, errors.Upgrade(joinedErr, "Unauthorized", errcode.ErrUnauthorized)
	}
	if !valid || userID == nil {
		return false, nil, nil // Token is invalid or user ID is not present
	}
	return true, userID, nil
}

// VerifyRefresh implements tokendomain.VerifyUsecase.
func (v *VerifyUsecase) VerifyRefresh(ctx context.Context, token string) (valid bool, userID *string, err error) {
	valid, userID, err = v.token.VerifyRefreshToken(ctx, token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify refresh token")
		return false, nil, errors.Upgrade(joinedErr, "Unauthorized", errcode.ErrUnauthorized)
	}
	if !valid || userID == nil {
		return false, nil, nil // Token is invalid or user ID is not present
	}
	return true, userID, nil
}

// NewVerifyUsecase creates a new instance of VerifyUsecase.
func NewVerifyUsecase(token tokenrepodomain.TokenRepository) tokendomain.VerifyUsecase {
	return &VerifyUsecase{
		token: token,
	}
}
