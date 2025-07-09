package oauthuser

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
	userdto "mandacode.com/accounts/auth/internal/app/user/dto"
)

type OAuthUserApp interface {
	// CreateUser creates a new OAuth user with the given userID, provider, and OAuth access token.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - provider: The OAuth provider used for authentication.
	// - oauthAccessToken: The OAuth access token for the user.
	//
	// Returns:
	// - *OAuthUser: The created OAuth user.
	// - error: An error if the creation fails, otherwise nil.
	CreateUser(userID uuid.UUID, provider oauthuser.Provider, oauthAccessToken string) (*userdto.OAuthUser, error)

	// GetUser retrieves an OAuth user by their userID and provider.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - provider: The OAuth provider used for authentication.
	//
	// Returns:
	// - *OAuthUser: The retrieved OAuth user.
	// - error: An error if the retrieval fails, otherwise nil.
	GetUser(userID uuid.UUID, provider oauthuser.Provider) (*userdto.OAuthUser, error)

	// SyncUser synchronizes an OAuth user with the given userID, provider, and OAuth access token.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - provider: The OAuth provider used for authentication.
	// - oauthAccessToken: The OAuth access token for the user.
	//
	// Returns:
	// - *OAuthUser: The synchronized OAuth user.
	// - error: An error if the synchronization fails, otherwise nil.
	SyncUser(userID uuid.UUID, provider oauthuser.Provider, oauthAccessToken string) (*userdto.OAuthUser, error)

	// UpdateActiveStatus updates the active status of an OAuth user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - provider: The OAuth provider used for authentication.
	// - isActive: The new active status to set for the user.
	//
	// Returns:
	// - *OAuthUser: The updated OAuth user.
	// - error: An error if the update fails, otherwise nil.
	UpdateActiveStatus(userID uuid.UUID, provider oauthuser.Provider, isActive bool) (*userdto.OAuthUser, error)

	// UpdateVerificationStatus updates the verification status of an OAuth user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - provider: The OAuth provider used for authentication.
	// - isVerified: The new verification status to set for the user.
	//
	// Returns:
	// - *OAuthUser: The updated OAuth user.
	// - error: An error if the update fails, otherwise nil.
	UpdateVerificationStatus(userID uuid.UUID, provider oauthuser.Provider, isVerified bool) (*userdto.OAuthUser, error)

	// DeleteUser deletes an OAuth user by their userID and provider.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - provider: The OAuth provider used for authentication.
	//
	// Returns:
	// - error: An error if the deletion fails, otherwise nil.
	DeleteUser(userID uuid.UUID, provider oauthuser.Provider) error

	// DeleteAllProviders deletes all OAuth providers associated with a user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	//
	// Returns:
	// - error: An error if the deletion fails, otherwise nil.
	DeleteAllProviders(userID uuid.UUID) error
}
