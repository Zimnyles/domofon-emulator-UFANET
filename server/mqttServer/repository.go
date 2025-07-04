package mqttserver

import (
	"context"
	"domofonEmulator/server/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MqttServerRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewMqttRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *MqttServerRepository {
	return &MqttServerRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *MqttServerRepository) CreateNewIntercom(intercom models.CreateIntercomCredentials, logger *zerolog.Logger) (int, bool, error) {
	query := `
        INSERT INTO intercoms (mac_address, address, number_of_apartments, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (mac_address) DO UPDATE SET
            updated_at = $5
        RETURNING id, (xmax = 0) AS is_new
    `
	var (
		id    int
		isNew bool
	)
	err := r.Dbpool.QueryRow(
		context.Background(),
		query,
		intercom.MAC,
		intercom.Address,
		intercom.NumberOfApartments,
		time.Now(),
		time.Now(),
	).Scan(&id, &isNew)
	if err != nil {
		logger.Error().
			Err(err).
			Str("mac", intercom.MAC).
			Msg("Failed to upsert intercom")
		return 0, false, fmt.Errorf("failed to upsert intercom: %w", err)
	}
	return id, isNew, nil
}
