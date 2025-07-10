package oauthuser

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
	userdto "mandacode.com/accounts/auth/internal/app/user/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/oauth"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	protoutil "mandacode.com/accounts/auth/internal/util/proto"
)

type oauthUserApp struct {
	providers    *map[oauthuser.Provider]oauthdomain.OAuthProvider
	oauthUserRep repodomain.OAuthUserRepository
}

func NewOAuthUserApp(
	providers *map[oauthuser.Provider]oauthdomain.OAuthProvider,
	oauthUserRep repodomain.OAuthUserRepository,
) OAuthUserApp {
	return &oauthUserApp{
		providers:    providers,
		oauthUserRep: oauthUserRep,
	}
}

func (a *oauthUserApp) CreateUser(
	userID uuid.UUID,
	provider oauthuser.Provider,
	oauthAccessToken string,
) (*userdto.OAuthUser, error) {
	providerApp, ok := (*a.providers)[provider]
	if !ok {
		return nil, protoutil.ErrUnsupportedProvider
	}

	userInfo, err := providerApp.GetUserInfo(oauthAccessToken)
	if err != nil {
		return nil, err
	}
	oauthUser, err := a.oauthUserRep.CreateUser(
		userID,
		provider,
		userInfo.ProviderID,
		userInfo.Email,
		nil,
		&userInfo.EmailVerified,
	)
	if err != nil {
		return nil, err
	}
	return userdto.NewOAuthUserFromEnt(oauthUser), nil
}

func (a *oauthUserApp) GetUser(
	userID uuid.UUID,
	provider oauthuser.Provider,
) (*userdto.OAuthUser, error) {
	oauthUser, err := a.oauthUserRep.GetUserByUserID(userID, provider)
	if err != nil {
		return nil, err
	}
	return userdto.NewOAuthUserFromEnt(oauthUser), nil
}

func (a *oauthUserApp) SyncUser(
	userID uuid.UUID,
	provider oauthuser.Provider,
	oauthAccessToken string,
) (*userdto.OAuthUser, error) {
	providerApp, ok := (*a.providers)[provider]
	if !ok {
		return nil, protoutil.ErrUnsupportedProvider
	}

	userInfo, err := providerApp.GetUserInfo(oauthAccessToken)
	if err != nil {
		return nil, err
	}

	oauthUser, err := a.oauthUserRep.UpdateUser(
		userID,
		provider,
		&userInfo.ProviderID,
		&userInfo.Email,
		nil,
		&userInfo.EmailVerified,
	)
	if err != nil {
		return nil, err
	}
	return userdto.NewOAuthUserFromEnt(oauthUser), nil
}

func (a *oauthUserApp) UpdateActiveStatus(
	userID uuid.UUID,
	provider oauthuser.Provider,
	isActive bool,
) (*userdto.OAuthUser, error) {
	oauthUser, err := a.oauthUserRep.UpdateUser(
		userID,
		provider,
		nil,
		nil,
		&isActive,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return userdto.NewOAuthUserFromEnt(oauthUser), nil
}

func (a *oauthUserApp) UpdateVerificationStatus(
	userID uuid.UUID,
	provider oauthuser.Provider,
	isVerified bool,
) (*userdto.OAuthUser, error) {
	oauthUser, err := a.oauthUserRep.UpdateUser(
		userID,
		provider,
		nil,
		nil,
		nil,
		&isVerified,
	)
	if err != nil {
		return nil, err
	}
	return userdto.NewOAuthUserFromEnt(oauthUser), nil
}

func (a *oauthUserApp) DeleteUser(
	userID uuid.UUID,
	provider oauthuser.Provider,
) error {
	return a.oauthUserRep.DeleteUserByProvider(userID, provider)
}

func (a *oauthUserApp) DeleteAllProviders(userID uuid.UUID) error {
	return a.oauthUserRep.DeleteUser(userID)
}
