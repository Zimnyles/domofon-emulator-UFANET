package mqttserver

import (
	"context"
	"domofonEmulator/server/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx"
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

func (r *MqttServerRepository) UpdateIntercomStatus(id int, status bool) error {
	query := `
        UPDATE intercoms 
        SET 
            intercom_status = $1,
            door_status = CASE WHEN $1 = false THEN true ELSE false END,
			is_calling = CASE WHEN $1 = true THEN false ELSE false END,
            updated_at = NOW()
        WHERE id = $2
    `
	result, err := r.Dbpool.Exec(context.Background(), query, status, id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no intercom found with id %d", id)
	}
	return nil
}

func (r *MqttServerRepository) UpdateIntercomDoorStatus(id int, status bool, apartment int) error {
	query := `
        UPDATE intercoms 
        SET 
            door_status = $1,
			openeddoorapartment = $3,
            updated_at = NOW()
        WHERE id = $2
    `
	result, err := r.Dbpool.Exec(context.Background(), query, status, id, apartment)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no intercom found with id %d", id)
	}
	return nil
}

func (r *MqttServerRepository) UpdateIntercomActiveStatus(id int, isActive bool) error {
	query := `
        UPDATE intercoms 
        SET 
            is_active = $1,
            updated_at = NOW()
        WHERE id = $2
    `
	result, err := r.Dbpool.Exec(context.Background(), query, isActive, id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no intercom found with id %d", id)
	}

	return nil
}

func (r *MqttServerRepository) UpdateCallStatus(id int, status bool, apartment int) error {
	query := `
        UPDATE intercoms 
        SET 
            is_calling = $1,
			calledapartment = $3,
            updated_at = NOW()
        WHERE id = $2
    `
	result, err := r.Dbpool.Exec(context.Background(), query, status, id, apartment)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no intercom found with id %d", id)
	}
	return nil
}

func (r *MqttServerRepository) GetIntercomByID(id int, logger *zerolog.Logger) (models.Intercom, error) {
	query := `
        SELECT 
            id,
            mac_address,
            intercom_status,
            door_status,
            address,
            number_of_apartments,
            is_calling,
            created_at,
            updated_at,
			calledapartment,
			openeddoorapartment,
			is_active
        FROM intercoms
        WHERE id = $1
    `
	var intercom models.Intercom
	err := r.Dbpool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&intercom.ID,
		&intercom.MAC,
		&intercom.IntercomStatus,
		&intercom.DoorStatus,
		&intercom.Address,
		&intercom.NumberOfApartments,
		&intercom.IsCalling,
		&intercom.CreatedAt,
		&intercom.UpdatedAt,
		&intercom.CalledApartment,
		&intercom.OpenedDoorApartment,
		&intercom.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Intercom{}, nil
		}
		if !strings.Contains(err.Error(), "no rows") {
			logger.Error().
				Err(err).
				Int("id", id).
				Msg("Database error while getting intercom")
		}
		return models.Intercom{}, fmt.Errorf("database operation failed")
	}

	return intercom, nil
}
