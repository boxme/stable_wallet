package main

import (
	"log"
	"net/http"
	"os"
	"stable_wallet/main/internal/app"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Mux *http.ServeMux
	App *app.App
}

func createServer(db *pgxpool.Pool) (*Server, error) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	server := &Server{
		Mux: http.NewServeMux(),
		App: &app.App{
			Db:       db,
			ErrorLog: errorLog,
			InfoLog:  infoLog,
		},
	}

	return server, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}
