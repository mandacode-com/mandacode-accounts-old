package repository

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/profile/ent"
	"mandacode.com/accounts/profile/ent/profile"
	repodomain "mandacode.com/accounts/profile/internal/domain/repository"
)

type ProfileRepository struct {
	db *ent.Client
}

func NewProfileRepository(db *ent.Client) repodomain.ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetProfileByID(userID uuid.UUID) (*ent.Profile, error) {
	profile, err := r.db.Profile.
		Query().
		Where(profile.IDEQ(userID)).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) InitializeProfile(userID uuid.UUID) (*ent.Profile, error) {
	profile, err := r.db.Profile.Create().
		SetID(userID).
		Save(context.Background())

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) UpdateProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string) (*ent.Profile, error) {
	update := r.db.Profile.UpdateOneID(userID)

	if email != nil {
		update.SetEmail(*email)
	}
	if displayName != nil {
		update.SetDisplayName(*displayName)
	}
	if bio != nil {
		update.SetNillableBio(bio)
	}
	if avatarURL != nil {
		update.SetNillableAvatarURL(avatarURL)
	}

	profile, err := update.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) DeleteProfile(userID uuid.UUID) error {
	err := r.db.Profile.DeleteOneID(userID).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
