package token

import (
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokengendomain "mandacode.com/accounts/token/internal/domain/infra/token"
	tokendomain "mandacode.com/accounts/token/internal/domain/usecase/token"
)

type tokenUsecase struct {
	accessTokenGenerator            tokengendomain.TokenGenerator
	refreshTokenGenerator           tokengendomain.TokenGenerator
	emailVerificationTokenGenerator tokengendomain.TokenGenerator
}

// GenerateAccessToken implements tokendomain.TokenUsecase.
func (t *tokenUsecase) GenerateAccessToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return t.accessTokenGenerator.GenerateToken(claims)
}

// GenerateEmailVerificationToken implements tokendomain.TokenUsecase.
func (t *tokenUsecase) GenerateEmailVerificationToken(userID string, email string, code string) (string, int64, error) {
	claims := map[string]string{
		"sub":   userID,
		"email": email,
		"code":  code,
	}
	return t.emailVerificationTokenGenerator.GenerateToken(claims)
}

// GenerateRefreshToken implements tokendomain.TokenUsecase.
func (t *tokenUsecase) GenerateRefreshToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return t.refreshTokenGenerator.GenerateToken(claims)
}

// VerifyAccessToken implements tokendomain.TokenUsecase.
func (t *tokenUsecase) VerifyAccessToken(token string) (*string, error) {
	claims, err := t.accessTokenGenerator.VerifyToken(token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify access token")
		return nil, errors.Upgrade(joinedErr, errcode.ErrInvalidToken, "Token Verification Error")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("access token does not contain user ID claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	return &userID, nil
}

// VerifyEmailVerificationToken implements tokendomain.TokenUsecase.
func (t *tokenUsecase) VerifyEmailVerificationToken(token string) (*string, *string, *string, error) {
	claims, err := t.emailVerificationTokenGenerator.VerifyToken(token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify email verification token")
		return nil, nil, nil, errors.Upgrade(joinedErr, errcode.ErrInvalidToken, "Token Verification Error")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, nil, nil, errors.New("email verification token does not contain user ID claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	email, ok := claims["email"]
	if !ok {
		return nil, nil, nil, errors.New("email verification token does not contain email claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	code, ok := claims["code"]
	if !ok {
		return nil, nil, nil, errors.New("email verification token does not contain code claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	return &userID, &email, &code, nil
}

// VerifyRefreshToken implements tokendomain.TokenUsecase.
func (t *tokenUsecase) VerifyRefreshToken(token string) (*string, error) {
	claims, err := t.refreshTokenGenerator.VerifyToken(token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify refresh token")
		return nil, errors.Upgrade(joinedErr, errcode.ErrInvalidToken, "Token Verification Error")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("refresh token does not contain user ID claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	return &userID, nil
}

// NewTokenUsecase creates a new instance of tokenUsecase with the provided TokenGenerators.
//
// Parameters:
//   - accessTokenGenerator: an instance of TokenGenerator used for generating and verifying access tokens.
//   - refreshTokenGenerator: an instance of TokenGenerator used for generating and verifying refresh tokens.
//   - emailVerificationTokenGenerator: an instance of TokenGenerator used for generating and verifying email verification tokens.
func NewTokenUsecase(
	accessTokenGenerator tokengendomain.TokenGenerator,
	refreshTokenGenerator tokengendomain.TokenGenerator,
	emailVerificationTokenGenerator tokengendomain.TokenGenerator,
) tokendomain.TokenUsecase {
	return &tokenUsecase{
		accessTokenGenerator:            accessTokenGenerator,
		refreshTokenGenerator:           refreshTokenGenerator,
		emailVerificationTokenGenerator: emailVerificationTokenGenerator,
	}
}
