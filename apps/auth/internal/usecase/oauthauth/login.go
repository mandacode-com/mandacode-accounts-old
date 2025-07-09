package oauthauth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthauth"

	oauthapidomain "mandacode.com/accounts/auth/internal/domain/infra/oauthapi"
	"mandacode.com/accounts/auth/internal/domain/models"
	dbmodels "mandacode.com/accounts/auth/internal/domain/models/database"
	codedomain "mandacode.com/accounts/auth/internal/domain/repository/code"
	dbdomain "mandacode.com/accounts/auth/internal/domain/repository/database"
	tokendomain "mandacode.com/accounts/auth/internal/domain/repository/token"
	oauthauthdomain "mandacode.com/accounts/auth/internal/domain/usecase/oauthauth"
)

type LoginUsecase struct {
	authAccount      dbdomain.AuthAccountRepository
	oauthAuth        dbdomain.OAuthAuthRepository
	token            tokendomain.TokenRepository
	loginCodeManager codedomain.CodeManager
	oauthApiMap      map[oauthauth.Provider]oauthapidomain.OAuthAPI
}

// GetLoginURL implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) GetLoginURL(ctx context.Context, provider string) (loginURL string, err error) {
	api, ok := l.oauthApiMap[oauthauth.Provider(provider)]
	if ok {
		loginURL = api.GetLoginURL()
		return loginURL, nil
	}
	if _, ok := l.oauthApiMap[oauthauth.Provider(provider)]; !ok {
		return "", errors.New(fmt.Sprintf("unsupported provider: %s", provider), "UnsupportedProvider", errcode.ErrInvalidInput)
	}
	return "", errors.New("unsupported login type", "UnsupportedLoginType", errcode.ErrInvalidInput)
}

// IssueLoginCode implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) IssueLoginCode(ctx context.Context, input oauthauthdomain.LoginInput) (code string, userID uuid.UUID, err error) {
	// Get Access Token
	var accessToken string
	if input.AccessToken == "" && input.Code != "" {
		// When code is provided, exchange it for an access token
		api, ok := l.oauthApiMap[input.Provider]
		if !ok {
			return "", uuid.Nil, errors.New(fmt.Sprintf("unsupported provider: %s", input.Provider), "UnsupportedProvider", errcode.ErrInvalidInput)
		}

		accessToken, err = api.GetAccessToken(input.Code)
		if err != nil {
			joinedErr := errors.Join(err, "failed to get access token from OAuth provider")
			return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, "AuthenticationFailed")
		}
	} else if input.AccessToken != "" {
		accessToken = input.AccessToken
	} else {
		return "", uuid.Nil, errors.New("either access token or code must be provided", "InvalidInput", errcode.ErrInvalidInput)
	}

	// Get User Info
	userInfo, err := l.oauthApiMap[input.Provider].GetUserInfo(accessToken)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get user info from OAuth provider")
		return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, "AuthenticationFailed")
	}
	if userInfo == nil {
		return "", uuid.Nil, errors.New("user info is nil", "InvalidUserInfo", errcode.ErrInvalidInput)
	}

	// Check if user exists in the database
	oauthEntity, err := l.oauthAuth.GetOAuthAuthByProviderID(ctx, input.Provider, userInfo.ProviderID)
	if err != nil {
		if errors.Is(err, errcode.ErrNotFound) { // Create new user if not found
			oauthEntity, err = l.createOAuthAuth(ctx, input.Provider, userInfo)
			if err != nil {
				return "", uuid.Nil, err
			}
		} else {
			joinedErr := errors.Join(err, "failed to get OAuth entity from database")
			return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, "AuthenticationFailed")
		}
	}

	// Generate and store login code
	code, err = l.loginCodeManager.IssueCode(ctx, oauthEntity.ID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to issue login code")
		return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	userID = oauthEntity.ID

	return code, userID, nil
}

// Login implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) Login(ctx context.Context, input oauthauthdomain.LoginInput) (accessToken string, refreshToken string, err error) {
	// Get Access Token
	var accessTokenStr string
	if input.AccessToken == "" && input.Code != "" {
		api, ok := l.oauthApiMap[input.Provider]
		if !ok {
			return "", "", errors.New(fmt.Sprintf("unsupported provider: %s", input.Provider), "UnsupportedProvider", errcode.ErrInvalidInput)
		}

		accessTokenStr, err = api.GetAccessToken(input.Code)
		if err != nil {
			joinedErr := errors.Join(err, "failed to get access token from OAuth provider")
			return "", "", errors.Upgrade(joinedErr, errcode.ErrUnauthorized, "AuthenticationFailed")
		}
	} else if input.AccessToken != "" {
		accessTokenStr = input.AccessToken
	} else {
		return "", "", errors.New("either access token or code must be provided", "InvalidInput", errcode.ErrInvalidInput)
	}

	// Get User Info
	userInfo, err := l.oauthApiMap[input.Provider].GetUserInfo(accessTokenStr)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get user info from OAuth provider")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrUnauthorized, "AuthenticationFailed")
	}
	if userInfo == nil {
		return "", "", errors.New("user info is nil", "InvalidUserInfo", errcode.ErrInvalidInput)
	}

	// Check if user exists in the database
	oauthEntity, err := l.oauthAuth.GetOAuthAuthByProviderID(ctx, input.Provider, userInfo.ProviderID)
	if err != nil {
		if errors.Is(err, errcode.ErrNotFound) { // Create new user if not found
			oauthEntity, err = l.createOAuthAuth(ctx, input.Provider, userInfo)
			if err != nil {
				return "", "", err
			}
		} else {
			joinedErr := errors.Join(err, "failed to get OAuth entity from database")
			return "", "", errors.Upgrade(joinedErr, errcode.ErrUnauthorized, "AuthenticationFailed")
		}
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, oauthEntity.ID)
}

// VerifyLoginCode implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error) {
	// Validate code
	valid, err := l.loginCodeManager.ValidateCode(ctx, userID, code)
	if err != nil {
		joinedErr := errors.Join(err, "failed to validate login code")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	if !valid {
		return "", "", errors.New("login code is invalid or expired", "LoginCodeInvalid", errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// issueToken generates access and refresh tokens for the user.
func (l *LoginUsecase) issueToken(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	// Generate access token
	accessToken, _, err = l.token.GenerateAccessToken(ctx, userID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate access token")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	// Generate refresh token
	refreshToken, _, err = l.token.GenerateRefreshToken(ctx, userID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate refresh token")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	return accessToken, refreshToken, nil
}

func (l *LoginUsecase) createOAuthAuth(ctx context.Context, provider oauthauth.Provider, userInfo *models.OAuthUserInfo) (*ent.OAuthAuth, error) {
	account, err := l.authAccount.CreateAuthAccount(
		ctx,
		&dbmodels.CreateAuthAccountInput{
			UserID:   uuid.New(), // Generate a new UUID for the user
		},
	)
	oauthAuth, err := l.oauthAuth.CreateOAuthAuth(
		ctx,
		&dbmodels.CreateOAuthAuthInput{
			AccountID:  account.ID,
			Provider:   provider,
			ProviderID: userInfo.ProviderID,
			Email:      userInfo.Email,
			IsVerified: userInfo.EmailVerified,
		},
	)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to create OAuthAuth", errcode.ErrInternalFailure)
	}
	return oauthAuth, nil
}

// NewLoginUsecase creates a new instance of LoginUsecase.
func NewLoginUsecase(
	authAccount dbdomain.AuthAccountRepository,
	oauthAuth dbdomain.OAuthAuthRepository,
	token tokendomain.TokenRepository,
	loginCodeManager codedomain.CodeManager,
	oauthApiMap map[oauthauth.Provider]oauthapidomain.OAuthAPI,
) oauthauthdomain.LoginUsecase {
	return &LoginUsecase{
		authAccount:      authAccount,
		oauthAuth:        oauthAuth,
		token:            token,
		loginCodeManager: loginCodeManager,
		oauthApiMap:      oauthApiMap,
	}
}
