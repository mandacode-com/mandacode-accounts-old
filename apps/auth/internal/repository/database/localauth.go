package dbrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/localauth"
	dbmodels "mandacode.com/accounts/auth/internal/domain/models/database"
	dbdomain "mandacode.com/accounts/auth/internal/domain/repository/database"
)

type localAuthRepository struct {
	client *ent.Client
}

// OnLoginFailed implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) OnLoginFailed(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	localAuth, err = localAuth.Update().
		SetFailedLoginAttempts(localAuth.FailedLoginAttempts + 1).
		SetLastFailedLoginAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to update failed login attempts")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// OnLoginSuccess implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) OnLoginSuccess(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	localAuth, err = localAuth.Update().
		SetLastLoginAt(time.Now()).
		SetFailedLoginAttempts(0).
		SetLastFailedLoginAt(time.Time{}).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to update last login time and reset failed attempts")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// ComparePassword implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) ComparePassword(ctx context.Context, authAccountID uuid.UUID, password string) (bool, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return false, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return false, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	err = bcrypt.CompareHashAndPassword([]byte(localAuth.Password), []byte(password))
	if err != nil {
		return false, nil // Password does not match
	}

	return true, nil // Password matches
}

// CreateLocalAuth implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) CreateLocalAuth(ctx context.Context, input *dbmodels.CreateLocalAuthInput) (*dbmodels.SecureLocalAuth, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to hash password")
	}

	localAuth, err := l.client.LocalAuth.Create().
		SetAuthAccountID(input.AccountID).
		SetEmail(input.Email).
		SetPassword(string(hashedPassword)).
		SetIsActive(input.IsActive).
		SetIsVerified(input.IsVerified).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			upgradedErr := errors.Upgrade(err, errcode.ErrConflict, "User already exists with this email")
			return nil, errors.Join(upgradedErr, "Failed to create local auth record")
		}
		upgradedErr := errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to create local auth record")
		return nil, errors.Join(upgradedErr, "Failed to create local auth record")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// DeleteLocalAuth implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) DeleteLocalAuth(ctx context.Context, id uuid.UUID) error {
	err := l.client.LocalAuth.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found")
		}
		return errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to delete local auth record")
	}
	return nil
}

// DeleteLocalAuthByAuthAccountID implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) DeleteLocalAuthByAuthAccountID(ctx context.Context, authAccountID uuid.UUID) error {
	_, err := l.client.LocalAuth.Delete().
		Where(localauth.AuthAccountID(authAccountID)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to delete local auth record by AuthAccountID")
	}
	return nil
}

// GetLocalAuthByAuthAccountID implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) GetLocalAuthByAuthAccountID(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// GetLocalAuthByEmail implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) GetLocalAuthByEmail(ctx context.Context, email string) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.Email(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given email")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by email")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// GetLocalAuthByID implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) GetLocalAuthByID(ctx context.Context, id uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given ID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by ID")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// ResetFailedLoginAttempts implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) ResetFailedLoginAttempts(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	localAuth, err = localAuth.Update().
		SetFailedLoginAttempts(0).
		SetLastFailedLoginAt(time.Time{}).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to reset failed login attempts")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// SetEmail implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) SetEmail(ctx context.Context, authAccountID uuid.UUID, email string) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	localAuth, err = localAuth.Update().
		SetEmail(email).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to update email in local auth record")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// SetIsActive implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) SetIsActive(ctx context.Context, authAccountID uuid.UUID, isActive bool) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	localAuth, err = localAuth.Update().
		SetIsActive(isActive).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to update active status in local auth record")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// SetIsVerified implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) SetIsVerified(ctx context.Context, authAccountID uuid.UUID, isVerified bool) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	localAuth, err = localAuth.Update().
		SetIsVerified(isVerified).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to update verification status in local auth record")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

// SetPassword implements dbdomain.LocalAuthRepository.
func (l *localAuthRepository) SetPassword(ctx context.Context, authAccountID uuid.UUID, password string) (*dbmodels.SecureLocalAuth, error) {
	localAuth, err := l.client.LocalAuth.Query().
		Where(localauth.AuthAccountID(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found for the given AuthAccountID")
		}
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to retrieve local auth record by AuthAccountID")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to hash password")
	}

	localAuth, err = localAuth.Update().
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to update password in local auth record")
	}

	return dbmodels.NewSecureLocalAuth(localAuth), nil
}

func NewLocalAuthRepository(client *ent.Client) dbdomain.LocalAuthRepository {
	return &localAuthRepository{client: client}
}
