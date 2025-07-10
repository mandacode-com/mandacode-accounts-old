package permissionapp

import "github.com/google/uuid"

type PermissionApp interface {
	// CheckAdmin checks if a user has admin privileges.
	//
	// Parameters:
	//   - userID: The unique identifier of the user.
	CheckAdmin(userID uuid.UUID) (bool, error)

	// CheckClientAccess checks if a client has access to a specific service.
	//
	// Parameters:
	//   - serviceID: The unique identifier of the service.
	//   - clientID: The unique identifier of the client.
	//   - clientSecret: The secret key for the client.
	CheckClientAccess(serviceID uuid.UUID, clientID string, clientSecret string) (bool, error)
}
