package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"mandacode.com/accounts/token/internal/app/token"
	"mandacode.com/accounts/token/internal/domain/token"
	"mandacode.com/accounts/token/internal/handler/token"
	proto "mandacode.com/accounts/token/proto/token/v1"
)

// mockTokenGenerator mocks the domain.TokenGenerator interface
type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) GenerateToken(claims map[string]string) (string, int64, error) {
	args := m.Called(claims)
	return args.String(0), int64(args.Int(1)), args.Error(2)
}

func (m *mockTokenGenerator) VerifyToken(token string) (map[string]string, error) {
	args := m.Called(token)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(map[string]string), args.Error(1)
}

func newTokenServiceWithMocks(m tokendomain.TokenGenerator) *token.TokenService {
	return token.NewTokenService(m, m, m)
}

func TestGenerateAccessTokenHandler(t *testing.T) {
	mockUserID := "d400cc04-32b6-4fdf-88a7-6246ab43ebb9"
	mockGen := new(mockTokenGenerator)
	mockGen.On("GenerateToken", mock.Anything).Return("tok", 1111, nil)

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	resp, err := h.GenerateAccessToken(context.Background(), &proto.GenerateAccessTokenRequest{UserId: mockUserID})
	require.NoError(t, err)
	require.Equal(t, "tok", resp.Token)
	require.Equal(t, int64(1111), resp.ExpiresAt)
}

func TestVerifyAccessTokenHandler_Error(t *testing.T) {
	mockGen := new(mockTokenGenerator)
	mockGen.On("VerifyToken", "bad").Return(nil, errors.New("fail"))

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	_, err := h.VerifyAccessToken(context.Background(), &proto.VerifyAccessTokenRequest{Token: "bad"})
	require.Error(t, err)
}

func TestVerifyAccessTokenHandler_Success(t *testing.T) {
	mockUserID := "d400cc04-32b6-4fdf-88a7-6246ab43ebb9"
	mockGen := new(mockTokenGenerator)
	mockGen.On("VerifyToken", "token").Return(map[string]string{"sub": mockUserID}, nil)

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	resp, err := h.VerifyAccessToken(context.Background(), &proto.VerifyAccessTokenRequest{Token: "token"})
	require.NoError(t, err)
	require.Equal(t, mockUserID, *resp.UserId)
	require.True(t, resp.Valid)
}

func TestGenerateRefreshTokenHandler(t *testing.T) {
	mockUserID := "f02578f2-b525-436b-8968-c0f7c8bf00d2"
	mockGen := new(mockTokenGenerator)
	mockGen.On("GenerateToken", mock.MatchedBy(func(claims map[string]string) bool {
		return claims["sub"] == mockUserID
	})).Return("rtok", 2222, nil)

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	resp, err := h.GenerateRefreshToken(context.Background(), &proto.GenerateRefreshTokenRequest{UserId: mockUserID})
	require.NoError(t, err)
	require.Equal(t, "rtok", resp.Token)
	require.Equal(t, int64(2222), resp.ExpiresAt)
}

func TestVerifyRefreshTokenHandler_Error(t *testing.T) {
	mockGen := new(mockTokenGenerator)
	mockGen.On("VerifyToken", "bad").Return(nil, errors.New("fail"))

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	_, err := h.VerifyRefreshToken(context.Background(), &proto.VerifyRefreshTokenRequest{Token: "bad"})
	require.Error(t, err)
}

func TestVerifyRefreshTokenHandler_Success(t *testing.T) {
	mockGen := new(mockTokenGenerator)
	mockGen.On("VerifyToken", "token").Return(map[string]string{"sub": "uid"}, nil)

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	resp, err := h.VerifyRefreshToken(context.Background(), &proto.VerifyRefreshTokenRequest{Token: "token"})
	require.NoError(t, err)
	require.Equal(t, "uid", *resp.UserId)
	require.True(t, resp.Valid)
}

func TestGenerateEmailVerificationTokenHandler(t *testing.T) {
	mockUserID := "db75a077-6945-48c7-b2cb-7d0f74040ff0"
	mockEmail := "test@test.com"
	mockCode := "code"

	mockGen := new(mockTokenGenerator)
	mockGen.On("GenerateToken", mock.Anything).Return("etok", 3333, nil)

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	resp, err := h.GenerateEmailVerificationToken(context.Background(), &proto.GenerateEmailVerificationTokenRequest{
		UserId: mockUserID, Email: mockEmail, Code: mockCode,
	})
	require.NoError(t, err)
	require.Equal(t, "etok", resp.Token)
	require.Equal(t, int64(3333), resp.ExpiresAt)
}

func TestVerifyEmailVerificationTokenHandler_Error(t *testing.T) {
	mockGen := new(mockTokenGenerator)
	mockGen.On("VerifyToken", "bad").Return(nil, errors.New("fail"))

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	_, err := h.VerifyEmailVerificationToken(context.Background(), &proto.VerifyEmailVerificationTokenRequest{Token: "bad"})
	require.Error(t, err)
}

func TestVerifyEmailVerificationTokenHandler_Success(t *testing.T) {
	mockUserID := "db75a077-6945-48c7-b2cb-7d0f74040ff0"
	mockEmail := "test@test.com"
	mockCode := "code"

	mockGen := new(mockTokenGenerator)
	mockGen.On("VerifyToken", "tok").Return(map[string]string{
		"sub":   mockUserID,
		"email": mockEmail,
		"code":  mockCode,
	}, nil)

	h := tokenhandler.NewTokenHandler(newTokenServiceWithMocks(mockGen), zap.NewNop())
	resp, err := h.VerifyEmailVerificationToken(context.Background(), &proto.VerifyEmailVerificationTokenRequest{Token: "tok"})
	require.NoError(t, err)
	require.True(t, resp.Valid)
	require.Equal(t, mockUserID, *resp.UserId)
	require.Equal(t, mockEmail, *resp.Email)
	require.Equal(t, mockCode, *resp.Code)
}
