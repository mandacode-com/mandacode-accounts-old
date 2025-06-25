package authrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
)

type OAuthAuthRepository struct {
	db *ent.Client
}

func NewOAuthUserRepository(db *ent.Client) repodomain.OAuthUserRepository {
	return &OAuthAuthRepository{db: db}
}

func (r *OAuthAuthRepository) GetUserByProvider(provider oauthuser.Provider, providerID string) (*ent.OAuthUser, error) {
	user, err := r.db.OAuthUser.
		Query().
		Where(
			oauthuser.ProviderEQ(provider),
			oauthuser.ProviderID(providerID),
		).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *OAuthAuthRepository) CreateOAuthUser(userID uuid.UUID, provider oauthuser.Provider, providerID string, email string, isActive *bool, isVerified *bool) (*ent.OAuthUser, error) {
	create := r.db.OAuthUser.Create()

	create.SetID(userID)
	create.SetProvider(provider)
	create.SetProviderID(providerID)
	create.SetEmail(email)
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

func (r *OAuthAuthRepository) DeleteOAuthUser(userID uuid.UUID) error {
	_, err := r.db.OAuthUser.Delete().Where(
		oauthuser.ID(userID),
	).Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *OAuthAuthRepository) DeleteOAuthUserByProvider(userID uuid.UUID, provider oauthuser.Provider) error {
	user, err := r.db.OAuthUser.Query().Where(
		oauthuser.ID(userID),
		oauthuser.ProviderEQ(provider),
	).Only(context.Background())

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}
	// Delete the user by ID and provider
	err = r.db.OAuthUser.DeleteOne(user).Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *OAuthAuthRepository) UpdateOAuthUser(userID uuid.UUID, provider *oauthuser.Provider, providerID *string, email *string, isActive *bool, isVerified *bool) (*ent.OAuthUser, error) {
	builder := r.db.OAuthUser.UpdateOneID(userID)

	if provider != nil {
		builder.SetProvider(*provider)
	}
	if providerID != nil {
		builder.SetProviderID(*providerID)
	}
	if email != nil {
		builder.SetEmail(*email)
	}
	if isActive != nil {
		builder.SetIsActive(*isActive)
	}
	if isVerified != nil {
		builder.SetIsVerified(*isVerified)
	}

	user, err := builder.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}
