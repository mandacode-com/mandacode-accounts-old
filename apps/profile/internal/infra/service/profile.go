package service

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/domain/dto"
	repodomain "mandacode.com/accounts/profile/internal/domain/repository"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
)

type ProfileService struct {
	repo repodomain.ProfileRepository
}

func NewProfileService(repo repodomain.ProfileRepository) svcdomain.ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfileByID(userID uuid.UUID) (*dto.Profile, error) {
	profile, err := s.repo.GetProfileByID(userID)
	if err != nil {
		return nil, err
	}

	dtoProfile := dto.NewProfile(
		profile.ID,
		&profile.Email,
		&profile.DisplayName,
		&profile.Bio,
		&profile.AvatarURL,
		profile.CreatedAt,
		profile.UpdatedAt,
	)
	err = dtoProfile.Validate()
	if err != nil {
		return nil, err
	}

	return dtoProfile, nil
}

func (s *ProfileService) InitializeProfile(userID uuid.UUID) (*dto.Profile, error) {
	profile, err := s.repo.InitializeProfile(userID)

	if err != nil {
		return nil, err
	}

	dtoProfile := dto.NewProfile(
		profile.ID,
		&profile.Email,
		&profile.DisplayName,
		&profile.Bio,
		&profile.AvatarURL,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	err = dtoProfile.Validate()
	if err != nil {
		return nil, err
	}

	return dtoProfile, nil
}

func (s *ProfileService) UpdateProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string) (*dto.Profile, error) {
	profile, err := s.repo.UpdateProfile(userID, email, displayName, bio, avatarURL)
	if err != nil {
		return nil, err
	}

	dtoProfile := dto.NewProfile(
		profile.ID,
		&profile.Email,
		&profile.DisplayName,
		&profile.Bio,
		&profile.AvatarURL,
		profile.CreatedAt,
		profile.UpdatedAt,
	)
	err = dtoProfile.Validate()
	if err != nil {
		return nil, err
	}

	return dtoProfile, nil
}

func (s *ProfileService) DeleteProfile(userID uuid.UUID) error {
	err := s.repo.DeleteProfile(userID)
	if err != nil {
		return err
	}
	return nil
}
