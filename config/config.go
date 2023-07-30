package config

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Db       *pgxpool.Pool
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
