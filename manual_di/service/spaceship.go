package service

import (
	"spaceships/config"
	"spaceships/model"
	"spaceships/repository"
)

// ISpaceshipService interface
type ISpaceshipService interface {
	GetAllSpaceships() ([]*model.Spaceship, error)
	GetSpaceshipByID(id int) (*model.Spaceship, error)
}

// SpaceshipService ...
type SpaceshipService struct {
	config     config.IConfig
	repository repository.ISpaceshipRepository
}

// NewSpaceshipService constructs a SpaceshipService
func NewSpaceshipService(cfg config.IConfig, rep repository.ISpaceshipRepository) *SpaceshipService {
	return &SpaceshipService{config: cfg, repository: rep}
}

// GetAllSpaceships returns all spaceships or no spaceships of the service is
// not enabled.
func (service *SpaceshipService) GetAllSpaceships() ([]*model.Spaceship, error) {
	if service.config.ServiceEnabled() {
		return service.repository.GetAllSpaceships()
	}

	return []*model.Spaceship{}, nil
}

// GetSpaceshipByID returns a spaceship by ID or no spaceship if the service is
// not enabled.
func (service *SpaceshipService) GetSpaceshipByID(id int) (*model.Spaceship, error) {
	if service.config.ServiceEnabled() {
		return service.repository.GetSpaceshipByID(id)
	}

	return &model.Spaceship{}, nil
}
