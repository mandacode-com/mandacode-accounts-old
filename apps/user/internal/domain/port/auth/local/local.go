package localuserdomain

import (
	"context"

	"github.com/google/uuid"
	localuserv1 "mandacode.com/accounts/proto/auth/user/local/v1"
)

type LocalUserService interface {
	// GetUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to retrieve.
	GetUser(ctx context.Context, userID uuid.UUID) (*localuserv1.GetUserResponse, error)


	// EnrollUser
	//
	// Parameters:
	//	 - ctx: The context for the operation.
	//   - userID: The ID of the user to enroll.
	//   - email: The email address of the user to enroll.
	//	 - password: The password for the user to enroll.
	//   - isActive: Indicates if the user is active.
	//   - isVerified: Indicates if the user is verified.
	EnrollUser(ctx context.Context, userID uuid.UUID, email, password string, isActive, isVerified bool) (*localuserv1.EnrollUserResponse, error)

	// UpdateEmail
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose email is to be updated.
	//	 - email: The new email address for the user.
	UpdateEmail(ctx context.Context, userID uuid.UUID, email string) (*localuserv1.UpdateEmailResponse, error)

	// DeleteUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to be deleted.
	DeleteUser(ctx context.Context, userID uuid.UUID) (*localuserv1.DeleteUserResponse, error)

	// UpdatePassword
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose password is to be updated.
	//   - currentPassword: The current password of the user.
	//   - newPassword: The new password for the user.
	UpdatePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) (*localuserv1.UpdatePasswordResponse, error)

	// UpdateActiveStatus
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose active status is to be updated.
	//	 - isActive: Indicates if the user is active.
	UpdateActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*localuserv1.UpdateActiveStatusResponse, error)

	// UpdateVerifiedStatus
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose verified status is to be updated.
	//   - isVerified: Indicates if the user is verified.
	UpdateVerifiedStatus(ctx context.Context, userID uuid.UUID, isVerified bool) (*localuserv1.UpdateVerifiedStatusResponse, error)
}
