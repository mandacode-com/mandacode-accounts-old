package dbrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/authaccount"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
)

type AuthAccountRepository struct {
	client *ent.Client
}

// CreateAuthAccount creates a new authentication account.
func (a *AuthAccountRepository) CreateAuthAccount(ctx context.Context, account *dbmodels.CreateAuthAccountInput) (*ent.AuthAccount, error) {
	create := a.client.AuthAccount.Create().
		SetID(uuid.New()).
		SetUserID(account.UserID)

	authAccount, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("AuthAccount already exists", "Conflict", errcode.ErrConflict)
		}
		return nil, errors.New(err.Error(), "Failed to create AuthAccount", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// DeleteAuthAccount deletes an authentication account by its ID.
func (a *AuthAccountRepository) DeleteAuthAccount(ctx context.Context, id uuid.UUID) error {
	delete := a.client.AuthAccount.DeleteOneID(id)

	err := delete.Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("AuthAccount not found", "NotFound", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete AuthAccount", errcode.ErrInternalFailure)
	}

	return nil
}

// GetAuthAccountByID retrieves an authentication account by its ID.
func (a *AuthAccountRepository) GetAuthAccountByID(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get AuthAccount by ID", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// GetAuthAccountByUserID retrieves an authentication account by the associated user ID.
func (a *AuthAccountRepository) GetAuthAccountByUserID(ctx context.Context, userID uuid.UUID) (*ent.AuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.UserIDEQ(userID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found for user ID", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get AuthAccount by User ID", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// OnLoginFailed sets the last failed login time and increments the failed login attempts for an authentication account.
func (a *AuthAccountRepository) OnLoginFailed(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error) {
	authAccount, err := a.client.AuthAccount.UpdateOneID(id).
		SetLastFailedLoginAt(time.Now()).
		AddFailedLoginAttempts(1).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to update AuthAccount on login failure", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// OnLoginSuccess sets the last login time, resets the last failed login time, and resets the failed login attempts for an authentication account.
func (a *AuthAccountRepository) OnLoginSuccess(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error) {
	authAccount, err := a.client.AuthAccount.UpdateOneID(id).
		SetLastLoginAt(time.Now()).
		SetLastFailedLoginAt(time.Time{}).
		SetFailedLoginAttempts(0).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to update AuthAccount on login success", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// ResetFailedLoginAttempts resets the failed login attempts for an authentication account.
func (a *AuthAccountRepository) ResetFailedLoginAttempts(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error) {
	authAccount, err := a.client.AuthAccount.UpdateOneID(id).
		SetFailedLoginAttempts(0).
		SetLastFailedLoginAt(time.Time{}).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to reset failed login attempts", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// UpdateAuthAccount updates an existing authentication account.
func (a *AuthAccountRepository) UpdateAuthAccount(ctx context.Context, id uuid.UUID, account *dbmodels.UpdateAuthAccountInput) (*ent.AuthAccount, error) {
	update := a.client.AuthAccount.UpdateOneID(id)

	if account.UserID != nil {
		update = update.SetUserID(*account.UserID)
	}

	if account.LastLoginAt != nil {
		update = update.SetLastLoginAt(*account.LastLoginAt)
	}

	if account.LastFailedLoginAt != nil {
		update = update.SetLastFailedLoginAt(*account.LastFailedLoginAt)
	}

	authAccount, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to update AuthAccount", errcode.ErrInternalFailure)
	}

	return authAccount, nil
}

// NewAuthAccountRepository creates a new instance of authAccountRepository.
func NewAuthAccountRepository(client *ent.Client) *AuthAccountRepository {
	return &AuthAccountRepository{
		client: client,
	}
}
