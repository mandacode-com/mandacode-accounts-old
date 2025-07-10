package locallogin

import (
	"context"
	logindto "mandacode.com/accounts/mobile-auth/internal/app/login/dto"
	authdomain "mandacode.com/accounts/mobile-auth/internal/domain/auth"
)

type localLoginApp struct {
	authenticator authdomain.Authenticator
}

// Login implements LocalLoginApp.
func (l *localLoginApp) Login(ctx context.Context, email string, password string) (*logindto.LoginToken, error) {
	loginToken, err := l.authenticator.LocalLogin(email, password)
	if err != nil {
		return nil, err
	}

	return &logindto.LoginToken{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}, nil
}

// NewLocalLoginApp creates a new instance of LocalLoginApp with the provided authenticator.
func NewLocalLoginApp(authenticator authdomain.Authenticator) LocalLoginApp {
	return &localLoginApp{
		authenticator: authenticator,
	}
}
