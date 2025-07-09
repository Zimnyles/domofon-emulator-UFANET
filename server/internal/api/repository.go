package api

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

type ApiRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewApiRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *ApiRepository {
	return &ApiRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *ApiRepository) SetIntercomDoorClose(id int) error {
	query := `
        UPDATE intercoms 
        SET 
            door_status = false,
			openeddoorapartment = 0,
            updated_at = NOW()
        WHERE id = $1
    `
	result, err := r.Dbpool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no intercom found with id %d", id)
	}
	return nil
}

func (r *ApiRepository) SetIntercomDoorOpened(id int) error {
	query := `
        UPDATE intercoms 
        SET 
            door_status = true,
			openeddoorapartment = 0,
            updated_at = NOW()
        WHERE id = $1
    `
	result, err := r.Dbpool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no intercom found with id %d", id)
	}
	return nil
}

func (r *ApiRepository) CreateNewIntercom(mac, address string, apartments int) (int, bool, error) {
	intercom := models.CreateIntercomCredentials{
		MAC:                mac,
		Address:            address,
		NumberOfApartments: apartments,
	}
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
		r.CustomLogger.Error().
			Err(err).
			Str("mac", intercom.MAC).
			Msg("Failed to upsert intercom")
		return 0, false, fmt.Errorf("failed to upsert intercom: %w", err)
	}
	return id, isNew, nil

}

func (r *ApiRepository) GetIntercomById(id int) (models.Intercom, error) {
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
			openeddoorapartment

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
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Intercom{}, nil
		}
		if !strings.Contains(err.Error(), "no rows") {
			r.CustomLogger.Error().
				Err(err).
				Int("id", id).
				Msg("Database error while getting intercom")
		}
		return models.Intercom{}, fmt.Errorf("database operation failed")
	}

	return intercom, nil
}
