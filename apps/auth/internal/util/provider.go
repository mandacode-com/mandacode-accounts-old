package util

import (
	providerv1 "github.com/mandacode-com/accounts-proto/go/oauth/provider/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent/oauthauth"
)

func ConvertToProto(provider string) providerv1.OAuthProvider {
	switch provider {
	case "google":
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE
	case "kakao":
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_KAKAO
	case "naver":
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_NAVER
	default:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED
	}
}

func ConvertToEnt(provider string) (oauthauth.Provider, error) {
	switch provider {
	case "google":
		return oauthauth.ProviderGoogle, nil
	case "kakao":
		return oauthauth.ProviderKakao, nil
	case "naver":
		return oauthauth.ProviderNaver, nil
	default:
		return "", errors.New("unsupported provider", "UnsupportedProvider", errcode.ErrInvalidInput)
	}
}
