package tokenhandler

import (
	"context"

	"go.uber.org/zap"
	"mandacode.com/accounts/token/internal/app/token"
	proto "mandacode.com/accounts/token/proto/token/v1"
)

// tokenHandler implements the JWTServiceServer gRPC interface
type tokenHandler struct {
	proto.UnimplementedTokenServiceServer
	tokenService *token.TokenService
	logger       *zap.Logger
}

// NewTokenHandler returns a gRPC handler ith dependencies injected
func NewTokenHandler(tokenService *token.TokenService, logger *zap.Logger) proto.TokenServiceServer {
	return &tokenHandler{
		tokenService: tokenService,
		logger:       logger,
	}
}

// GenerateAccessToken handles the gRPC call for generating access tokens
func (h *tokenHandler) GenerateAccessToken(ctx context.Context, req *proto.GenerateAccessTokenRequest) (*proto.GenerateAccessTokenResponse, error) {

	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for access token generation", zap.Error(err))
		return nil, err
	}

	token, expiresAt, err := h.tokenService.GenerateAccessToken(req.UserId)
	if err != nil {
		h.logger.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	return &proto.GenerateAccessTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyAccessToken handles the gRPC call for verifying access tokens
func (h *tokenHandler) VerifyAccessToken(ctx context.Context, req *proto.VerifyAccessTokenRequest) (*proto.VerifyAccessTokenResponse, error) {

	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for access token verification", zap.Error(err))
		return nil, err
	}

	userId, err := h.tokenService.VerifyAccessToken(req.Token)

	if err != nil {
		h.logger.Error("failed to verify access token", zap.Error(err))
		return nil, err
	}

	return &proto.VerifyAccessTokenResponse{
		Valid:  true,
		UserId: userId,
	}, nil
}

// GenerateRefreshToken handles the gRPC call for generating refresh tokens
func (h *tokenHandler) GenerateRefreshToken(ctx context.Context, req *proto.GenerateRefreshTokenRequest) (*proto.GenerateRefreshTokenResponse, error) {

	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for refresh token generation", zap.Error(err))
		return nil, err
	}

	token, expiresAt, err := h.tokenService.GenerateRefreshToken(req.UserId)
	if err != nil {
		h.logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	return &proto.GenerateRefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyRefreshToken handles the gRPC call for verifying refresh tokens
func (h *tokenHandler) VerifyRefreshToken(ctx context.Context, req *proto.VerifyRefreshTokenRequest) (*proto.VerifyRefreshTokenResponse, error) {

	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for refresh token verification", zap.Error(err))
		return nil, err
	}

	userId, err := h.tokenService.VerifyRefreshToken(req.Token)

	if err != nil {
		h.logger.Error("failed to verify refresh token", zap.Error(err))
		return nil, err
	}

	return &proto.VerifyRefreshTokenResponse{
		Valid:  true,
		UserId: userId,
	}, nil
}

// GenerateEmailVerificationToken handles the gRPC call for generating email verification tokens
func (h *tokenHandler) GenerateEmailVerificationToken(ctx context.Context, req *proto.GenerateEmailVerificationTokenRequest) (*proto.GenerateEmailVerificationTokenResponse, error) {

	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for email verification token generation", zap.Error(err))
		return nil, err
	}

	token, expiresAt, err := h.tokenService.GenerateEmailVerificationToken(req.UserId, req.Email, req.Code)
	if err != nil {
		h.logger.Error("failed to generate email verification token", zap.Error(err))
		return nil, err
	}

	return &proto.GenerateEmailVerificationTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyEmailVerificationToken handles the gRPC call for verifying email verification tokens
func (h *tokenHandler) VerifyEmailVerificationToken(ctx context.Context, req *proto.VerifyEmailVerificationTokenRequest) (*proto.VerifyEmailVerificationTokenResponse, error) {

	if err := req.Validate(); err != nil {
		h.logger.Error("invalid request for email verification token verification", zap.Error(err))
		return nil, err
	}

	userID, email, code, err := h.tokenService.VerifyEmailVerificationToken(req.Token)

	if err != nil {
		h.logger.Error("failed to verify email verification token", zap.Error(err))
		return nil, err
	}

	return &proto.VerifyEmailVerificationTokenResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
		Code:   code,
	}, nil
}
