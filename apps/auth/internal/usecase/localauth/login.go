package localauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	codedomain "mandacode.com/accounts/auth/internal/domain/repository/code"
	dbdomain "mandacode.com/accounts/auth/internal/domain/repository/database"
	tokendomain "mandacode.com/accounts/auth/internal/domain/repository/token"
	localauthdomain "mandacode.com/accounts/auth/internal/domain/usecase/localauth"
)

type loginUsecase struct {
	authAccount      dbdomain.AuthAccountRepository
	localAuth        dbdomain.LocalAuthRepository
	token            tokendomain.TokenRepository
	loginCodeManager codedomain.CodeManager
}

// IssueLoginCode implements localauthdomain.LoginUsecase.
func (l *loginUsecase) IssueLoginCode(ctx context.Context, input localauthdomain.LoginInput) (code string, userID uuid.UUID, err error) {
	auth, err := l.localAuth.GetLocalAuthByEmail(ctx, input.Email)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get user by ID")
		return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}
	if !auth.IsVerified {
		return "", uuid.Nil, errors.New("user is not verified", PubUserNotVerified, errcode.ErrUnauthorized)
	}

	account, err := l.authAccount.GetAuthAccountByUserID(ctx, auth.AccountID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get auth account by user ID")
		return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}

	verified, err := l.localAuth.ComparePassword(ctx, account.ID, input.Password)
	if err != nil {
		joinedErr := errors.Join(err, "failed to compare password")
		return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}
	if !verified {
		return "", uuid.Nil, errors.New("invalid password", PubAuthenticationFailed, errcode.ErrUnauthorized)
	}

	code, err = l.loginCodeManager.IssueCode(ctx, account.UserID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to issue login code")
		return "", uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	userID = account.UserID
	return code, userID, nil
}

// VerifyLoginCode implements localauthdomain.LoginUsecase.
func (l *loginUsecase) VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error) {
	valid, err := l.loginCodeManager.ValidateCode(ctx, userID, code)
	if err != nil {
		joinedErr := errors.Join(err, "failed to validate login code")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	if !valid {
		return "", "", errors.New("login code is invalid or expired", PubAuthenticationFailed, errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// Login implements localauthdomain.LoginUsecase.
func (l *loginUsecase) Login(ctx context.Context, input localauthdomain.LoginInput) (accessToken string, refreshToken string, err error) {
	auth, err := l.localAuth.GetLocalAuthByEmail(ctx, input.Email)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get user by email")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}
	if !auth.IsVerified {
		return "", "", errors.New("user is not verified", PubUserNotVerified, errcode.ErrUnauthorized)
	}

	account, err := l.authAccount.GetAuthAccountByUserID(ctx, auth.AccountID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get auth account by user ID")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}

	verified, err := l.localAuth.ComparePassword(ctx, account.ID, input.Password)
	if err != nil {
		joinedErr := errors.Join(err, "failed to compare password")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}
	if !verified {
		return "", "", errors.New("invalid password", PubAuthenticationFailed, errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, account.UserID)
}

// issueToken issues a new access token and refresh token for the user.
func (l *loginUsecase) issueToken(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	accessToken, _, err = l.token.GenerateAccessToken(ctx, userID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate access token")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	refreshToken, _, err = l.token.GenerateRefreshToken(ctx, userID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate refresh token")
		return "", "", errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	return accessToken, refreshToken, nil
}

func NewLoginUsecase(
	authAccount dbdomain.AuthAccountRepository,
	localAuth dbdomain.LocalAuthRepository,
	token tokendomain.TokenRepository,
	loginCodeManager codedomain.CodeManager,
) localauthdomain.LoginUsecase {
	return &loginUsecase{
		authAccount:      authAccount,
		localAuth:        localAuth,
		token:            token,
		loginCodeManager: loginCodeManager,
	}
}
