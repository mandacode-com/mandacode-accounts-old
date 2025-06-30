package locallogin

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	logindto "mandacode.com/accounts/auth/internal/app/login/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	tokendomain "mandacode.com/accounts/auth/internal/domain/token"
	tokenprovider "mandacode.com/accounts/auth/internal/infra/token"
)

type localLoginApp struct {
	tokenProvider tokendomain.TokenProvider
	localLoginRep repodomain.LocalUserRepository
}

func NewLocalLoginApp(tokenProvider tokendomain.TokenProvider, localLoginRep repodomain.LocalUserRepository) LocalLoginApp {
	return &localLoginApp{
		tokenProvider: tokenProvider,
		localLoginRep: localLoginRep,
	}
}

func (a *localLoginApp) Login(ctx context.Context, email, password string) (*logindto.LoginToken, error) {
	user, err := a.localLoginRep.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, logindto.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, logindto.ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, logindto.ErrUserNotActive
	}
	if !user.IsVerified {
		return nil, logindto.ErrUserNotVerified
	}

	userID := user.ID.String()

	// Generate access and refresh tokens for the LoginLocalUser
	accessToken, _, err := a.tokenProvider.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, tokenprovider.ErrTokenGenerationFailed
	}
	refreshToken, _, err := a.tokenProvider.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return nil, tokenprovider.ErrTokenGenerationFailed
	}

	return &logindto.LoginToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
