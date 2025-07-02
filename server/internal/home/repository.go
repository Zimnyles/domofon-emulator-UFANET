package home

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type HomeRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewHomeRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *HomeRepository {
	return &HomeRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
