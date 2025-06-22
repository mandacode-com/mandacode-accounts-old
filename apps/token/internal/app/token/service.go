package token

import tokendomain "mandacode.com/accounts/token/internal/domain/token"

// TokenService handles application logic for JWT generation
type TokenService struct {
	accessTokenGen            tokendomain.TokenGenerator
	refreshTokenGen           tokendomain.TokenGenerator
	emailVerificationTokenGen tokendomain.TokenGenerator
}

// NewTokenService constructs the TokenService
func NewTokenService(
	accessTokenGen tokendomain.TokenGenerator,
	refreshTokenGen tokendomain.TokenGenerator,
	emailVerificationTokenGen tokendomain.TokenGenerator,
) *TokenService {
	return &TokenService{
		accessTokenGen:            accessTokenGen,
		refreshTokenGen:           refreshTokenGen,
		emailVerificationTokenGen: emailVerificationTokenGen,
	}
}
