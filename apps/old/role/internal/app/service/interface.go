package serviceapp

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
)

type ServiceApp interface {
	// CreateService
	//
	// Parameters:
	//   - name: The name of the service.
	//   - description: A brief description of the service. (optional)
	//
	// Returns:
	//   - *model.Service: The created service entity.
	//   - error: An error if the service could not be created.
	CreateService(name string, description *string) (*model.Service, error)

	// GetServiceByID
	//
	// Parameters:
	//   - id: The unique uuid of the service.
	//
	// Returns:
	//   - *model.Service: The service entity if found.
	//   - error: An error if the service could not be found or another error occurred.
	GetServiceByID(id uuid.UUID) (*model.Service, error)

	// GetAllServices
	//
	// Returns:
	//   - []*model.Service: A slice of all service entities.
	//   - error: An error if the services could not be retrieved.
	GetAllServices() ([]*model.Service, error)

	// UpdateService
	//
	// Parameters:
	//   - id: The unique uuid of the service.
	//   - name: The new name of the service. (optional)
	//   - description: The new description of the service. (optional)
	//
	// Returns:
	//   - *model.Service: The updated service entity.
	//   - error: An error if the service could not be updated.
	UpdateService(id uuid.UUID, name *string, description *string) (*model.Service, error)

	// DeleteService
	//
	// Parameters:
	//   - id: The unique uuid of the service.
	//
	// Returns:
	//   - error: An error if the service could not be deleted.
	DeleteService(id uuid.UUID) error
}
