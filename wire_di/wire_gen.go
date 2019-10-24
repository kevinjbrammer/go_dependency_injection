// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"spaceships/config"
	"spaceships/repository"
	"spaceships/server"
	"spaceships/service"
)

import (
	_ "github.com/mattn/go-sqlite3"
)

// Injectors from container.go:

func buildServer(config2 config.IConfig) (*server.Server, error) {
	db, err := openDatabaseConnection(config2)
	if err != nil {
		return nil, err
	}
	spaceshipRepository := repository.NewSpaceshipRepository(db)
	spaceshipService := service.NewSpaceshipService(config2, spaceshipRepository)
	serverServer := server.NewServer(config2, spaceshipService)
	return serverServer, nil
}
