package dbdomain

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthauth"
	dbmodels "mandacode.com/accounts/auth/internal/domain/models/database"
)

type OAuthAuthRepository interface {
	// GetOAuthAuthByProviderID retrieves an OAuthAuth record by its provider and provider ID.
	GetOAuthAuthByProviderID(ctx context.Context, provider oauthauth.Provider, providerID string) (*ent.OAuthAuth, error)

	// GetOAuthAuthByAccountID retrieves an OAuthAuth record by its account ID.
	GetOAuthAuthByAuthAccountID(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error)

	// CreateOAuthAuth creates a new OAuthAuth record.
	CreateOAuthAuth(ctx context.Context, input *dbmodels.CreateOAuthAuthInput) (*ent.OAuthAuth, error)

	// SetEmail updates the email of an OAuthAuth record by its account ID.
	SetEmail(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID, email string) (*ent.OAuthAuth, error)

	// SetIsActive updates the active status of an OAuthAuth record by its account ID.
	SetIsActive(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID, isActive bool) (*ent.OAuthAuth, error)

	// SetIsVerified updates the verification status of an OAuthAuth record by its account ID.
	SetIsVerified(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID, isVerified bool) (*ent.OAuthAuth, error)

	// OnLoginSuccess handles the logic for a successful login attempt, updating the last login time and resetting failed attempts.
	OnLoginSuccess(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error)

	// OnLoginFailed handles the logic for a failed login attempt, updating the last failed login time and incrementing the failed attempts.
	OnLoginFailed(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error)

	// ResetFailedLoginAttempts resets the failed login attempts for an OAuthAuth record by its account ID.
	ResetFailedLoginAttempts(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) (*ent.OAuthAuth, error)

	// DeleteOAuthAuthByProviderID deletes an OAuthAuth record by its provider and provider ID.
	DeleteOAuthAuthByProviderID(ctx context.Context, provider oauthauth.Provider, providerID string) error

	// DeleteOAuthAuth deletes an OAuthAuth record by its ID.
	DeleteOAuthAuth(ctx context.Context, provider oauthauth.Provider, authAccountID uuid.UUID) error
}
