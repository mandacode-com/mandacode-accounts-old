package userdomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
)

type OAuthUserService interface {
	// CreateOAuthUser creates a new OAuth user with the given details.
	//
	// Parameters:
	//   - userID: The unique identifier for the user.
	//   - provider: The OAuth provider (e.g., "google", "github").
	//   - providerID: The unique identifier provided by the OAuth provider.
	//   - code: The authorization code received from the OAuth provider.
	//   - isActive: A pointer to a boolean indicating if the user is active.
	//   - isVerified: A pointer to a boolean indicating if the user's email is verified.
	//
	// Returns:
	//   - user: The created OAuth user if the creation is successful.
	//   - error: An error if the creation fails, otherwise nil.
	CreateOAuthUser(userID uuid.UUID, provider oauthuser.Provider, providerID, accessToken string, isActive, isVerified *bool) (*dto.OAuthUser, error)

	// DeleteOAuthUser deletes an OAuth user by their userID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be deleted.
	//
	// Returns:
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteOAuthUser(userID uuid.UUID) error

	// DeleteOAuthUserByProvider deletes an OAuth user by their userID and provider.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be deleted.
	//   - provider: The OAuth provider of the user to be deleted.
	//
	// Returns:
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteOAuthUserByProvider(userID uuid.UUID, provider oauthuser.Provider) error

	// UpdateOAuthUser updates the details of an OAuth user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be updated.
	//   - provider: A pointer to the new OAuth provider (nil if not updating).
	//   - code: A pointer to the new authorization code (nil if not updating).
	//   - isActive: A pointer to a boolean indicating if the user is active (nil if not updating).
	//   - isVerified: A pointer to a boolean indicating if the user's email is verified (nil if not updating).
	//
	// Returns:
	//   - user: The updated OAuth user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateOAuthUser(userID uuid.UUID, provider *oauthuser.Provider, providerID, email *string, isActive, isVerified *bool) (*dto.OAuthUser, error)
}
