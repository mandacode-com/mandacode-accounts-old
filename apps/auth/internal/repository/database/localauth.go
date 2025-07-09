package dbrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/localauth"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
)

type LocalAuthRepository struct {
	client *ent.Client
}

// OnLoginFailed implements the logic for handling a failed login attempt.
func (l *LocalAuthRepository) OnLoginFailed(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
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

// OnLoginSuccess implements the logic for handling a successful login attempt.
func (l *LocalAuthRepository) OnLoginSuccess(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
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

// ComparePassword checks if the provided password matches the stored hashed password for the given AuthAccountID.
func (l *LocalAuthRepository) ComparePassword(ctx context.Context, authAccountID uuid.UUID, password string) (bool, error) {
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

// CreateLocalAuth creates a new local authentication record with the provided input.
func (l *LocalAuthRepository) CreateLocalAuth(ctx context.Context, input *dbmodels.CreateLocalAuthInput) (*dbmodels.SecureLocalAuth, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to hash password")
	}

	localAuth, err := l.client.LocalAuth.Create().
		SetAuthAccountID(input.AccountID).
		SetEmail(input.Email).
		SetPassword(string(hashedPassword)).
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

// DeleteLocalAuth deletes a local authentication record by its ID.
func (l *LocalAuthRepository) DeleteLocalAuth(ctx context.Context, id uuid.UUID) error {
	err := l.client.LocalAuth.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.Upgrade(err, errcode.ErrNotFound, "LocalAuth not found")
		}
		return errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to delete local auth record")
	}
	return nil
}

// DeleteLocalAuthByAuthAccountID deletes a local authentication record by its associated AuthAccountID.
func (l *LocalAuthRepository) DeleteLocalAuthByAuthAccountID(ctx context.Context, authAccountID uuid.UUID) error {
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

// GetLocalAuthByAuthAccountID retrieves a local authentication record by its associated AuthAccountID.
func (l *LocalAuthRepository) GetLocalAuthByAuthAccountID(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
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

// GetLocalAuthByEmail retrieves a local authentication record by its email address.
func (l *LocalAuthRepository) GetLocalAuthByEmail(ctx context.Context, email string) (*dbmodels.SecureLocalAuth, error) {
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

// GetLocalAuthByID retrieves a local authentication record by its ID.
func (l *LocalAuthRepository) GetLocalAuthByID(ctx context.Context, id uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
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

// ResetFailedLoginAttempts resets the failed login attempts for a local authentication record.
func (l *LocalAuthRepository) ResetFailedLoginAttempts(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error) {
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

// SetEmail updates the email address for a local authentication record.
func (l *LocalAuthRepository) SetEmail(ctx context.Context, authAccountID uuid.UUID, email string) (*dbmodels.SecureLocalAuth, error) {
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

// SetIsVerified updates the verification status of a local authentication record.
func (l *LocalAuthRepository) SetIsVerified(ctx context.Context, authAccountID uuid.UUID, isVerified bool) (*dbmodels.SecureLocalAuth, error) {
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

// SetPassword updates the password for a local authentication record.
func (l *LocalAuthRepository) SetPassword(ctx context.Context, authAccountID uuid.UUID, password string) (*dbmodels.SecureLocalAuth, error) {
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

func NewLocalAuthRepository(client *ent.Client) *LocalAuthRepository {
	return &LocalAuthRepository{client: client}
}
