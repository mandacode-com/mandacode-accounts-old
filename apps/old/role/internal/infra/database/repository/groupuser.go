package repository

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
	"mandacode.com/accounts/role/ent/groupuser"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type GroupUserRepository struct {
	db *ent.Client
}

// DeleteGroupUser implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) DeleteGroupUser(userID uuid.UUID, groupID uuid.UUID) error {
	_, err := g.db.GroupUser.
		Delete().
		Where(
			groupuser.UserID(userID),
			groupuser.GroupID(groupID),
		).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// DeleteGroupUserByGroupID implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) DeleteGroupUserByGroupID(groupID uuid.UUID) error {
	_, err := g.db.GroupUser.
		Delete().
		Where(groupuser.GroupID(groupID)).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// DeleteGroupUserByUserID implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) DeleteGroupUserByUserID(userID uuid.UUID) error {
	_, err := g.db.GroupUser.
		Delete().
		Where(groupuser.UserID(userID)).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// GetGroupUsersByGroupID implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) GetGroupUsersByGroupID(groupID uuid.UUID) ([]*ent.GroupUser, error) {
	groupUsers, err := g.db.GroupUser.
		Query().
		Where(groupuser.GroupID(groupID)).
		All(context.Background())

	if err != nil {
		return nil, err
	}

	return groupUsers, nil
}

// GetGroupUsersByUserID implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) GetGroupUsersByUserID(userID uuid.UUID) ([]*ent.GroupUser, error) {
	groupUsers, err := g.db.GroupUser.
		Query().
		Where(groupuser.UserID(userID)).
		All(context.Background())

	if err != nil {
		return nil, err
	}

	return groupUsers, nil
}

// GetGroupUser implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) GetGroupUser(userID uuid.UUID, groupID uuid.UUID) (*ent.GroupUser, error) {
	groupUser, err := g.db.GroupUser.
		Query().
		Where(
			groupuser.UserID(userID),
			groupuser.GroupID(groupID),
		).
		Only(context.Background())
	if err != nil {
		return nil, err
	}

	return groupUser, nil
}

// CheckGroupUserExists implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) CheckGroupUserExists(userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	exists, err := g.db.GroupUser.
		Query().
		Where(
			groupuser.UserID(userID),
			groupuser.GroupID(groupID),
		).
		Exist(context.Background())

	return exists, err
}

// CreateGroupUser implements repodomain.GroupUserRepository.
func (g *GroupUserRepository) CreateGroupUser(userID uuid.UUID, groupID uuid.UUID) (*ent.GroupUser, error) {
	create := g.db.GroupUser.Create()

	create.SetUserID(userID)
	create.SetGroupID(groupID)

	groupUser, err := create.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return groupUser, nil
}

// NewGroupUserRepository creates a new GroupUserRepository.
func NewGroupUserRepository(db *ent.Client) repodomain.GroupUserRepository {
	return &GroupUserRepository{
		db: db,
	}
}
