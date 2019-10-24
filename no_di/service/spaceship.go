package service

import (
	"spaceships/config"
	"spaceships/model"
	"spaceships/repository"
)

// SpaceshipService ...
type SpaceshipService struct {
	config     *config.Config
	repository *repository.SpaceshipRepository
}

// NewSpaceshipService constructs a SpaceshipService
func NewSpaceshipService() *SpaceshipService {
	cfg := config.NewConfig()
	rep := repository.NewSpaceshipRepository()
	return &SpaceshipService{config: cfg, repository: rep}
}

// GetAllSpaceships returns all spaceships or no spaceships if the service is
// not enabled.
func (service *SpaceshipService) GetAllSpaceships() ([]*model.Spaceship, error) {
	if service.config.Enabled {
		return service.repository.GetAllSpaceships()
	}

	return []*model.Spaceship{}, nil
}

// GetSpaceshipByID returns a spaceship by ID or no spaceship if the service is
// not enabled.
func (service *SpaceshipService) GetSpaceshipByID(id int) (*model.Spaceship, error) {
	if service.config.Enabled {
		return service.repository.GetSpaceshipByID(id)
	}

	return &model.Spaceship{}, nil
}
