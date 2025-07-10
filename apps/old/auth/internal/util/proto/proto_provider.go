package protoutil

import (
	"errors"

	"mandacode.com/accounts/auth/ent/oauthuser"
	providerv1 "github.com/mandacode-com/accounts-proto/common/provider/v1"
)

func FromProtoToEntProvider(provider providerv1.OAuthProvider) (oauthuser.Provider, error) {
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
		return "", ErrUnsupportedProvider
	}
}

func FromEntToProtoProvider(provider oauthuser.Provider) (providerv1.OAuthProvider, error) {
	switch provider {
	case oauthuser.ProviderGoogle:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE, nil
	case oauthuser.ProviderKakao:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_KAKAO, nil
	case oauthuser.ProviderNaver:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_NAVER, nil
	case oauthuser.ProviderFacebook:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_FACEBOOK, nil
	case oauthuser.ProviderGithub:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB, nil
	case oauthuser.ProviderApple:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_APPLE, nil
	default:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED, ErrUnsupportedProvider
	}
}

var ErrUnsupportedProvider = errors.New("unsupported provider")
