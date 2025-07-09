package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type ClientAccessRepository interface {
	// CreateClientAccess
	//
	// Parameters:
	//   - serviceID: The unique uuid of the service.
	//   - clientID: The unique identifier for the client (e.g., a UUID or a string).
	//   - clientSecret: The secret key for the client.
	//   - isActive(optional): Indicates whether the service client is active or not.
	CreateClientAccess(serviceID uuid.UUID, clientID string, clientSecret string, isActive *bool) (*ent.ClientAccess, error)

	// GetClientAccessByID
	//
	// Parameters:
	//   - id: The unique uuid of the client access to retrieve.
	GetClientAccessByID(id uuid.UUID) (*ent.ClientAccess, error)

	// GetClientAccess
	//
	// Parameters:
	//   - serviceID: The unique uuid of the service.
	//   - clientID: The unique identifier for the client (e.g., a UUID or a string).
	GetClientAccess(serviceID uuid.UUID, clientID string) (*ent.ClientAccess, error)

	// GetClientAccessesByServiceID
	//
	// Parameters:
	//   - serviceID: The unique uuid of the service for which client accesses should be retrieved.
	GetClientAccessesByServiceID(serviceID uuid.UUID) ([]*ent.ClientAccess, error)

	// UpdateClientAccess
	//
	// Parameters:
	//   - id: The unique uuid of the client access to update.
	//   - serviceID: The unique uuid of the service (optional).
	//   - clientID: The unique identifier for the client (optional).
	//   - clientSecret: The secret key for the client (optional).
	//   - isActive: Indicates whether the service client is active or not (optional).
	UpdateClientAccess(id uuid.UUID, serviceID *uuid.UUID, clientID *string, clientSecret *string, isActive *bool) (*ent.ClientAccess, error)

	// DeleteClientAccess
	//
	// Parameters:
	//   - id: The unique uuid of the client access to delete.
	DeleteClientAccess(id uuid.UUID) error

	// DeleteClientAccessByServiceID
	//
	// Parameters:
	//   - serviceID: The unique uuid of the service for which all client accesses should be deleted.
	DeleteClientAccessByServiceID(serviceID uuid.UUID) error
}
