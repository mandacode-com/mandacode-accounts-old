package models

type OAuthUserInfo struct {
	ProviderID    string `validate:"required"`
	Email         string `validate:"required,email"`
	Name          string `validate:"required"`
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
