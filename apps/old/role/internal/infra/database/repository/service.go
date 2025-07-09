package repository

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
	"mandacode.com/accounts/role/ent/service"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type ServiceRepository struct {
	db *ent.Client
}

func NewServiceRepository(db *ent.Client) repodomain.ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

// CreateService implements repodomain.ServiceRepository.
func (s *ServiceRepository) CreateService(name string, description *string) (*ent.Service, error) {
	create := s.db.Service.Create()

	create.SetID(uuid.New())
	create.SetName(name)

	if description != nil {
		create.SetDescription(*description)
	}

	service, err := create.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetServiceByID implements repodomain.ServiceRepository.
func (s *ServiceRepository) GetServiceByID(id uuid.UUID) (*ent.Service, error) {
	service, err := s.db.Service.
		Query().
		Where(service.ID(id)).
		Only(context.Background())

	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetAllServices implements repodomain.ServiceRepository.
func (s *ServiceRepository) GetAllServices() ([]*ent.Service, error) {
	services, err := s.db.Service.
		Query().
		All(context.Background())
	if err != nil {
		return nil, err
	}

	return services, nil
}

// UpdateService implements repodomain.ServiceRepository.
func (s *ServiceRepository) UpdateService(id uuid.UUID, name *string, description *string) (*ent.Service, error) {
	update := s.db.Service.UpdateOneID(id)

	if name != nil {
		update.SetName(*name)
	}

	if description != nil {
		update.SetDescription(*description)
	}

	service, err := update.Save(context.Background())
	if err != nil {
		return nil, err
	}

	return service, nil
}

// DeleteService implements repodomain.ServiceRepository.
func (s *ServiceRepository) DeleteService(id uuid.UUID) error {
	err := s.db.Service.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
