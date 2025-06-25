package util

import (
	"errors"

	"mandacode.com/accounts/auth/ent/oauthuser"
	providerv1 "mandacode.com/accounts/auth/proto/common/provider/v1"
)

func FromStrToProvider(provider string) (oauthuser.Provider, error) {
	switch provider {
	case "google":
		return oauthuser.ProviderGoogle, nil
	case "kakao":
		return oauthuser.ProviderKakao, nil
	case "naver":
		return oauthuser.ProviderNaver, nil
	case "facebook":
		return oauthuser.ProviderFacebook, nil
	case "github":
		return oauthuser.ProviderGithub, nil
	case "apple":
		return oauthuser.ProviderApple, nil
	default:
		return "", errors.New("unsupported provider: " + provider)
	}
}

func FromProtoToProvider(provider providerv1.OAuthProvider) (oauthuser.Provider, error) {
	switch provider {
	case providerv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE:
		return oauthuser.ProviderGoogle, nil
	case providerv1.OAuthProvider_O_AUTH_PROVIDER_KAKAO:
		return oauthuser.ProviderKakao, nil
	case providerv1.OAuthProvider_O_AUTH_PROVIDER_NAVER:
		return oauthuser.ProviderNaver, nil
	case providerv1.OAuthProvider_O_AUTH_PROVIDER_FACEBOOK:
		return oauthuser.ProviderFacebook, nil
	case providerv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB:
		return oauthuser.ProviderGithub, nil
	case providerv1.OAuthProvider_O_AUTH_PROVIDER_APPLE:
		return oauthuser.ProviderApple, nil
	default:
		return "", errors.New("unsupported provider: " + provider.String())
	}
}
