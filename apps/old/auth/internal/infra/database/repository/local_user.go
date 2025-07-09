package repository

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/localuser"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
)

type LocalUserRepository struct {
	db *ent.Client
}

func NewLocalUserRepository(db *ent.Client) repodomain.LocalUserRepository {
	return &LocalUserRepository{db: db}
}

func (r *LocalUserRepository) GetUserByEmail(email string) (*ent.LocalUser, error) {
	user, err := r.db.LocalUser.
		Query().
		Where(localuser.EmailEQ(email)).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *LocalUserRepository) GetUserByID(userID uuid.UUID) (*ent.LocalUser, error) {
	user, err := r.db.LocalUser.
		Query().
		Where(localuser.IDEQ(userID)).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *LocalUserRepository) CreateUser(userID uuid.UUID, email string, password string, isActive *bool, isVerified *bool) (*ent.LocalUser, error) {
	create := r.db.LocalUser.Create()

	create.SetID(userID)
	create.SetEmail(email)
	create.SetPassword(password)
	if isActive != nil {
		create.SetIsActive(*isActive)
	}
	if isVerified != nil {
		create.SetIsVerified(*isVerified)
	}
	user, err := create.Save(context.Background())

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *LocalUserRepository) DeleteUser(userID uuid.UUID) error {
	err := r.db.LocalUser.DeleteOneID(userID).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *LocalUserRepository) UpdateUser(userID uuid.UUID, email *string, password *string, isActive *bool, isVerified *bool) (*ent.LocalUser, error) {
	update := r.db.LocalUser.UpdateOneID(userID)

	if email != nil {
		update.SetEmail(*email)
	}
	if password != nil {
		update.SetPassword(*password)
	}
	if isActive != nil {
		update.SetIsActive(*isActive)
	}
	if isVerified != nil {
		update.SetIsVerified(*isVerified)
	}

	user, err := update.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}
