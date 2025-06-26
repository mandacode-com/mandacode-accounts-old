package dto

type OAuthUserInfo struct {
	ProviderID    string
	Email         string
	Name          string
	EmailVerified bool
}

func NewOAuthUserInfo(providerID string, email string, name string, emailVerified bool) *OAuthUserInfo {
	return &OAuthUserInfo{
		ProviderID:    providerID,
		Email:         email,
		Name:          name,
		EmailVerified: emailVerified,
	}
}
