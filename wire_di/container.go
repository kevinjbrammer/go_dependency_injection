//+build wireinject

package main

import (
	"spaceships/config"
	"spaceships/repository"
	"spaceships/server"
	"spaceships/service"

	"github.com/google/wire"
)

func buildServer(config config.IConfig) (*server.Server, error) {
	wire.Build(
		openDatabaseConnection,

		repository.NewSpaceshipRepository,
		wire.Bind(new(repository.ISpaceshipRepository), new(*repository.SpaceshipRepository)),

		service.NewSpaceshipService,
		wire.Bind(new(service.ISpaceshipService), new(*service.SpaceshipService)),

		server.NewServer)

	return &server.Server{}, nil
}
