package util

import providerv1 "github.com/mandacode-com/accounts-proto/common/provider/v1"

// ConvertToOAuthProvider converts a string representation of an OAuth provider to the corresponding OAuthProvider enum value.
//
// Parameter:
//   - provider: The string representation of the OAuth provider (e.g., "google", "naver", "kakao").
func ConvertToOAuthProvider(provider string) providerv1.OAuthProvider {
	switch provider {
	case "google":
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE
	case "naver":
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_NAVER
	case "kakao":
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_KAKAO
	default:
		return providerv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED
	}
}
