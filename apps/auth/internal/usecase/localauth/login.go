package localauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"

	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	localauthdto "mandacode.com/accounts/auth/internal/usecase/localauth/dto"
)

type LoginUsecase struct {
	authAccount      *dbrepo.AuthAccountRepository
	localAuth        *dbrepo.LocalAuthRepository
	token            *tokenrepo.TokenRepository
	loginCodeManager *coderepo.CodeManager
}

// IssueLoginCode implements localauthdomain.LoginUsecase.
func (l *LoginUsecase) IssueLoginCode(ctx context.Context, input localauthdto.LoginInput) (code string, userID uuid.UUID, err error) {
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
func (l *LoginUsecase) VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error) {
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
func (l *LoginUsecase) Login(ctx context.Context, input localauthdto.LoginInput) (accessToken string, refreshToken string, err error) {
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
func (l *LoginUsecase) issueToken(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
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
	authAccount *dbrepo.AuthAccountRepository,
	localAuth *dbrepo.LocalAuthRepository,
	token *tokenrepo.TokenRepository,
	loginCodeManager *coderepo.CodeManager,
) *LoginUsecase {
	return &LoginUsecase{
		authAccount:      authAccount,
		localAuth:        localAuth,
		token:            token,
		loginCodeManager: loginCodeManager,
	}
}
