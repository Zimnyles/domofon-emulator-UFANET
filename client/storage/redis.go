package storage

import (
	"domofonEmulator/client/models"
	"domofonEmulator/config"
	"encoding/json"
	"fmt"
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

func (s *SessionStorage) GetActiveIntercomData(c *fiber.Ctx) (models.Intercom, error) {
	sess, err := s.GetSession(c)
	if err != nil {
		return models.Intercom{}, err
	}

	raw := sess.Get("intercom_data")
	jsonStr, ok := raw.(string)
	if !ok {
		return models.Intercom{}, err
	}

	var intercomData models.Intercom
	if err := json.Unmarshal([]byte(jsonStr), &intercomData); err != nil {
		return models.Intercom{}, err
	}

	return intercomData, nil
}

func (s *SessionStorage) SetActiveIntercomData(c *fiber.Ctx, intercom models.Intercom) error {
	sess, err := s.GetSession(c)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	jsonData, err := json.Marshal(intercom)
	if err != nil {
		return fmt.Errorf("failed to marshal intercom data: %w", err)
	}
	sess.Set("intercom_data", string(jsonData))
	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}


