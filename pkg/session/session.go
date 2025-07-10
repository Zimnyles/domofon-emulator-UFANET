package session

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSessionStorage(dbpool *pgxpool.Pool) *postgres.Storage {
	storage := postgres.New(postgres.Config{
		DB:         dbpool,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})
	return storage
}

func NewSession(storage *postgres.Storage) *session.Store {
	store := session.New(session.Config{
		Storage:        storage,
		CookieHTTPOnly: true,
		Expiration:     24 * time.Hour, 
		CookieSecure:   true,         
		CookieSameSite: "Lax",
		CookieName: "server_session_3031",
	})

	return store
}
