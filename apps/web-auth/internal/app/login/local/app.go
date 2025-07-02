package locallogin

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	logindto "mandacode.com/accounts/web-auth/internal/app/login/dto"
	authdomain "mandacode.com/accounts/web-auth/internal/domain/auth"
	"mandacode.com/accounts/web-auth/internal/util"
)

type localLoginApp struct {
	codeStore     *redis.Client
	authenticator authdomain.Authenticator
	codeTTL       time.Duration
}

// NewLocalLoginApp creates a new instance of LocalLoginApp with the provided Redis client and authenticator.
func NewLocalLoginApp(codeStore *redis.Client, authenticator authdomain.Authenticator, codeTTL time.Duration) LocalLoginApp {
	return &localLoginApp{
		codeStore:     codeStore,
		authenticator: authenticator,
		codeTTL:       codeTTL,
	}
}

func (app *localLoginApp) Login(ctx context.Context, email, password string) (code string, err error) {
	loginToken, err := app.authenticator.LocalLogin(email, password)
	if err != nil {
		return "", err
	}

	// Generate a login code (token) for the authenticated user.
	code, err = util.GenerateSecureRandomCode(16)
	if err != nil {
		return "", err
	}

	targetLoginToken := &logindto.LoginToken{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}
	targetJson, err := json.Marshal(targetLoginToken)
	if err != nil {
		return "", err
	}

	// Store the login token in Redis with the generated code as the key.
	err = app.codeStore.Set(ctx, code, string(targetJson), app.codeTTL).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (app *localLoginApp) VerifyCode(ctx context.Context, code string) (*logindto.LoginToken, error) {
	loginTokenJSON, err := app.codeStore.Get(ctx, code).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("verification code not found or expired")
		}
		return nil, err
	}

	// Deserialize the login token from JSON.
	var loginToken *logindto.LoginToken
	err = json.Unmarshal([]byte(loginTokenJSON), &loginToken)
	if err != nil {
		return nil, err
	}

	// Optionally delete the code after verification to prevent reuse.
	err = app.codeStore.Del(ctx, code).Err()
	if err != nil {
		return nil, err
	}

	return loginToken, nil
}
