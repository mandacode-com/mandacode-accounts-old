package permissionapp

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type permissionApp struct {
	groupUserRepo    repodomain.GroupUserRepository
	clientAccessRepo repodomain.ClientAccessRepository
	adminGroupID     uuid.UUID
}

// CheckAdmin implements PermissionApp.
func (p *permissionApp) CheckAdmin(userID uuid.UUID) (bool, error) {
	groupUser, err := p.groupUserRepo.CheckGroupUserExists(userID, p.adminGroupID)
	if err != nil {
		return false, err
	}
	return groupUser, nil
}

// CheckClientAccess implements PermissionApp.
func (p *permissionApp) CheckClientAccess(serviceID uuid.UUID, clientID string, clientSecret string) (bool, error) {
	clientAccess, err := p.clientAccessRepo.GetClientAccess(serviceID, clientID)
	if err != nil {
		return false, err
	}
	if clientAccess == nil {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(clientAccess.ClientSecret), []byte(clientSecret)); err != nil {
		return false, nil
	}

	return true, nil
}

// NewPermissionApp creates a new instance of permissionApp.
func NewPermissionApp(
	groupUserRepo repodomain.GroupUserRepository,
	clientAccessRepo repodomain.ClientAccessRepository,
	adminGroupID uuid.UUID,
) PermissionApp {
	return &permissionApp{
		groupUserRepo:    groupUserRepo,
		clientAccessRepo: clientAccessRepo,
		adminGroupID:     adminGroupID,
	}
}
