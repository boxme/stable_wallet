package main

import (
	"log"
	"net/http"
	"os"
	"stable_wallet/main/internal/app"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Server struct {
	Mux *http.ServeMux
	App *app.App
}

func createServer(db *pgxpool.Pool) (*Server, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	server := &Server{
		Mux: http.NewServeMux(),
		App: &app.App{
			JwtSecretKey: []byte(jwtSecretKey),
			Db:           db,
			ErrorLog:     errorLog,
			InfoLog:      infoLog,
		},
	}

	return server, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}
