package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework1.git/cmd/app/dto"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

var (
	ErrEmptyLogin    = errors.New("Login is empty")
	ErrEmptyPassword = errors.New("Password is empty")
	ErrParse         = errors.New("Body parse error")
	ErrServer        = errors.New("Internal server error")
	ErrWrongUser     = errors.New("Username already exist")
)

type Server struct {
	mux  chi.Router
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewServer(mux chi.Router, pool *pgxpool.Pool, ctx context.Context) *Server {
	return &Server{mux: mux, pool: pool, ctx: ctx}
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

	if user.Login == "" {
		log.Println(ErrEmptyLogin)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		log.Println(ErrEmptyPassword)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := s.pool.Acquire(s.ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Release()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return
	}

	var id int
	err = conn.QueryRow(s.ctx,
		"INSERT INTO USERS (login, password, roles) VALUES($1, $2, $3) ON CONFLICT DO NOTHING  RETURNING id",
		user.Login, hash, user.Roles).Scan(&id)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "23505") {
			w.WriteHeader(http.StatusConflict)
			err = json.NewEncoder(w).Encode(dto.ErrResp{ErrWrongUser.Error()})
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrResp{ErrServer.Error()})
		if err != nil {
			log.Println(err)
		}
		return
	}

	response := dto.UserId{Id: id}
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
