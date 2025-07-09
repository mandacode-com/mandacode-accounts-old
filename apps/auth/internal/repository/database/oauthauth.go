package dbrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthauth"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
)

type OAuthAuthRepository struct {
	client *ent.Client
}

// OnLoginFailed sets the last failed login time and increments the failed login attempts for an OAuthAuth record.
func (o *OAuthAuthRepository) OnLoginFailed(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth for login failure", errcode.ErrInternalFailure)
	}

	auth, err = o.client.OAuthAuth.UpdateOne(auth).
		SetFailedLoginAttempts(auth.FailedLoginAttempts + 1).
		SetLastFailedLoginAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to update OAuthAuth on login failure", errcode.ErrInternalFailure)
	}

	return auth, nil
}

// OnLoginSuccess sets the last login time, resets the last failed login time, and resets the failed login attempts for an OAuthAuth record.
func (o *OAuthAuthRepository) OnLoginSuccess(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth for login success", errcode.ErrInternalFailure)
	}

	auth, err = o.client.OAuthAuth.UpdateOne(auth).
		SetFailedLoginAttempts(0).
		SetLastFailedLoginAt(time.Time{}).
		SetLastLoginAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to update OAuthAuth on login success", errcode.ErrInternalFailure)
	}

	return auth, nil
}

// DeleteOAuthAuthByProviderID deletes an OAuthAuth record by provider and provider ID.
func (o *OAuthAuthRepository) DeleteOAuthAuthByProviderID(ctx context.Context, provider oauthauth.Provider, providerID string) error {
	_, err := o.client.OAuthAuth.Delete().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.ProviderIDEQ(providerID)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete OAuthAuth record", errcode.ErrInternalFailure)
	}
	return nil
}

// GetOAuthAuthByAuthAccountID retrieves an OAuthAuth record by provider and auth account ID.
func (o *OAuthAuthRepository) GetOAuthAuthByAuthAccountID(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth by account ID", errcode.ErrInternalFailure)
	}
	return auth, nil
}

// GetOAuthAuthByProviderID retrieves an OAuthAuth record by provider and provider ID.
func (o *OAuthAuthRepository) GetOAuthAuthByProviderID(ctx context.Context, provider oauthauth.Provider, providerID string) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.ProviderIDEQ(providerID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth by ID", errcode.ErrInternalFailure)
	}
	return auth, nil
}

// ResetFailedLoginAttempts resets the failed login attempts and last failed login time for an OAuthAuth record.
func (o *OAuthAuthRepository) ResetFailedLoginAttempts(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth for resetting failed login attempts", errcode.ErrInternalFailure)
	}

	auth, err = o.client.OAuthAuth.UpdateOne(auth).
		SetFailedLoginAttempts(0).
		SetLastFailedLoginAt(time.Time{}).
		Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to reset failed login attempts", errcode.ErrInternalFailure)
	}

	return auth, nil
}

// SetEmail updates the email address for an OAuthAuth record.
func (o *OAuthAuthRepository) SetEmail(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID, email string) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth for updating email", errcode.ErrInternalFailure)
	}

	auth, err = o.client.OAuthAuth.UpdateOne(auth).
		SetEmail(email).
		Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to update OAuthAuth email", errcode.ErrInternalFailure)
	}

	return auth, nil
}

// SetIsVerified updates the verification status for an OAuthAuth record.
func (o *OAuthAuthRepository) SetIsVerified(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID, isVerified bool) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Query().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get OAuthAuth for updating verification status", errcode.ErrInternalFailure)
	}

	auth, err = o.client.OAuthAuth.UpdateOne(auth).
		SetIsVerified(isVerified).
		Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to update OAuthAuth isVerified status", errcode.ErrInternalFailure)
	}

	return auth, nil
}

// CreateOAuthAuth creates a new OAuthAuth record in the database.
func (o *OAuthAuthRepository) CreateOAuthAuth(ctx context.Context, input *dbmodels.CreateOAuthAuthInput) (*ent.OAuthAuth, error) {
	auth, err := o.client.OAuthAuth.Create().
		SetProvider(input.Provider).
		SetProviderID(input.ProviderID).
		SetAuthAccountID(input.AccountID).
		SetEmail(input.Email).
		SetIsVerified(input.IsVerified).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			upgradedErr := errors.Upgrade(err, errcode.ErrConflict, "User already exists with this provider ID")
			return nil, errors.Join(upgradedErr, "Failed to create local auth record")
		}
		upgradedErr := errors.Upgrade(err, errcode.ErrInternalFailure, "Failed to create oauth record")
		return nil, errors.Join(upgradedErr, "Failed to create OAuthAuth record")
	}
	return auth, nil
}

// DeleteOAuthAuth deletes an OAuthAuth record by provider and auth account ID.
func (o *OAuthAuthRepository) DeleteOAuthAuth(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) error {
	_, err := o.client.OAuthAuth.Delete().
		Where(oauthauth.ProviderEQ(provider)).
		Where(oauthauth.AuthAccountIDEQ(authAccountID)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("OAuthAuth not found", "NotFound", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete OAuthAuth record", errcode.ErrInternalFailure)
	}
	return nil
}

// NewOAuthAuthRepository creates a new instance of OAuthAuthRepository.
func NewOAuthAuthRepository(client *ent.Client) *OAuthAuthRepository {
	return &OAuthAuthRepository{
		client: client,
	}
}
