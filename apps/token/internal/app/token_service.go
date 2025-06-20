package app

import "mandacode.com/accounts/token/internal/domain"

// TokenService handles application logic for JWT generation
type TokenService struct {
	// tokenGen domain.TokenGenerator
	accessTokenGen            domain.TokenGenerator
	refreshTokenGen           domain.TokenGenerator
	emailVerificationTokenGen domain.TokenGenerator
}

// NewTokenService constructs the TokenService
func NewTokenService(
	accessTokenGen domain.TokenGenerator,
	refreshTokenGen domain.TokenGenerator,
	emailVerificationTokenGen domain.TokenGenerator,
) *TokenService {
	return &TokenService{
		accessTokenGen:            accessTokenGen,
		refreshTokenGen:           refreshTokenGen,
		emailVerificationTokenGen: emailVerificationTokenGen,
	}
}

func (s *TokenService) GenerateAccessToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"user_id": userID,
	}
	return s.accessTokenGen.GenerateToken(
		claims,
	)
}

func (s *TokenService) VerifyAccessToken(token string) (map[string]string, error) {
	claims, err := s.accessTokenGen.VerifyToken(
		token,
	)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		return nil, nil
	}
	return claims, nil
}

func (s *TokenService) GenerateRefreshToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"user_id": userID,
	}
	return s.refreshTokenGen.GenerateToken(
		claims,
	)
}

func (s *TokenService) VerifyRefreshToken(token string) (map[string]string, error) {
	claims, err := s.refreshTokenGen.VerifyToken(
		token,
	)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		return nil, nil
	}
	return claims, nil
}

func (s *TokenService) GenerateEmailVerificationToken(email string, code string) (string, int64, error) {
	claims := map[string]string{
		"email": email,
		"code":  code,
	}
	return s.emailVerificationTokenGen.GenerateToken(
		claims,
	)
}

func (s *TokenService) VerifyEmailVerificationToken(token string) (map[string]string, error) {
	claims, err := s.emailVerificationTokenGen.VerifyToken(
		token,
	)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		return nil, nil
	}
	return claims, nil
}
