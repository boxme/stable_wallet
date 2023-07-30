package server

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Mux *http.ServeMux
	Db  *pgxpool.Pool
}

func CreateNewServer(db *pgxpool.Pool) (*Server, error) {
	server := &Server{
		Mux: http.NewServeMux(),
		Db:  db,
	}

	return server, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}
