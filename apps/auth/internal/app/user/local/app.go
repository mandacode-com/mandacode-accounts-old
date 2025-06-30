package localuser

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	userdto "mandacode.com/accounts/auth/internal/app/user/dto"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
)

type localUserApp struct {
	localUserRep repodomain.LocalUserRepository
}

// NewLocalUserApp creates a new instance of LocalUserApp.
func NewLocalUserApp(localUserRep repodomain.LocalUserRepository) LocalUserApp {
	return &localUserApp{
		localUserRep: localUserRep,
	}
}

// CreateUser implements LocalUserApp.
func (l *localUserApp) CreateUser(userID uuid.UUID, email string, password string, isActive *bool, isVerified *bool) (*userdto.LocalUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	localUser, err := l.localUserRep.CreateUser(
		userID,
		email,
		string(hashedPassword),
		isActive,
		isVerified,
	)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}

// DeleteUser implements LocalUserApp.
func (l *localUserApp) DeleteUser(userID uuid.UUID) error {
	err := l.localUserRep.DeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser implements LocalUserApp.
func (l *localUserApp) GetUser(userID uuid.UUID) (*userdto.LocalUser, error) {
	localUser, err := l.localUserRep.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}

// GetUserByEmail implements LocalUserApp.
func (l *localUserApp) GetUserByEmail(email string) (*userdto.LocalUser, error) {
	localUser, err := l.localUserRep.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}

// UpdateActiveStatus implements LocalUserApp.
func (l *localUserApp) UpdateActiveStatus(userID uuid.UUID, isActive bool) (*userdto.LocalUser, error) {
	localUser, err := l.localUserRep.UpdateUser(
		userID,
		nil, // email
		nil, // password
		&isActive,
		nil, // isVerified
	)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}

// UpdateEmail implements LocalUserApp.
func (l *localUserApp) UpdateEmail(userID uuid.UUID, email string) (*userdto.LocalUser, error) {
	localUser, err := l.localUserRep.UpdateUser(
		userID,
		&email, // email
		nil,    // password
		nil,    // isActive
		nil,    // isVerified
	)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}

// UpdatePassword implements LocalUserApp.
func (l *localUserApp) UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) (*userdto.LocalUser, error) {
	// Fetch the current user to verify the current password
	localUser, err := l.localUserRep.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Verify the current password
	err = bcrypt.CompareHashAndPassword([]byte(localUser.Password), []byte(currentPassword))
	if err != nil {
		return nil, err // Current password is incorrect
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	stringHashedNewPassword := string(hashedNewPassword)

	// Update the user's password
	localUser, err = l.localUserRep.UpdateUser(
		userID,
		nil,                      // email
		&stringHashedNewPassword, // password
		nil,                      // isActive
		nil,                      // isVerified
	)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}

// UpdateVerificationStatus implements LocalUserApp.
func (l *localUserApp) UpdateVerificationStatus(userID uuid.UUID, isVerified bool) (*userdto.LocalUser, error) {
	localUser, err := l.localUserRep.UpdateUser(
		userID,
		nil, // email
		nil, // password
		nil, // isActive
		&isVerified,
	)
	if err != nil {
		return nil, err
	}

	return userdto.NewLocalUserFromEnt(localUser), nil
}
