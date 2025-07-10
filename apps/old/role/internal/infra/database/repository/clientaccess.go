package repository

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/role/ent"
	"mandacode.com/accounts/role/ent/clientaccess"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type ClientAccessRepository struct {
	db *ent.Client
}

// CreateClientAccess implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) CreateClientAccess(serviceID uuid.UUID, clientID string, clientSecret string, isActive *bool) (*ent.ClientAccess, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)

	create := c.db.ClientAccess.Create().
		SetServiceID(serviceID).
		SetClientID(clientID).
		SetClientSecret(string(hashedSecret))

	if isActive != nil {
		create.SetIsActive(*isActive)
	}

	clientAccess, err := create.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return clientAccess, nil
}

// DeleteClientAccess implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) DeleteClientAccess(id uuid.UUID) error {
	err := c.db.ClientAccess.
		DeleteOneID(id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// DeleteClientAccessByServiceID implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) DeleteClientAccessByServiceID(serviceID uuid.UUID) error {
	_, err := c.db.ClientAccess.
		Delete().
		Where(clientaccess.ServiceID(serviceID)).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// GetClientAccess implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) GetClientAccess(serviceID uuid.UUID, clientID string) (*ent.ClientAccess, error) {
	clientAccess, err := c.db.ClientAccess.
		Query().
		Where(
			clientaccess.ServiceID(serviceID),
			clientaccess.ClientID(clientID),
		).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return clientAccess, nil
}

// GetClientAccessByID implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) GetClientAccessByID(id uuid.UUID) (*ent.ClientAccess, error) {
	clientAccess, err := c.db.ClientAccess.
		Query().
		Where(clientaccess.ID(id)).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return clientAccess, nil
}

// GetClientAccessesByServiceID implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) GetClientAccessesByServiceID(serviceID uuid.UUID) ([]*ent.ClientAccess, error) {
	clientAccesses, err := c.db.ClientAccess.
		Query().
		Where(clientaccess.ServiceID(serviceID)).
		All(context.Background())

	if err != nil {
		return nil, err
	}

	return clientAccesses, nil
}

// UpdateClientAccess implements repodomain.ClientAccessRepository.
func (c *ClientAccessRepository) UpdateClientAccess(id uuid.UUID, serviceID *uuid.UUID, clientID *string, clientSecret *string, isActive *bool) (*ent.ClientAccess, error) {
	update := c.db.ClientAccess.
		UpdateOneID(id)

	if serviceID != nil {
		update.SetServiceID(*serviceID)
	}

	if clientID != nil {
		update.SetClientID(*clientID)
	}

	if clientSecret != nil {
		hashedSecret, err := bcrypt.GenerateFromPassword([]byte(*clientSecret), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		update.SetClientSecret(string(hashedSecret))
	}

	if isActive != nil {
		update.SetIsActive(*isActive)
	}

	clientAccess, err := update.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return clientAccess, nil
}

// NewClientAccessRepository creates a new ClientAccessRepository.
func NewClientAccessRepository(db *ent.Client) repodomain.ClientAccessRepository {
	return &ClientAccessRepository{db: db}
}
