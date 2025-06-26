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

func (a *OAuthUserApp) CreateUser(ctx context.Context, userID string, provider providerv1.OAuthProvider, accessToken string, isActive *bool, isVerified *bool) (*dto.OAuthUser, error) {
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
	return a.oauthUserService.CreateUser(userUUID, providerEnum, userInfo.ProviderID, userInfo.Email, isActive, isVerified)
}

func (a *OAuthUserApp) GetUser(ctx context.Context, userID string, provider providerv1.OAuthProvider) (*dto.OAuthUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	providerEnum, err := util.FromProtoToProvider(provider)
	if err != nil {
		return nil, err
	}

	return a.oauthUserService.GetUserByUserID(userUUID, providerEnum)
}

func (a *OAuthUserApp) DeleteUser(ctx context.Context, userID string, provider *providerv1.OAuthProvider) (*dto.OAuthDeletedUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// If provider is nil, delete the user regardless of the provider
	if provider == nil {
		return a.oauthUserService.DeleteUser(userUUID)
	}

	providerEnum, err := util.FromProtoToProvider(*provider)
	if err != nil {
		return nil, err
	}

	return a.oauthUserService.DeleteUserByProvider(userUUID, providerEnum)
}

func (a *OAuthUserApp) SyncUser(ctx context.Context, userID string, provider providerv1.OAuthProvider, accessToken string) (*dto.OAuthUser, error) {
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

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return a.oauthUserService.UpdateUserBase(userUUID, providerEnum, userInfo.ProviderID, userInfo.Email, userInfo.EmailVerified)
}

// Update Active Status
func (a *OAuthUserApp) UpdateActiveStatus(ctx context.Context, userID string, provider providerv1.OAuthProvider, isActive bool) (*dto.OAuthUser, error) {
	providerEnum, err := util.FromProtoToProvider(provider)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return a.oauthUserService.UpdateActiveStatus(userUUID, providerEnum, isActive)
}

// Update Verified Status
func (a *OAuthUserApp) UpdateVerifiedStatus(ctx context.Context, userID string, provider providerv1.OAuthProvider, isVerified bool) (*dto.OAuthUser, error) {
	providerEnum, err := util.FromProtoToProvider(provider)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return a.oauthUserService.UpdateVerifiedStatus(userUUID, providerEnum, isVerified)
}
