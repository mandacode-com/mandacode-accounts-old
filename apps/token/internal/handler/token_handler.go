package handler

import (
	"context"
	"log"

	"mandacode.com/accounts/token/internal/app"
	proto "mandacode.com/accounts/token/proto/token/v1"
)

// tokenHandler implements the JWTServiceServer gRPC interface
type tokenHandler struct {
	proto.UnimplementedTokenServiceServer
	tokenService *app.TokenService
}

// NewTokenHandler returns a gRPC handler with dependencies injected
func NewTokenHandler(tokenService *app.TokenService) proto.TokenServiceServer {
	return &tokenHandler{
		tokenService: tokenService,
	}
}

// GenerateAccessToken handles the gRPC call for generating access tokens
func (h *tokenHandler) GenerateAccessToken(ctx context.Context, req *proto.GenerateAccessTokenRequest) (*proto.GenerateAccessTokenResponse, error) {

	if err := req.Validate(); err != nil {
		log.Printf("invalid request for access token generation: %v", err)
		return nil, err
	}

	token, expiresAt, err := h.tokenService.GenerateAccessToken(req.UserId)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		return nil, err
	}

	return &proto.GenerateAccessTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateRefreshToken handles the gRPC call for generating refresh tokens
func (h *tokenHandler) GenerateRefreshToken(ctx context.Context, req *proto.GenerateRefreshTokenRequest) (*proto.GenerateRefreshTokenResponse, error) {

	if err := req.Validate(); err != nil {
		log.Printf("invalid request for refresh token generation: %v", err)
		return nil, err
	}

	token, expiresAt, err := h.tokenService.GenerateRefreshToken(req.UserId)
	if err != nil {
		log.Printf("failed to generate refresh token: %v", err)
		return nil, err
	}

	return &proto.GenerateRefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateEmailVerificationToken handles the gRPC call for generating email verification tokens
func (h *tokenHandler) GenerateEmailVerificationToken(ctx context.Context, req *proto.GenerateEmailVerificationTokenRequest) (*proto.GenerateEmailVerificationTokenResponse, error) {

	if err := req.Validate(); err != nil {
		log.Printf("invalid request for email verification token generation: %v", err)
		return nil, err
	}

	token, expiresAt, err := h.tokenService.GenerateEmailVerificationToken(req.Email, req.Code)
	if err != nil {
		log.Printf("failed to generate email verification token: %v", err)
		return nil, err
	}

	return &proto.GenerateEmailVerificationTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
