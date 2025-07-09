package mainhandler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MainRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewMainRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *MainRepository {
	return &MainRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
