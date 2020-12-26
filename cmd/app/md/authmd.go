package md

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework1.git/cmd/app/dto"
	"log"
	"net/http"
)

var IdentifierContextKey = &contextKey{"identifier context"}
var UserIdContextKey = &contextKey{"user id"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

func IdentMD(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var token *dto.TokenDTO
		err := json.NewDecoder(r.Body).Decode(&token)
		if err != nil {
			log.Println(err)
			return
		}

		if token.Token != "" {
			ctx := context.WithValue(r.Context(), IdentifierContextKey, &token.Token)
			r = r.WithContext(ctx)
		}

		handler.ServeHTTP(w, r)

	})
}

func AuthMD(dbCtx context.Context, pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, _ := r.Context().Value(IdentifierContextKey).(*string)

			conn, err := pool.Acquire(dbCtx)
			if err != nil {
				log.Println(err)
				return
			}
			if err != nil {
				log.Println(err)
				return
			}
			defer conn.Release()

			var userid int
			err = conn.QueryRow(dbCtx, "SELECT id FROM users WHERE token=$1", token).Scan(&userid)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if userid != 0 {
				ctx := context.WithValue(r.Context(), UserIdContextKey, userid)
				r = r.WithContext(ctx)
			}

			handler.ServeHTTP(w, r)
		})
	}
}
