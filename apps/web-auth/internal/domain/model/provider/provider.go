package provider

import (
	"errors"

	providerv1 "mandacode.com/accounts/proto/common/provider/v1"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderNaver  Provider = "naver"
	ProviderKakao  Provider = "kakao"
)

func (p Provider) String() string {
	switch p {
	case ProviderGoogle:
		return "google"
	case ProviderNaver:
		return "naver"
	case ProviderKakao:
		return "kakao"
	default:
		return "unknown"
	}
}

func (p Provider) ToProto() providerv1.OAuthProvider {
	switch p {
	case ProviderGoogle:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE
	case ProviderNaver:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_NAVER
	case ProviderKakao:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_KAKAO
	default:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED
	}
}

func FromString(providerStr string) (Provider, error) {
	switch providerStr {
	case "google":
		return ProviderGoogle, nil
	case "naver":
		return ProviderNaver, nil
	case "kakao":
		return ProviderKakao, nil
	default:
		return "", errors.New("unsupported OAuth provider: " + providerStr)
	}
}
