package handler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	token "mandacode.com/accounts/token/internal/app"
	tokenv1 "github.com/mandacode-com/accounts-proto/token/v1"
)

type TokenHandler struct {
	tokenv1.UnimplementedTokenServiceServer
	accessTokenApp            *token.AccessTokenApp
	refreshTokenApp           *token.RefreshTokenApp
	emailVerificationTokenApp *token.EmailVerificationTokenApp
	logger                    *zap.Logger
}

func NewTokenHandler(
	accessTokenApp *token.AccessTokenApp,
	refreshTokenApp *token.RefreshTokenApp,
	emailVerificationTokenApp *token.EmailVerificationTokenApp,
	logger *zap.Logger,
) (tokenv1.TokenServiceServer, error) {
	if accessTokenApp == nil || refreshTokenApp == nil || emailVerificationTokenApp == nil {
		return nil, errors.New("token applications cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &TokenHandler{
		accessTokenApp:            accessTokenApp,
		refreshTokenApp:           refreshTokenApp,
		emailVerificationTokenApp: emailVerificationTokenApp,
		logger:                    logger,
	}, nil
}

func (h *TokenHandler) GenerateAccessToken(ctx context.Context, req *tokenv1.GenerateAccessTokenRequest) (*tokenv1.GenerateAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for access token generation", zap.Error(err))
		return nil, err
	}

	token, expiresAt, err := h.accessTokenApp.GenerateToken(req.UserId)
	if err != nil {
		h.logger.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	return &tokenv1.GenerateAccessTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *TokenHandler) VerifyAccessToken(ctx context.Context, req *tokenv1.VerifyAccessTokenRequest) (*tokenv1.VerifyAccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for access token verification", zap.Error(err))
		return nil, err
	}

	userId, err := h.accessTokenApp.VerifyToken(req.Token)
	if err != nil {
		h.logger.Error("failed to verify access token", zap.Error(err))
		return nil, err
	}

	return &tokenv1.VerifyAccessTokenResponse{
		Valid:   true,
		UserId: userId,
	}, nil
}

func (h *TokenHandler) GenerateRefreshToken(ctx context.Context, req *tokenv1.GenerateRefreshTokenRequest) (*tokenv1.GenerateRefreshTokenResponse, error) {
	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for refresh token generation", zap.Error(err))
		return nil, err
	}

	token, expiresAt, err := h.refreshTokenApp.GenerateToken(req.UserId)
	if err != nil {
		h.logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	return &tokenv1.GenerateRefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *TokenHandler) VerifyRefreshToken(ctx context.Context, req *tokenv1.VerifyRefreshTokenRequest) (*tokenv1.VerifyRefreshTokenResponse, error) {
	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for refresh token verification", zap.Error(err))
		return nil, err
	}

	userId, err := h.refreshTokenApp.VerifyToken(req.Token)
	if err != nil {
		h.logger.Error("failed to verify refresh token", zap.Error(err))
		return nil, err
	}

	return &tokenv1.VerifyRefreshTokenResponse{
		Valid:   true,
		UserId: userId,
	}, nil
}

func (h *TokenHandler) GenerateEmailVerificationToken(ctx context.Context, req *tokenv1.GenerateEmailVerificationTokenRequest) (*tokenv1.GenerateEmailVerificationTokenResponse, error) {
	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for email verification token generation", zap.Error(err))
		return nil, err
	}

	token, expiresAt, err := h.emailVerificationTokenApp.GenerateToken(req.UserId, req.Email, req.Code)
	if err != nil {
		h.logger.Error("failed to generate email verification token", zap.Error(err))
		return nil, err
	}

	return &tokenv1.GenerateEmailVerificationTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *TokenHandler) VerifyEmailVerificationToken(ctx context.Context, req *tokenv1.VerifyEmailVerificationTokenRequest) (*tokenv1.VerifyEmailVerificationTokenResponse, error) {
	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for email verification token verification", zap.Error(err))
		return nil, err
	}

	userID, email, code, err := h.emailVerificationTokenApp.VerifyToken(req.Token)
	if err != nil {
		h.logger.Error("failed to verify email verification token", zap.Error(err))
		return nil, err
	}

	return &tokenv1.VerifyEmailVerificationTokenResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
		Code:   code,
	}, nil
}
