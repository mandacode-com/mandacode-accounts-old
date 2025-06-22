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
	expiresIn  time.Duration
}

// NewTokenGenerator creates a new tokenGenerator instance with the provided RSA keys and expiration duration
//
// Parameters:
//   - privateKey: the RSA private key used for signing the token
//   - expiresIn: the duration after which the token will expire
//
// Returns:
//   - *tokenGenerator: a pointer to the newly created tokenGenerator instance
//   - error: an error if any of the parameters are invalid or if key parsing fails
func NewTokenGenerator(
	privateKey *rsa.PrivateKey,
	expiresIn time.Duration) (*tokenGenerator, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	if expiresIn <= 0 {
		return nil, errors.New("expiration duration must be greater than zero")
	}

	return &tokenGenerator{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
		expiresIn:  expiresIn,
	}, nil
}

// NewTokenGeneratorByStr creates a new tokenGenerator using RSA keys provided as PEM formatted strings
//
// Parameters:
//   - privateKeyStr: the PEM formatted RSA private key string
//   - expiresIn: the duration after which the token will expiresIn
//
// Returns:
//   - *tokenGenerator: a pointer to the newly created tokenGenerator instance
//   - error: an error if any of the parameters are invalid or if key parsing fails
func NewTokenGeneratorByStr(
	privateKeyStr string,
	expiresIn time.Duration) (*tokenGenerator, error) {
	privateKey, err := LoadRSAPrivateKeyFromPEM(privateKeyStr)
	if err != nil {
		return nil, err
	}
	return NewTokenGenerator(privateKey, expiresIn)
}

// GenerateToken generates a signed JWT token with the provided claims
//
// Parameters:
//   - claims: a map of additional claims to include in the token
//
// Returns:
//   - token: the generated JWT token as a string
//   - expiresAt: the Unix timestamp (in seconds) indicating the expiration time of the token
//   - error: an error if token generation fails (e.g., if the private key is not initialized or signing fails
func (j *tokenGenerator) GenerateToken(
	claims map[string]string,
) (string, int64, error) {

	if j.privateKey == nil {
		return "", 0, errors.New("private key is not initialized")
	}

	now := time.Now()
	expiresAt := now.Add(j.expiresIn)

	tokenClaims := jwt.MapClaims{
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

// VerifyToken verifies the provided JWT token and extracts the claims
//
// Parameters:
//   - token: the JWT token to verify
//
// Returns:
//   - claims: a map of claims extracted from the token if valid
//   - error: an error if verification fails (e.g., if the public key is not initialized or parsing fails)
func (j *tokenGenerator) VerifyToken(
	token string,
) (map[string]string, error) {

	if j.publicKey == nil {
		return nil, errors.New("public key is not initialized")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
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
