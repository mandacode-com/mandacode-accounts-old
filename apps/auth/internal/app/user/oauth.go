package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	userdomain "mandacode.com/accounts/auth/internal/domain/service/user"
	"mandacode.com/accounts/auth/internal/util"
	providerv1 "mandacode.com/accounts/auth/proto/common/provider/v1"
)

type OAuthUserApp struct {
	providers        *map[oauthuser.Provider]oauthdomain.OAuthService
	oauthUserService userdomain.OAuthUserService
}

func NewOAuthUserApp(providers *map[oauthuser.Provider]oauthdomain.OAuthService, oauthUserService userdomain.OAuthUserService) *OAuthUserApp {
	return &OAuthUserApp{
		providers:        providers,
		oauthUserService: oauthUserService,
	}
}

func (a *OAuthUserApp) CreateOAuthUser(ctx context.Context, userID string, provider providerv1.OAuthProvider, accessToken string, isActive *bool, isVerified *bool) (*dto.OAuthUser, error) {
	providerEnum, err := util.FromProtoToProvider(provider)
	if err != nil {
		return nil, err
	}

	oauthService, exists := (*a.providers)[providerEnum]
	if !exists {
		return nil, errors.New("unsupported provider")
	}

	userInfo, err := oauthService.GetUserInfo(accessToken)

	if err != nil {
		return nil, err
	}
	if userInfo.ProviderID == "" {
		return nil, errors.New("user info does not contain a valid provider ID")
	}
	if userInfo.Email == "" {
		return nil, errors.New("user info does not contain a valid email")
	}
	if userInfo.Name == "" {
		return nil, errors.New("user info does not contain a valid name")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := a.oauthUserService.CreateOAuthUser(userUUID, providerEnum, userInfo.ProviderID, userInfo.Email, isActive, isVerified)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *OAuthUserApp) DeleteOAuthUser(ctx context.Context, userID string, provider *providerv1.OAuthProvider) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// If provider is nil, delete the user regardless of the provider
	if provider == nil {
		return a.oauthUserService.DeleteOAuthUser(userUUID)
	}

	providerEnum, err := util.FromProtoToProvider(*provider)
	if err != nil {
		return err
	}
	return a.oauthUserService.DeleteOAuthUserByProvider(userUUID, providerEnum)
}

func (a *OAuthUserApp) UpdateOAuthUser(
	ctx context.Context,
	userID string,
	provider *providerv1.OAuthProvider,
	providerID *string,
	email *string,
	isActive *bool,
	isVerified *bool,
) (*dto.OAuthUser, error) {
	var providerEnum *oauthuser.Provider
	if provider != nil {
		providerEnumValue, err := util.FromProtoToProvider(*provider)
		if err != nil {
			return nil, err
		}
		providerEnum = &providerEnumValue
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := a.oauthUserService.UpdateOAuthUser(userUUID, providerEnum, providerID, email, isActive, isVerified)
	if err != nil {
		return nil, err
	}

	return user, nil
}
