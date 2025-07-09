package groupapp

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type groupApp struct {
	groupRepo repodomain.GroupRepository
}

// CreateGroup implements GroupApp.
func (g *groupApp) CreateGroup(name string, serviceID uuid.UUID, isActive *bool, description *string) (*model.Group, error) {
	group, err := g.groupRepo.CreateGroup(name, serviceID, isActive, description)
	if err != nil {
		return nil, err
	}

	return model.GroupFromEnt(group), nil
}

// DeleteGroup implements GroupApp.
func (g *groupApp) DeleteGroup(id uuid.UUID) error {
	err := g.groupRepo.DeleteGroup(id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroupsByServiceID implements GroupApp.
func (g *groupApp) DeleteGroupsByServiceID(serviceID uuid.UUID) error {
	err := g.groupRepo.DeleteGroupsByServiceID(serviceID)
	if err != nil {
		return err
	}

	return nil
}

// GetGroupByID implements GroupApp.
func (g *groupApp) GetGroupByID(id uuid.UUID) (*model.Group, error) {
	group, err := g.groupRepo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, nil // or return an error if you prefer
	}

	return model.GroupFromEnt(group), nil
}

// GetGroupsByServiceID implements GroupApp.
func (g *groupApp) GetGroupsByServiceID(serviceID uuid.UUID) ([]*model.Group, error) {
	groups, err := g.groupRepo.GetGroupsByServiceID(serviceID)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, nil // or return an error if you prefer
	}

	modelGroups := make([]*model.Group, len(groups))
	for i, group := range groups {
		modelGroups[i] = model.GroupFromEnt(group)
	}

	return modelGroups, nil
}

// UpdateGroup implements GroupApp.
func (g *groupApp) UpdateGroup(id uuid.UUID, name *string, serviceID *uuid.UUID, isActive *bool, description *string) (*model.Group, error) {
	group, err := g.groupRepo.UpdateGroup(id, name, serviceID, isActive, description)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, nil // or return an error if you prefer
	}

	return model.GroupFromEnt(group), nil
}

func NewGroupApp(groupRepo repodomain.GroupRepository) GroupApp {
	return &groupApp{
		groupRepo: groupRepo,
	}
}
