package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/user/ent"
	"mandacode.com/accounts/user/ent/user"
	repodomain "mandacode.com/accounts/user/internal/domain/port/repository"
)

type UserRepository struct {
	db *ent.Client
}

// UpdateEmailVerificationCode implements repodomain.UserRepository.
func (u *UserRepository) UpdateEmailVerificationCode(userID uuid.UUID, emailVerificationCode string) (*ent.User, error) {
	_, err := u.db.User.
		UpdateOneID(userID).
		SetEmailVerificationCode(emailVerificationCode).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	user, err := u.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser implements repodomain.UserRepository.
func (u *UserRepository) CreateUser(userID uuid.UUID, syncCode string) (*ent.User, error) {
	create := u.db.User.Create().
		SetID(userID).
		SetSyncCode(syncCode).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	user, err := create.Save(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser implements repodomain.UserRepository.
func (u *UserRepository) DeleteUser(userID uuid.UUID) error {
	_, err := u.db.User.
		Delete().
		Where(user.IDEQ(userID)).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// GetUser implements repodomain.UserRepository.
func (u *UserRepository) GetUser(userID uuid.UUID) (*ent.User, error) {
	user, err := u.db.User.
		Query().
		Where(user.IDEQ(userID)).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ArchiveUser implements repodomain.UserRepository.
func (u *UserRepository) ArchiveUser(userID uuid.UUID, syncCode string) error {
	_, err := u.db.User.
		UpdateOneID(userID).
		SetDeletedAt(time.Now()).
		SetSyncCode(syncCode).
		Save(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *ent.Client) repodomain.UserRepository {
	return &UserRepository{db: db}
}
