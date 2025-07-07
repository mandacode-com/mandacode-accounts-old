package tokenhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"mandacode.com/accounts/auth/internal/app/token"
	tokenv1 "github.com/mandacode-com/accounts-proto/auth/token/v1"
)

type TokenHandler struct {
	tokenv1.UnimplementedTokenServiceServer
	app token.TokenApp
	logger *zap.Logger
}

// NewTokenHandler returns a new token service handler
func NewTokenHandler(app token.TokenApp, logger *zap.Logger) (tokenv1.TokenServiceServer, error) {
	if app == nil {
		return nil, errors.New("TokenApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &TokenHandler{
		app:    app,
		logger: logger,
	}, nil
}

func (h *TokenHandler) RefreshToken(ctx context.Context, req *tokenv1.RefreshTokenRequest) (*tokenv1.RefreshTokenResponse, error) {
	newToken, err := h.app.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Error("failed to refresh token", zap.Error(err), zap.String("refresh_token", req.RefreshToken))
		return nil, err
	}

	return &tokenv1.RefreshTokenResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}, nil
}

func (h *TokenHandler) VerifyToken(ctx context.Context, req *tokenv1.VerifyTokenRequest) (*tokenv1.VerifyTokenResponse, error) {
	verifyResult, err := h.app.VerifyToken(ctx, req.AccessToken)
	if err != nil {
		h.logger.Error("failed to verify token", zap.Error(err), zap.String("access_token", req.AccessToken))
		return nil, err
	}

	return &tokenv1.VerifyTokenResponse{
		UserId: verifyResult.UserID,
		Valid:  verifyResult.Valid,
	}, nil
}
