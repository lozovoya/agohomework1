package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const dbcon = "postgres://app:pass@localhost:5432/db"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port), dbcon); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, dbcon string) error {

	r := chi.NewRouter()

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dbcon)
	if err != nil {
		return err
	}

	return nil
}