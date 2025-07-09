package groupuserapp

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type groupUserApp struct {
	groupUserRepo repodomain.GroupUserRepository
}

// DeleteGroupUser implements GroupUserApp.
func (g *groupUserApp) DeleteGroupUser(userID uuid.UUID, groupID uuid.UUID) error {
	err := g.groupUserRepo.DeleteGroupUser(userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGroupUserByGroupID implements GroupUserApp.
func (g *groupUserApp) DeleteGroupUserByGroupID(groupID uuid.UUID) error {
	err := g.groupUserRepo.DeleteGroupUserByGroupID(groupID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGroupUserByUserID implements GroupUserApp.
func (g *groupUserApp) DeleteGroupUserByUserID(userID uuid.UUID) error {
	err := g.groupUserRepo.DeleteGroupUserByUserID(userID)
	if err != nil {
		return err
	}
	return nil
}

// CheckGroupUserExists implements GroupUserApp.
func (g *groupUserApp) CheckGroupUserExists(userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	exists, err := g.groupUserRepo.CheckGroupUserExists(userID, groupID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateGroupUser implements GroupUserApp.
func (g *groupUserApp) CreateGroupUser(userID uuid.UUID, groupID uuid.UUID) (*model.GroupUser, error) {
	groupUser, err := g.groupUserRepo.CreateGroupUser(userID, groupID)
	if err != nil {
		return nil, err
	}

	return model.GroupUserFromEnt(groupUser), nil
}

// GetGroupUser implements GroupUserApp.
func (g *groupUserApp) GetGroupUser(userID uuid.UUID, groupID uuid.UUID) (*model.GroupUser, error) {
	groupUser, err := g.groupUserRepo.GetGroupUser(userID, groupID)
	if err != nil {
		return nil, err
	}

	return model.GroupUserFromEnt(groupUser), nil
}


// GetGroupUsersByUserID implements GroupUserApp.
func (g *groupUserApp) GetGroupUsersByUserID(userID uuid.UUID) ([]*model.GroupUser, error) {
	groupUsers, err := g.groupUserRepo.GetGroupUsersByUserID(userID)
	if err != nil {
		return nil, err
	}

	var groupUserModels []*model.GroupUser
	for _, groupUser := range groupUsers {
		groupUserModels = append(groupUserModels, model.GroupUserFromEnt(groupUser))
	}

	return groupUserModels, nil
}

// GetGroupUsersByGroupID implements GroupUserApp.
func (g *groupUserApp) GetGroupUsersByGroupID(groupID uuid.UUID) ([]*model.GroupUser, error) {
	groupUsers, err := g.groupUserRepo.GetGroupUsersByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	var groupUserModels []*model.GroupUser
	for _, groupUser := range groupUsers {
		groupUserModels = append(groupUserModels, model.GroupUserFromEnt(groupUser))
	}

	return groupUserModels, nil
}

// NewGroupUserApp creates a new GroupUserApp.
func NewGroupUserApp(groupUserRepo repodomain.GroupUserRepository) GroupUserApp {
	return &groupUserApp{
		groupUserRepo: groupUserRepo,
	}
}
