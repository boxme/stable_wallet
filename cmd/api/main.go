package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"stable_wallet/main/remote"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	boolPointer := flag.Bool("prod", false, "Provide true for this flag in production")
	flag.Parse()

	cfg := remote.LoadConfig(*boolPointer)
	dbpool, err := pgxpool.New(context.Background(), cfg.Database.GetDbConnectionString())
	if err != nil {
		return err
	}
	defer dbpool.Close()

	server, err := createServer(dbpool)
	server.startRouting()
	if err != nil {
		return err
	}

	port := ":" + strconv.Itoa((cfg.Port))
	srv := &http.Server{
		Addr:         port,
		Handler:      server.App.LogRequest(server),
		ErrorLog:     server.App.ErrorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	server.App.InfoLog.Printf("Starting server on %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		server.App.ErrorLog.Fatal(err)
	}
	return err
}
