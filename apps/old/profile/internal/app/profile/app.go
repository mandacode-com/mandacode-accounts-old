package profile

import (
	"errors"

	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/domain/model"
	repodomain "mandacode.com/accounts/profile/internal/domain/repository"
)

type profileApp struct {
	repo repodomain.ProfileRepository
}

func NewProfileApp(repo repodomain.ProfileRepository) ProfileApp {
	return &profileApp{
		repo: repo,
	}
}

// UpdateProfile implements ProfileApp.
func (p *profileApp) UpdateProfile(userID uuid.UUID, email *string, displayName *string, bio *string, avatarURL *string) (*model.Profile, error) {
	if email == nil && displayName == nil && bio == nil && avatarURL == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	entProfile, err := p.repo.UpdateProfile(userID, email, displayName, bio, avatarURL)
	if err != nil {
		return nil, err
	}

	profile := model.NewProfileFromEnt(entProfile)

	return profile, nil
}

// InitializeProfile implements ProfileApp.
func (p *profileApp) InitializeProfile(userID uuid.UUID) (*model.Profile, error) {
	entProfile, err := p.repo.InitializeProfile(userID)
	if err != nil {
		return nil, err
	}

	profile := model.NewProfileFromEnt(entProfile)

	return profile, nil
}

// GetProfile implements ProfileApp.
func (p *profileApp) GetProfile(userID uuid.UUID) (*model.Profile, error) {
	entProfile, err := p.repo.GetProfileByID(userID)
	if err != nil {
		return nil, err
	}

	profile := model.NewProfileFromEnt(entProfile)

	return profile, nil
}

// DeleteProfile implements ProfileApp.
func (p *profileApp) DeleteProfile(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("userID cannot be nil")
	}

	err := p.repo.DeleteProfile(userID)
	if err != nil {
		return err
	}

	return nil
}
