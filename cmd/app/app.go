package app

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	mux  chi.Router
	pool *pgxpool.Pool
}

func NewServer(mux chi.Router, pool *pgxpool.Pool) *Server {
	return &Server{mux: mux, pool: pool}
}
