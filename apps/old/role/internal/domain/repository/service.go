package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type ServiceRepository interface {
	// CreateService
	//
	// Parameters:
	//   - name: The name of the service.
	//   - description(optional): A brief description of the service.
	//
	// Returns:
	//   - *ent.Service: The created service entity.
	//   - error: An error if the service could not be created.
	CreateService(name string, description *string) (*ent.Service, error)

	// GetServiceByID
	//
	// Parameters:
	//   - id: The unique uuid of the service.
	//
	// Returns:
	//   - *ent.Service: The service entity if found.
	//   - error: An error if the service could not be found or another error occurred.
	GetServiceByID(id uuid.UUID) (*ent.Service, error)

	// GetAllServices
	//
	// Returns:
	//   - []*ent.Service: A slice of all service entities.
	//   - error: An error if the services could not be retrieved.
	GetAllServices() ([]*ent.Service, error)

	// UpdateService
	//
	// Parameters:
	//   - id: The unique uuid of the service.
	//   - name(optional): The new name of the service.
	//   - description(optional): The new description of the service.
	//
	// Returns:
	//   - *ent.Service: The updated service entity.
	//   - error: An error if the service could not be updated.
	UpdateService(id uuid.UUID, name *string, description *string) (*ent.Service, error)

	// DeleteService
	//
	// Parameters:
	//   - id: The unique uuid of the service.
	//
	// Returns:
	//   - error: An error if the service could not be deleted.
	DeleteService(id uuid.UUID) error
}
