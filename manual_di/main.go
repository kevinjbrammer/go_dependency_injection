package main

import (
	"database/sql"
	"log"
	"spaceships/config"
	"spaceships/repository"
	"spaceships/server"
	"spaceships/service"

	_ "github.com/mattn/go-sqlite3"
)

func openDatabaseConnection(cfg config.IConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", cfg.GetDatabasePath())
}

func main() {
	var cfg config.IConfig = config.NewConfig("./spaceships.db", "8000", true)

	db, err := openDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repository := repository.NewSpaceshipRepository(db)
	service := service.NewSpaceshipService(cfg, repository)
	server := server.NewServer(cfg, service)
	server.Run()
}
