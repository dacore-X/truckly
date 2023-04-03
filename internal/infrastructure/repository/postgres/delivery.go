package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dacore-x/truckly/internal/entity"
)

// DeliveryRepo is a struct that provides
// all functions to execute SQL queries
// related to `Delivery` requests
type DeliveryRepo struct {
	*sql.DB
}

func NewDeliveryRepo(db *sql.DB) *DeliveryRepo {
	return &DeliveryRepo{db}
}

func (dr *DeliveryRepo) CreateDelivery(ctx context.Context, delivery *entity.Delivery) error {
	q := `
	INSERT INTO deliveries(client_id, status_id, type_id, from_longitude, to_longitude, from_latitude, to_latitude, distance, price, has_loader)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	res, err := dr.ExecContext(ctx, q, delivery.ClientID, 1, delivery.TypeID, delivery.FromLongitude, delivery.ToLongitude, delivery.FromLatitude, delivery.ToLatitude, delivery.Distance, delivery.Price, delivery.HasLoader)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}
