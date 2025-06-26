package userdomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
)

type OAuthUserService interface {
	// GetUserByProvider retrieves an OAuth user by their provider and provider ID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user.
	//   - provider: The OAuth provider (e.g., "google", "github").
	//
	// Returns:
	//   - user: The OAuth user if found, otherwise nil.
	//   - error: An error if the retrieval fails, otherwise nil.
	GetUserByProvider(userID uuid.UUID, provider oauthuser.Provider) (*dto.OAuthUser, error)

	// CreateUser creates a new OAuth user with the given details.
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
	CreateUser(userID uuid.UUID, provider oauthuser.Provider, providerID, accessToken string, isActive, isVerified *bool) (*dto.OAuthUser, error)

	// DeleteUser deletes an OAuth user by their userID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be deleted.
	//
	// Returns:
	//   - user: The deleted OAuth user if the deletion is successful.
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteUser(userID uuid.UUID) (*dto.OAuthDeletedUser, error)

	// DeleteUserByProvider deletes an OAuth user by their userID and provider.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be deleted.
	//   - provider: The OAuth provider of the user to be deleted.
	//
	// Returns:
	//   - user: The deleted OAuth user if the deletion is successful.
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteUserByProvider(userID uuid.UUID, provider oauthuser.Provider) (*dto.OAuthDeletedUser, error)

	// UpdateUserBase updates an OAuth user's details.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be updated.
	//   - provider: The OAuth provider.
	//   - providerID: The unique identifier provided by the OAuth provider.
	//   - email: The new email address for the user.
	//   - isVerified: A boolean indicating if the user's email is verified.
	//
	// Returns:
	//   - user: The updated OAuth user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateUserBase(userID uuid.UUID, provider oauthuser.Provider, providerID string, email string, isVerified bool) (*dto.OAuthUser, error)

	// UpdateActiveStatus updates the active status of an OAuth user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose active status is to be updated.
	//   - provider: The OAuth provider of the user whose active status is to be updated.
	//   - isActive: A boolean indicating if the user is active.
	//
	// Returns:
	//   - user: The updated OAuth user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateActiveStatus(userID uuid.UUID, provider oauthuser.Provider, isActive bool) (*dto.OAuthUser, error)

	// UpdateVerifiedStatus updates the verified status of an OAuth user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose verified status is to be updated.
	//   - provider: The OAuth provider of the user whose verified status is to be updated.
	//   - isVerified: A boolean indicating if the user's email is verified.
	//
	// Returns:
	//   - user: The updated OAuth user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateVerifiedStatus(userID uuid.UUID, provider oauthuser.Provider, isVerified bool) (*dto.OAuthUser, error)
}
