package oauthdomain

import oauthdto "mandacode.com/accounts/auth/internal/infra/oauth/dto"


type OAuthProvider interface {
	// GetProfile retrieves the profile information of the authenticated user.
	//
	// Parameters:
	// - accessToken: The OAuth access token for the user.
	//
	// Returns:
	// - userInfo: The user's profile information.
	// - error: An error if the retrieval fails, otherwise nil.
	GetUserInfo(accessToken string) (*oauthdto.OAuthUserInfo, error)
}
