package adminhandlerdto

import "mandacode.com/accounts/role/internal/domain/model"

type CreateServiceRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}
type CreateServiceResponse struct {
	*model.Service // Embedding model.Service to include service details in the response
}

type GetServiceResponse struct {
	*model.Service // Embedding model.Service to include service details in the response
}

type GetAllServicesResponse struct {
	Services []*model.Service `json:"services"` // Slice of model.Service to hold multiple service details
}

type UpdateServiceRequest struct {
	Name        *string `json:"name,omitempty"`        // Optional field, can be nil
	Description *string `json:"description,omitempty"` // Optional field, can be nil
}
type UpdateServiceResponse struct {
	*model.Service // Embedding model.Service to include service details in the response
}

type DeleteServiceResponse struct {
	Success bool `json:"success"` // Indicates whether the deletion was successful
}
