package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"spaceships/config"
	"spaceships/service"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	config  *config.Config
	service *service.SpaceshipService
}

// NewServer constructs a Server
func NewServer() *Server {
	cfg := config.NewConfig()
	service := service.NewSpaceshipService()
	return &Server{
		config:  cfg,
		service: service,
	}
}

// Run starts the server
func (s *Server) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/spaceships", s.getAllSpaceships).Methods("GET")
	router.HandleFunc("/spaceships/{id}", s.getSpaceshipByID).Methods("GET")

	address := ":" + s.config.Port

	fmt.Printf("Server is running on port %s\n", s.config.Port)
	log.Fatal(http.ListenAndServe(address, router))
}

func (s *Server) getAllSpaceships(w http.ResponseWriter, r *http.Request) {
	spaceships, err := s.service.GetAllSpaceships()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spaceships)
}

func (s *Server) getSpaceshipByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	spaceship, err := s.service.GetSpaceshipByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spaceship)
}
