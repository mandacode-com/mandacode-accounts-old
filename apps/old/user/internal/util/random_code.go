package util

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type RandomCodeGenerator interface {
	// GenerateRandomCode generates a random code of the specified length.
	//
	// Parameters:
	//   - length: The length of the random code to generate.
	GenerateRandomCode(length int) (string, error)

	// GenerateSyncCode generates a synchronous code of the specified length.
	GenerateSyncCode() (string, error)

	// GenerateEmailVerificationCode generates a random email verification code of the specified length.
	GenerateEmailVerificationCode() (string, error)
}

type randomCodeGenerator struct {
	syncCodeLength              int
	emailVerificationCodeLength int
}

func NewRandomCodeGenerator(syncCodeLength, emailVerificationCodeLength int) RandomCodeGenerator {
	return &randomCodeGenerator{
		syncCodeLength:              syncCodeLength,
		emailVerificationCodeLength: emailVerificationCodeLength,
	}
}

func (r *randomCodeGenerator) GenerateRandomCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than zero")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}

	return string(code), nil
}

func (r *randomCodeGenerator) GenerateSyncCode() (string, error) {
	return r.GenerateRandomCode(r.syncCodeLength)
}

func (r *randomCodeGenerator) GenerateEmailVerificationCode() (string, error) {
	return r.GenerateRandomCode(r.emailVerificationCodeLength)
}
