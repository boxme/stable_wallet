package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"stable_wallet/main/server"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
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

	cfg := LoadConfig(*boolPointer)
	dbpool, err := pgxpool.New(context.Background(), cfg.Database.GetDbConnectionString())
	if err != nil {
		return err
	}
	defer dbpool.Close()

	server, err := server.CreateServer(dbpool)
	server.StartRouting()
	if err != nil {
		return err
	}

	port := ":" + strconv.Itoa((cfg.Port))
	srv := &http.Server{
		Addr:     port,
		Handler:  server,
		ErrorLog: server.App.ErrorLog,
	}
	server.App.InfoLog.Printf("Starting server on %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		server.App.ErrorLog.Fatal(err)
	}
	return err
}
