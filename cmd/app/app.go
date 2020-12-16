package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework1.git/cmd/app/dto"
	"log"
	"net/http"
)

var (
	ErrEmptyLogin    = errors.New("Login is empty")
	ErrEmptyPassword = errors.New("Password is empty")
	ErrParse         = errors.New("Body parse error")
)

type Server struct {
	mux  chi.Router
	pool *pgxpool.Pool
}

func NewServer(mux chi.Router, pool *pgxpool.Pool) *Server {
	return &Server{mux: mux, pool: pool}
}

func (s *Server) Init() error {
	s.mux.With(middleware.Logger).Put("/api/users", s.AddUser)
	return nil
}

func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {

	var user *dto.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(ErrParse)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(user)
	return
}
