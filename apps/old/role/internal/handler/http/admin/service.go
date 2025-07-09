package adminhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	serviceapp "mandacode.com/accounts/role/internal/app/service"
	adminhandlerdto "mandacode.com/accounts/role/internal/handler/http/admin/dto"
)

type ServiceHandler struct {
	serviceApp      serviceapp.ServiceApp
	validator       *validator.Validate
	adminMiddleware gin.HandlerFunc
}

func NewServiceHandler(
	serviceApp serviceapp.ServiceApp,
	validator *validator.Validate,
	adminMiddleware gin.HandlerFunc,
) (*ServiceHandler, error) {
	return &ServiceHandler{
		serviceApp:      serviceApp,
		validator:       validator,
		adminMiddleware: adminMiddleware,
	}, nil
}

func (h *ServiceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(h.adminMiddleware)
	rg.POST("/", h.CreateService)
	rg.GET("/", h.GetAllServices)
	rg.GET("/:id", h.GetServiceByID)
	rg.PUT("/:id", h.UpdateService)
	rg.DELETE("/:id", h.DeleteService)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var dto adminhandlerdto.CreateServiceRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	service, err := h.serviceApp.CreateService(dto.Name, dto.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create service"})
		return
	}

	response := adminhandlerdto.CreateServiceResponse{
		Service: service,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID format"})
		return
	}

	service, err := h.serviceApp.GetServiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get service"})
		return
	}

	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	response := adminhandlerdto.GetServiceResponse{
		Service: service,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ServiceHandler) GetAllServices(c *gin.Context) {
	services, err := h.serviceApp.GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get services"})
		return
	}

	response := adminhandlerdto.GetAllServicesResponse{
		Services: services,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID format"})
		return
	}

	var dto adminhandlerdto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	service, err := h.serviceApp.UpdateService(id, dto.Name, dto.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update service"})
		return
	}

	response := adminhandlerdto.UpdateServiceResponse{
		Service: service,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID format"})
		return
	}

	err = h.serviceApp.DeleteService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete service"})
		return
	}

	response := adminhandlerdto.DeleteServiceResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
