package storage

import (
	"domofonEmulator/config"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v2"
)

type SessionStorage struct {
	store *session.Store
}

func NewRedisStorage(config config.RedisConfig) *SessionStorage {
	redisStore := redis.New(redis.Config{
		Host:      config.Url,
		Port:      config.Port,
		Database:  config.DB,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	return &SessionStorage{
		store: session.New(session.Config{
			Storage:    redisStore,
			Expiration: 24 * time.Hour,
			KeyLookup:  "cookie:session_id",
		}),
	}

}

func (s *SessionStorage) GetSession(c *fiber.Ctx) (*session.Session, error) {
	return s.store.Get(c)
}

func (s *SessionStorage) SaveSession(sess *session.Session) error {
	return sess.Save()
}
