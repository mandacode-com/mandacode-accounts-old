package dbdomain

import (
	"context"

	"github.com/google/uuid"
	dbmodels "mandacode.com/accounts/auth/internal/domain/models/database"
)

type LocalAuthRepository interface {
	// GetLocalAuthByEmail retrieves a LocalAuth record by email.
	GetLocalAuthByEmail(ctx context.Context, email string) (*dbmodels.SecureLocalAuth, error)

	// GetLocalAuthByID retrieves a LocalAuth record by its ID.
	GetLocalAuthByID(ctx context.Context, id uuid.UUID) (*dbmodels.SecureLocalAuth, error)

	// GetLocalAuthByAuthAccountID retrieves a LocalAuth record by the associated AuthAccount ID.
	GetLocalAuthByAuthAccountID(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error)

	// ComparePassword compares a plaintext password with the hashed password of a LocalAuth record.
	ComparePassword(ctx context.Context, authAccountID uuid.UUID, password string) (bool, error)

	// CreateLocalAuth creates a new LocalAuth record.
	CreateLocalAuth(ctx context.Context, input *dbmodels.CreateLocalAuthInput) (*dbmodels.SecureLocalAuth, error)

	// SetEmail updates the email of a LocalAuth record by its account ID.
	SetEmail(ctx context.Context, authAccountID uuid.UUID, email string) (*dbmodels.SecureLocalAuth, error)

	// SetPassword updates the password of a LocalAuth record by its account ID.
	SetPassword(ctx context.Context, authAccountID uuid.UUID, password string) (*dbmodels.SecureLocalAuth, error)

	// SetIsVerified updates the verification status of a LocalAuth record by its account ID.
	SetIsVerified(ctx context.Context, authAccountID uuid.UUID, isVerified bool) (*dbmodels.SecureLocalAuth, error)

	// OnLoginSuccess handles the logic for a successful login attempt, updating the last login time and resetting failed attempts.
	OnLoginSuccess(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error)

	// OnLoginFailed handles the logic for a failed login attempt, updating the last failed login time and incrementing the failed attempts.
	OnLoginFailed(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error)

	// ResetFailedLoginAttempts resets the failed login attempts for a LocalAuth record by its account ID.
	ResetFailedLoginAttempts(ctx context.Context, authAccountID uuid.UUID) (*dbmodels.SecureLocalAuth, error)

	// DeleteLocalAuth deletes a LocalAuth record by its ID.
	DeleteLocalAuth(ctx context.Context, id uuid.UUID) error

	// DeleteLocalAuthByAuthAccountID deletes a LocalAuth record by its associated AuthAccount ID.
	DeleteLocalAuthByAuthAccountID(ctx context.Context, authAccountID uuid.UUID) error
}
