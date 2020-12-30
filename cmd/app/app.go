package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework1/cmd/app/dto"
	"github.com/lozovoya/agohomework1/cmd/app/md"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

var (
	ErrEmptyLogin    = errors.New("login is empty")
	ErrEmptyPassword = errors.New("password is empty")
	ErrParse         = errors.New("body parse error")
	ErrServer        = errors.New("internal server error")
	ErrWrongUser     = errors.New("username already exist")
	ErrWrongCred     = errors.New("wrong credentials")
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

	identMD := md.IdentMD
	logMD := middleware.Logger

	s.mux.With(logMD).Post("/api/users", s.AddUser)
	s.mux.With(logMD).Post("/api/token", s.Token)
	s.mux.With(logMD, identMD, md.AuthMD(s.ctx, s.pool)).Post("/api/getcards", s.GetCards)

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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return
	}

	var id int
	err = s.pool.QueryRow(s.ctx,
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
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (s *Server) Token(w http.ResponseWriter, r *http.Request) {

	var user *dto.TokenRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(ErrParse)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var hash *[]byte
	err = s.pool.QueryRow(s.ctx, "SELECT password FROM USERS WHERE login=$1", user.Login).Scan(&hash)
	if err != nil {
		log.Println(ErrServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(*hash, []byte(user.Password))
	if err != nil {
		log.Println(ErrWrongCred)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	token, err := uuid.NewRandom()
	if err != nil {
		log.Println(ErrServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = s.pool.Query(s.ctx, "UPDATE users SET token=$1 WHERE login=$2", token.String(), user.Login)
	if err != nil {
		log.Println(ErrServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(dto.TokenDTO{Token: token.String()})
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) GetCards(w http.ResponseWriter, r *http.Request) {

	userid, ok := r.Context().Value(md.UserIdContextKey).(int)
	if !ok {
		log.Println(ok)
		log.Println(ErrServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var roles []string
	err := s.pool.QueryRow(s.ctx, "SELECT roles FROM users WHERE id=$1", userid).Scan(&roles)
	if err != nil {
		log.Println(ErrServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(roles[0], "ADMIN") {
		rows, err := s.pool.Query(s.ctx, "SELECT number, balance FROM cards")
		if err != nil {
			log.Println(ErrServer)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var cards []*dto.CardDTO
		for rows.Next() {
			card := &dto.CardDTO{}
			err = rows.Scan(&card.Number, &card.Balance)
			if err != nil {
				log.Println(ErrServer)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cards = append(cards, card)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err = json.NewEncoder(w).Encode(cards)
		if err != nil {
			log.Println(err)
			return
		}
		return

	} else {
		rows, err := s.pool.Query(s.ctx, "SELECT number, balance FROM cards WHERE owner=$1", userid)
		if err != nil {
			log.Println(ErrServer)
			return
		}
		defer rows.Close()

		var cards []*dto.CardDTO
		for rows.Next() {
			card := &dto.CardDTO{}
			err = rows.Scan(&card.Number, &card.Balance)
			if err != nil {
				log.Println(err)
				return
			}
			cards = append(cards, card)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err = json.NewEncoder(w).Encode(cards)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

}
