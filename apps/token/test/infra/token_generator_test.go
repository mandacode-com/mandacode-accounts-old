package infra_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"mandacode.com/accounts/token/internal/infra"
)

func generateTestKeys(t *testing.T) (string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	return string(privPEM)
}

func TestGenerateAndVerifyToken(t *testing.T) {
	priv := generateTestKeys(t)
	gen, err := infra.NewTokenGeneratorByStr(priv, time.Second)
	require.NoError(t, err)

	claims := map[string]string{"sub": "user1", "role": "admin"}
	token, exp, err := gen.GenerateToken(claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.True(t, exp > time.Now().Unix())

	verified, err := gen.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, claims["sub"], verified["sub"])
	require.Equal(t, claims["role"], verified["role"])
}

func TestVerifyInvalidToken(t *testing.T) {
	priv := generateTestKeys(t)
	gen, err := infra.NewTokenGeneratorByStr(priv, time.Minute)
	require.NoError(t, err)

	_, err = gen.VerifyToken("invalid.token.here")
	require.Error(t, err)
}

func TestExpiredToken(t *testing.T) {
	priv := generateTestKeys(t)
	gen, err := infra.NewTokenGeneratorByStr(priv, 1*time.Millisecond)
	require.NoError(t, err)

	token, _, err := gen.GenerateToken(map[string]string{"sub": "user1"})
	require.NoError(t, err)

	// wait for token to expire
	time.Sleep(5 * time.Millisecond)

	_, err = gen.VerifyToken(token)
	require.Error(t, err)
}

func TestNewTokenGeneratorFailures(t *testing.T) {
	// Test with invalid private key
	_, err := infra.NewTokenGeneratorByStr("invalid-key", time.Second)
	require.Error(t, err)

	// Test with empty private key
	_, err = infra.NewTokenGeneratorByStr("", time.Second)
	require.Error(t, err)

	// Test with zero expiration time
	priv := generateTestKeys(t)
	_, err = infra.NewTokenGeneratorByStr(priv, 0)
	require.Error(t, err)
}
