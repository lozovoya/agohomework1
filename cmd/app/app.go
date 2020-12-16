package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type Server struct {
	mux  chi.Router
	pool *pgxpool.Pool
}

func NewServer(mux chi.Router, pool *pgxpool.Pool) *Server {
	return &Server{mux: mux, pool: pool}
}

func (s *Server) Init() error {
	s.mux.With(middleware.Logger).Get("/api/users", s.AddUser)
	return nil
}

func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
