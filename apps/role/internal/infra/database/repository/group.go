package repository

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
	"mandacode.com/accounts/role/ent/group"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type GroupRepository struct {
	db *ent.Client
}

// CreateGroup implements repodomain.GroupRepository.
func (g *GroupRepository) CreateGroup(name string, serviceID uuid.UUID, isActive *bool, description *string) (*ent.Group, error) {
	create := g.db.Group.Create()

	create.SetID(uuid.New())
	create.SetName(name)
	create.SetServiceID(serviceID)

	if isActive != nil {
		create.SetIsActive(*isActive)
	}

	if description != nil {
		create.SetDescription(*description)
	}

	group, err := create.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return group, nil
}

// DeleteGroup implements repodomain.GroupRepository.
func (g *GroupRepository) DeleteGroup(id uuid.UUID) error {
	err := g.db.Group.
		DeleteOneID(id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// DeleteGroupsByServiceID implements repodomain.GroupRepository.
func (g *GroupRepository) DeleteGroupsByServiceID(serviceID uuid.UUID) error {
	_, err := g.db.Group.
		Delete().
		Where(group.ServiceID(serviceID)).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// GetGroupByID implements repodomain.GroupRepository.
func (g *GroupRepository) GetGroupByID(id uuid.UUID) (*ent.Group, error) {
	group, err := g.db.Group.
		Query().
		Where(group.ID(id)).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroupsByServiceID implements repodomain.GroupRepository.
func (g *GroupRepository) GetGroupsByServiceID(serviceID uuid.UUID) ([]*ent.Group, error) {
	groups, err := g.db.Group.
		Query().
		Where(group.ServiceID(serviceID)).
		All(context.Background())

	if err != nil {
		return nil, err
	}

	return groups, nil
}

// UpdateGroup implements repodomain.GroupRepository.
func (g *GroupRepository) UpdateGroup(id uuid.UUID, name *string, serviceID *uuid.UUID, isActive *bool, description *string) (*ent.Group, error) {
	update := g.db.Group.UpdateOneID(id)

	if name != nil {
		update.SetName(*name)
	}
	if serviceID != nil {
		update.SetServiceID(*serviceID)
	}
	if isActive != nil {
		update.SetIsActive(*isActive)
	}
	if description != nil {
		update.SetDescription(*description)
	}

	group, err := update.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return group, nil
}

func NewGroupRepository(db *ent.Client) repodomain.GroupRepository {
	return &GroupRepository{db: db}
}
