package main

import (
	"database/sql"
	"log"
	"spaceships/config"
	"spaceships/repository"
	"spaceships/server"
	"spaceships/service"

	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/dig"
)

func openDatabaseConnection(cfg config.IConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", cfg.GetDatabasePath())
}

func buildContainer(databasePath, serverPort string, serviceEnabled bool) *dig.Container {
	container := dig.New()

	//Note: dig has the dig.As provider option that is slated for
	//next release which will make it so that you do not have
	//to wrap your constructors with functions that return the
	//appropriate interface.

	container.Provide(func() config.IConfig {
		return config.NewConfig(databasePath, serverPort, serviceEnabled)
	})

	container.Provide(openDatabaseConnection)
	container.Provide(func(db *sql.DB) repository.ISpaceshipRepository {
		return repository.NewSpaceshipRepository(db)
	})
	container.Provide(func(cfg config.IConfig, rep repository.ISpaceshipRepository) service.ISpaceshipService {
		return service.NewSpaceshipService(cfg, rep)
	})
	container.Provide(server.NewServer)

	return container
}

func main() {
	container := buildContainer("./spaceships.db", "8000", true)

	err := container.Invoke(func(server *server.Server) {
		server.Run()
	})

	if err != nil {
		log.Fatal(err)
	}
}
