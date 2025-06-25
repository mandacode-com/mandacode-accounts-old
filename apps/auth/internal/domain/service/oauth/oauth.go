package oauthdomain

import "mandacode.com/accounts/auth/internal/domain/dto"

type OAuthService interface {
	// GetProfile retrieves the profile information of the authenticated user.
	//
	// Parameters:
	// - accessToken: The OAuth access token for the user.
	//
	// Returns:
	// - userInfo: The user's profile information.
	// - error: An error if the retrieval fails, otherwise nil.
	GetUserInfo(accessToken string) (*dto.OAuthUserInfo, error)
}
