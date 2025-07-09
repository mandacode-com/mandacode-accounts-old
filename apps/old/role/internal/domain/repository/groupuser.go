package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type GroupUserRepository interface {
	// CreateGroupUser
	//
	// Parameters:
	//   - userID: The unique uuid of the user.
	//   - groupID: The unique uuid of the group.
	//
	// Returns:
	//   - *ent.GroupUser: The created group user entity.
	//   - error: An error if the group user could not be created.
	CreateGroupUser(userID uuid.UUID, groupID uuid.UUID) (*ent.GroupUser, error)

	// CheckGroupUserExists
	//
	// Parameters:
	//   - userID: The unique uuid of the user.
	//   - groupID: The unique uuid of the group.
	// Returns:
	//   - bool: True if the group user exists, false otherwise.
	//   - error: An error if the check could not be performed.
	CheckGroupUserExists(userID uuid.UUID, groupID uuid.UUID) (bool, error)

	// GetGroupUser
	//
	// Parameters:
	//   - userID: The unique uuid of the user.
	//   - groupID: The unique uuid of the group.
	//
	// Returns:
	//   - *ent.GroupUser: The group user entity if found.
	//   - error: An error if the group user could not be found or another error occurred.
	GetGroupUser(userID uuid.UUID, groupID uuid.UUID) (*ent.GroupUser, error)

	// GetGroupUsersByUserID
	//
	// Parameters:
	//   - userID: The unique uuid of the user.
	// Returns:
	//   - []*ent.GroupUser: A slice of group user entities associated with the user.
	//   - error: An error if the group users could not be retrieved.
	GetGroupUsersByUserID(userID uuid.UUID) ([]*ent.GroupUser, error)

	// GetGroupUsersByGroupID
	//
	// Parameters:
	//   - groupID: The unique uuid of the group.
	//
	// Returns:
	//   - []*ent.GroupUser: A slice of group user entities associated with the group.
	//   - error: An error if the group users could not be retrieved.
	GetGroupUsersByGroupID(groupID uuid.UUID) ([]*ent.GroupUser, error)

	// DeleteGroupUserByUserID
	//
	// Parameters:
	//   - userID: The unique uuid of the user.
	//
	// Returns:
	//   - error: An error if the group user could not be deleted.
	DeleteGroupUserByUserID(userID uuid.UUID) error

	// DeleteGroupUserByGroupID
	//
	// Parameters:
	//   - groupID: The unique uuid of the group.
	//
	// Returns:
	//   - error: An error if the group user could not be deleted.
	DeleteGroupUserByGroupID(groupID uuid.UUID) error

	// DeleteGroupUser
	//
	// Parameters:
	//   - userID: The unique uuid of the user.
	//   - groupID: The unique uuid of the group.
	//
	// Returns:
	//   - error: An error if the group user could not be deleted.
	DeleteGroupUser(userID uuid.UUID, groupID uuid.UUID) error
}
