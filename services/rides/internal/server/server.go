package server

import (
	"net/http"
	"rides/internal/database"
	"rides/internal/services"
)

type Server struct {
	db             *database.Database
	userService    *services.UserService
	paymentService *services.PaymentService
}

func NewServer(db *database.Database, userService *services.UserService, paymentService *services.PaymentService) *Server {
	return &Server{db: db, userService: userService, paymentService: paymentService}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /rides", s.createRide)
	mux.HandleFunc("GET /rides/{id}", s.getRide)
	mux.HandleFunc("PATCH /rides/{id}/status", s.updateRideStatus)

	mux.ServeHTTP(w, r)
}
