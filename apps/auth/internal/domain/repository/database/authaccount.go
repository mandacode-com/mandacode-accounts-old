package dbdomain

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	dbmodels "mandacode.com/accounts/auth/internal/domain/models/database"
)

type AuthAccountRepository interface {
	// GetAuthAccountByID retrieves an authentication account by its ID.
	GetAuthAccountByID(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error)

	// GetAuthAccountByUserID retrieves an authentication account by the associated user ID.
	GetAuthAccountByUserID(ctx context.Context, userID uuid.UUID) (*ent.AuthAccount, error)

	// CreateAuthAccount creates a new authentication account.
	CreateAuthAccount(ctx context.Context, account *dbmodels.CreateAuthAccountInput) (*ent.AuthAccount, error)

	// UpdateAuthAccount updates an existing authentication account.
	UpdateAuthAccount(ctx context.Context, id uuid.UUID, account *dbmodels.UpdateAuthAccountInput) (*ent.AuthAccount, error)

	// SetActiveStatus sets the active status of an authentication account.
	SetActiveStatus(ctx context.Context, id uuid.UUID, isActive bool) (*ent.AuthAccount, error)

	// OnLoginSuccess handles the logic for a successful login attempt, updating the last login time and resetting failed attempts.
	OnLoginSuccess(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error)

	// OnLoginFailed handles the logic for a failed login attempt, updating the last failed login time and incrementing the failed attempts.
	OnLoginFailed(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error)

	// ResetFailedLoginAttempts resets the failed login attempts for an authentication account.
	ResetFailedLoginAttempts(ctx context.Context, id uuid.UUID) (*ent.AuthAccount, error)

	// DeleteAuthAccount deletes an authentication account by its ID.
	DeleteAuthAccount(ctx context.Context, id uuid.UUID) error
}
