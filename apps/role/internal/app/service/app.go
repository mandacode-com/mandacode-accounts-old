package serviceapp

import (
	"errors"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
	repodomain "mandacode.com/accounts/role/internal/domain/repository"
)

type serviceApp struct {
	serviceRepo repodomain.ServiceRepository
}

// CreateService implements ServiceApp.
func (s *serviceApp) CreateService(name string, description *string) (*model.Service, error) {
	service, err := s.serviceRepo.CreateService(name, description)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, errors.New("failed to create service")
	}

	return model.ServiceFromEnt(service), nil
}

// DeleteService implements ServiceApp.
func (s *serviceApp) DeleteService(id uuid.UUID) error {
	err := s.serviceRepo.DeleteService(id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllServices implements ServiceApp.
func (s *serviceApp) GetAllServices() ([]*model.Service, error) {
	services, err := s.serviceRepo.GetAllServices()
	if err != nil {
		return nil, err
	}

	var result []*model.Service
	for _, service := range services {
		result = append(result, model.ServiceFromEnt(service))
	}

	return result, nil
}

// GetServiceByID implements ServiceApp.
func (s *serviceApp) GetServiceByID(id uuid.UUID) (*model.Service, error) {
	service, err := s.serviceRepo.GetServiceByID(id)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, nil // or return an error if you prefer
	}

	return model.ServiceFromEnt(service), nil
}

// UpdateService implements ServiceApp.
func (s *serviceApp) UpdateService(id uuid.UUID, name *string, description *string) (*model.Service, error) {
	service, err := s.serviceRepo.UpdateService(id, name, description)
	if err != nil {
		return nil, err
	}

	return model.ServiceFromEnt(service), nil
}

// NewServiceApp creates a new ServiceApp.
func NewServiceApp(
	serviceRepo repodomain.ServiceRepository,
) ServiceApp {
	return &serviceApp{
		serviceRepo: serviceRepo,
	}
}
