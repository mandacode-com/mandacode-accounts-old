package usereventdomain

type UserEventUsecase interface {
	// DeleteUser deletes a user by their ID.
	DeleteUser(userID string) error

	// SetActiveStatus sets the active status of a user by their ID.
	SetActiveStatus(userID string, active bool) error

	// SetEmailVerifiedStatus sets the email verified status of a user by their ID.
	SetEmailVerifiedStatus(userID string, verified bool) error
}
