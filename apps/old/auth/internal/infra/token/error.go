package tokenprovider

import "errors"

var (
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrTokenValidationFailed = errors.New("token validation failed")
	ErrTokenVerificationFailed = errors.New("token verification failed")
)
