package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
