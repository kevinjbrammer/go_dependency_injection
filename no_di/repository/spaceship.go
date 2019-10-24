package repository

import (
	"database/sql"
	"log"
	"spaceships/config"
	"spaceships/model"

	_ "github.com/mattn/go-sqlite3"
)

// SpaceshipRepository ...
type SpaceshipRepository struct {
	database *sql.DB
}

// NewSpaceshipRepository constructs a SpaceshipRepository struct
func NewSpaceshipRepository() *SpaceshipRepository {
	cfg := config.NewConfig()
	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	return &SpaceshipRepository{db}
}

// GetAllSpaceships retrieves all spaceships from the database
func (repository *SpaceshipRepository) GetAllSpaceships() ([]*model.Spaceship, error) {
	rows, err := repository.database.Query("SELECT id, name FROM spaceships;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	spaceships := []*model.Spaceship{}
	for rows.Next() {
		var (
			id   int
			name string
		)

		rows.Scan(&id, &name)

		spaceships = append(spaceships, &model.Spaceship{
			ID:   id,
			Name: name,
		})
	}
	return spaceships, nil
}

// GetSpaceshipByID retrieves a spaceship from the database by ID
func (repository *SpaceshipRepository) GetSpaceshipByID(id int) (*model.Spaceship, error) {
	var spaceship model.Spaceship
	row := repository.database.QueryRow("SELECT id, name FROM spaceships WHERE id=$1;", id)

	err := row.Scan(&spaceship.ID, &spaceship.Name)
	if err != nil {
		return nil, err
	}

	return &spaceship, nil
}
