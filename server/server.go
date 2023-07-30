package server

import (
	"log"
	"net/http"
	"os"
	"stable_wallet/main/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Mux *http.ServeMux
	App *config.App
}

func CreateServer(db *pgxpool.Pool) (*Server, error) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	server := &Server{
		Mux: http.NewServeMux(),
		App: &config.App{
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
