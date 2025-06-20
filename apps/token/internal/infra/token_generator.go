package infra

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// jwtGenerator is the concrete implementation of TokenGenerator
type tokenGenerator struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	sub        string
	expiresIn  time.Duration
}

func LoadRSAPublicKeyFromPEM(keyStr string) (*rsa.PublicKey, error) {
	if keyStr == "" {
		return nil, errors.New("public key string cannot be empty")
	}
	// Load the public key from PEM format
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
	if err != nil {
		return nil, errors.New("failed to parse public key: " + err.Error())
	}
	return publicKey, nil
}

func LoadRSAPrivateKeyFromPEM(keyStr string) (*rsa.PrivateKey, error) {
	if keyStr == "" {
		return nil, errors.New("private key string cannot be empty")
	}
	// Load the private key from PEM format
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyStr))
	if err != nil {
		return nil, errors.New("failed to parse private key: " + err.Error())
	}
	return privateKey, nil
}

// NewTokenGenerator creates a new tokenGenerator with the provided RSA private key, subject, and expiration duration
func NewTokenGenerator(
	publicKey *rsa.PublicKey,
	privateKey *rsa.PrivateKey,
	sub string,
	expiresIn time.Duration) (*tokenGenerator, error) {
	if publicKey == nil {
		return nil, errors.New("public key cannot be nil")
	}
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	if sub == "" {
		return nil, errors.New("subject (sub) cannot be empty")
	}
	if expiresIn <= 0 {
		return nil, errors.New("expiration duration must be greater than zero")
	}

	return &tokenGenerator{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
		sub:        sub,
		expiresIn:  expiresIn,
	}, nil
}

func NewTokenGeneratorByStr(
	publicKeyStr string,
	privateKeyStr string,
	sub string,
	expiresIn time.Duration) (*tokenGenerator, error) {
	publicKey, err := LoadRSAPublicKeyFromPEM(publicKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := LoadRSAPrivateKeyFromPEM(privateKeyStr)
	if err != nil {
		return nil, err
	}
	return NewTokenGenerator(publicKey, privateKey, sub, expiresIn)
}

// GenerateToken creates a signed JWT token using RSA private key
func (j *tokenGenerator) GenerateToken(
	claims map[string]string,
) (string, int64, error) {

	if j.privateKey == nil {
		return "", 0, errors.New("private key is not initialized")
	}

	now := time.Now()
	expiresAt := now.Add(j.expiresIn)

	tokenClaims := jwt.MapClaims{
		"sub": j.sub,
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	for key, value := range claims {
		tokenClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresAt.Unix(), nil
}

func (j *tokenGenerator) VerifyToken(
	token string,
) (map[string]string, error) {

	if j.publicKey == nil {
		return nil, errors.New("public key is not initialized")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		result := make(map[string]string)
		for key, value := range claims {
			if strValue, ok := value.(string); ok {
				result[key] = strValue
			}
		}
		return result, nil
	}

	return nil, errors.New("invalid token")
}
