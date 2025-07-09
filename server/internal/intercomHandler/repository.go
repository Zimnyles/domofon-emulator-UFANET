package intercom

import (
	"context"
	"domofonEmulator/server/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type IntercomRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewIntercomRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *IntercomRepository {
	return &IntercomRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *IntercomRepository) GetIntercomById(id int) (models.Intercom, error) {
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
			r.CustomLogger.Error().
				Err(err).
				Int("id", id).
				Msg("Database error while getting intercom")
		}
		return models.Intercom{}, fmt.Errorf("database operation failed")
	}

	return intercom, nil
}
