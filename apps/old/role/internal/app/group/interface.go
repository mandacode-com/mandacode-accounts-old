package groupapp

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
)

type GroupApp interface {
	// CreateGroup
	//
	// Parameters:
	//   - name: The name of the group.
	//   - serviceID: The unique uuid of the service to which the group belongs.
	//   - isActive(optional): Indicates whether the group is active or not.
	//   - description(optional): A brief description of the group.
	//
	// Returns:
	//   - *model.Group: The created group entity.
	//   - error: An error if the group could not be created.
	CreateGroup(name string, serviceID uuid.UUID, isActive *bool, description *string) (*model.Group, error)

	// GetGroupByID
	//
	// Parameters:
	//   - id: The unique uuid of the group.
	//
	// Returns:
	//   - *model.Group: The group entity if found.
	//   - error: An error if the group could not be found or another error occurred.
	GetGroupByID(id uuid.UUID) (*model.Group, error)

	// GetGroupsByServiceID
	//
	// Parameters:
	//   - serviceID: The unique uuid of the service for which groups should be retrieved.
	//
	// Returns:
	//   - []*model.Group: A slice of group entities associated with the service.
	//	 - error: An error if the groups could not be retrieved.
	GetGroupsByServiceID(serviceID uuid.UUID) ([]*model.Group, error)

	// UpdateGroup
	//
	// Parameters:
	//   - id: The unique uuid of the group.
	//   - name(optional): The new name of the group.
	//   - serviceID(optional): The new service ID to which the group belongs.
	//   - isActive(optional): The new active status of the group.
	//   - description(optional): The new description of the group.
	//
	// Returns:
	//   - *model.Group: The updated group entity.
	//   - error: An error if the group could not be updated.
	UpdateGroup(id uuid.UUID, name *string, serviceID *uuid.UUID, isActive *bool, description *string) (*model.Group, error)

	// DeleteGroup
	//
	// Parameters:
	//   - id: The unique uuid of the group.
	//
	// Returns:
	//   - error: An error if the group could not be deleted.
	DeleteGroup(id uuid.UUID) error

	// DeleteGroupsByServiceID
	//
	// Parameters:
	//   - serviceID: The unique uuid of the service for which groups should be deleted.
	//
	// Returns:
	//   - error: An error if the groups could not be deleted.
	DeleteGroupsByServiceID(serviceID uuid.UUID) error
}
