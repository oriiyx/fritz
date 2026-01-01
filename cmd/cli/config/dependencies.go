package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type Dependencies struct {
	DB      *pgxpool.Pool
	Queries *db.Queries
	Logger  *zerolog.Logger
}
