package api

import (
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/go-chi/chi"
)

type Server struct {
	store  db.Store
	router chi.Router
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := chi.NewRouter()

	server.router = router
	return server
}

func (server *Server) Run(addr string) error {
	return http.ListenAndServe(addr, server.router)
}
