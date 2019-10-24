package main

import (
	"database/sql"
	"log"
	"spaceships/config"

	_ "github.com/mattn/go-sqlite3"
)

func openDatabaseConnection(cfg config.IConfig) (*sql.DB, error) {
	return sql.Open("sqlite3", cfg.GetDatabasePath())
}

func main() {
	cfg := config.NewConfig("./spaceships.db", "8000", true)
	server, err := buildServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
