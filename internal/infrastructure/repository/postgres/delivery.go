package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/entity"
)

// DeliveryRepo is a struct that provides
// all functions to execute SQL queries
// related to `Delivery` requests
type DeliveryRepo struct {
	*sql.DB
	appLogger *logger.Logger
}

func NewDeliveryRepo(db *sql.DB, l *logger.Logger) *DeliveryRepo {
	return &DeliveryRepo{db, l}
}

func (dr *DeliveryRepo) CreateDelivery(ctx context.Context, delivery *entity.Delivery) error {
	tx, err := dr.Begin()
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}
	defer tx.Rollback()
	//from_longitude, to_longitude, from_latitude, to_latitude, distance
	q1 := `
		INSERT INTO geo(from_longitude, from_latitude, from_object, to_longitude, to_latitude, to_object, distance)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	lastInsertID := 0
	err = tx.QueryRowContext(ctx, q1, delivery.Geo.FromLongitude, delivery.Geo.FromLatitude, delivery.Geo.FromObject, delivery.Geo.ToLongitude, delivery.Geo.ToLatitude, delivery.Geo.ToObject, delivery.Geo.Distance).Scan(&lastInsertID)
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}

	q2 := `
	INSERT INTO deliveries(client_id, status_id, type_id, geo_id, price, has_loader)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	res, err := tx.ExecContext(ctx, q2, delivery.ClientID, 1, delivery.TypeID, lastInsertID, delivery.Price, delivery.HasLoader)
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}

	if rows != 1 {
		err := fmt.Errorf("expected to affect 1 row, affected %d", rows)
		dr.appLogger.Error(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		dr.appLogger.Error(err)
		return err
	}
	return nil
}
