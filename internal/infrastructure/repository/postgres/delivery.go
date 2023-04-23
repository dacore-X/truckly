package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
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

func (dr *DeliveryRepo) GetDeliveryByID(ctx context.Context, clientID, deliveryID int) (*dto.DeliveryFullInfoResponse, error) {
	query := `
		SELECT deliveries.id, type_id, courier_id, status_id, price, has_loader,
       	geo.from_latitude, geo.from_longitude, geo.from_object, geo.to_latitude, geo.to_longitude, geo.to_object,
       	geo.distance, deliveries.created_at
		FROM deliveries
		INNER JOIN geo ON deliveries.geo_id = geo.id
		WHERE (client_id = $1 OR courier_id = $1 OR $1 IN (
		    SELECT users.id
		    FROM users INNER JOIN meta ON users.id = meta.user_id
		    WHERE is_admin = true
		)) AND deliveries.id = $2`

	queryCourier := `
		SELECT name, phone_number, rating, created_at
		FROM users INNER JOIN meta
		ON users.id = meta.user_id
		WHERE users.id = $1`

	response := &dto.DeliveryFullInfoResponse{}
	row := dr.QueryRowContext(ctx, query, clientID, deliveryID)
	var courierID sql.NullInt64
	err := row.Scan(&response.ID, &response.TypeID, &courierID, &response.StatusID, &response.Price, &response.HasLoader,
		&response.FromObject.Latitude, &response.FromObject.Longitude, &response.FromObject.Object, &response.ToObject.Latitude,
		&response.ToObject.Longitude, &response.ToObject.Object, &response.Distance, &response.Time)

	if err == sql.ErrNoRows {
		err = fmt.Errorf("user with this id doesn't have permission to get delivery")
		dr.appLogger.Error(err)
		return nil, err
	}

	if err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}

	if courierID.Valid {
		row = dr.QueryRowContext(ctx, queryCourier, courierID.Int64)
		err = row.Scan(&response.Courier.Name, &response.Courier.PhoneNumber, &response.Courier.Rating, &response.Courier.CreatedAt)
		if err != nil {
			dr.appLogger.Error(err)
			return nil, err
		}
	}

	return response, nil
}

func (dr *DeliveryRepo) GetActiveDeliveryAmount(ctx context.Context, courierID int) (int, error) {
	query := `SELECT COUNT(id) FROM deliveries WHERE status_id = 2 AND courier_id = $1`
	var amount int
	row := dr.QueryRowContext(ctx, query, courierID)
	err := row.Scan(&amount)
	if err != nil {
		dr.appLogger.Error(err)
		return 0, err
	}
	return amount, nil
}

func (dr *DeliveryRepo) AcceptDelivery(ctx context.Context, courierID, deliveryID int) error {
	query := `UPDATE deliveries SET courier_id = $1, status_id = 2 WHERE id = $2`
	result, err := dr.ExecContext(ctx, query, courierID, deliveryID)
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}

	if rows != 1 {
		err = fmt.Errorf("error accepting delivery")
		dr.appLogger.Error(err)
		return err
	}
	return nil
}

func (dr *DeliveryRepo) IsDeliveryPerformer(ctx context.Context, courierID, deliveryID int) (bool, error) {
	query := `SELECT COUNT(id) FROM deliveries WHERE courier_id = $1 AND id = $2`
	var amount int
	row := dr.QueryRowContext(ctx, query, courierID, deliveryID)
	err := row.Scan(&amount)
	if err != nil {
		dr.appLogger.Error(err)
		return false, err
	}

	if amount != 1 {
		return false, nil
	}
	return true, nil
}

func (dr *DeliveryRepo) IsDeliveryOwner(ctx context.Context, clientID, deliveryID int) (bool, error) {
	query := `SELECT COUNT(id) FROM deliveries WHERE client_id = $1 AND id = $2`
	var amount int
	row := dr.QueryRowContext(ctx, query, clientID, deliveryID)
	err := row.Scan(&amount)
	if err != nil {
		dr.appLogger.Error(err)
		return false, err
	}

	if amount != 1 {
		return false, nil
	}
	return true, nil
}

func (dr *DeliveryRepo) ChangeDeliveryStatus(ctx context.Context, statusID, deliveryID int) error {
	query := `UPDATE deliveries SET status_id = $1 WHERE id = $2`
	result, err := dr.ExecContext(ctx, query, statusID, deliveryID)
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		dr.appLogger.Error(err)
		return err
	}

	if rows != 1 {
		err = fmt.Errorf("error changing status")
		dr.appLogger.Error(err)
		return err
	}
	return nil
}

func (dr *DeliveryRepo) GetDeliveriesByGeolocation(ctx context.Context, q *dto.DeliveryListGeolocationQuery, searchD float64) ([]*dto.DeliveryBriefResponse, error) {
	query := `
	SELECT deliveries.id, type_id, has_loader, status_id, price, geo.from_object, geo.to_object, geo.distance, created_at
	FROM deliveries INNER JOIN geo ON deliveries.geo_id = geo.id
	WHERE (geo.from_latitude BETWEEN $1 AND $2) AND (geo.from_longitude BETWEEN $3 AND $4) AND status_id = 1
	LIMIT 10 OFFSET $5
	`

	latSearchFrom := q.Latitude - searchD
	latSearchTo := q.Latitude + searchD

	lonSearchFrom := q.Longitude - searchD
	lonSearchTo := q.Longitude + searchD

	rows, err := dr.QueryContext(ctx, query, latSearchFrom, latSearchTo, lonSearchFrom, lonSearchTo, (q.Page-1)*10)
	if err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}
	defer rows.Close()

	results := make([]*dto.DeliveryBriefResponse, 0)
	for rows.Next() {
		result := &dto.DeliveryBriefResponse{}
		err = rows.Scan(&result.ID, &result.TypeID, &result.HasLoader, &result.StatusID, &result.Price, &result.FromObject, &result.ToObject, &result.Distance, &result.Time)
		if err != nil {
			dr.appLogger.Error(err)
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}

	if len(results) == 0 {
		err = fmt.Errorf("results not found")
		dr.appLogger.Error(err)
		return nil, err
	}

	return results, nil
}

func (dr *DeliveryRepo) GetDeliveriesByClientID(ctx context.Context, clientID int, page int) ([]*dto.DeliveryBriefResponse, error) {
	query := `
	SELECT deliveries.id, type_id, has_loader, status_id, price, geo.from_object, geo.to_object, geo.distance, created_at
	FROM deliveries INNER JOIN geo ON deliveries.geo_id = geo.id
	WHERE client_id = $1
	ORDER BY deliveries.id DESC
	LIMIT 10 OFFSET $2
	`

	rows, err := dr.QueryContext(ctx, query, clientID, (page-1)*10)
	if err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}
	defer rows.Close()

	results := make([]*dto.DeliveryBriefResponse, 0)
	for rows.Next() {
		result := &dto.DeliveryBriefResponse{}
		err = rows.Scan(&result.ID, &result.TypeID, &result.HasLoader, &result.StatusID, &result.Price, &result.FromObject, &result.ToObject, &result.Distance, &result.Time)
		if err != nil {
			dr.appLogger.Error(err)
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}

	if len(results) == 0 {
		err = fmt.Errorf("results not found")
		dr.appLogger.Error(err)
		return nil, err
	}

	return results, nil
}

func (dr *DeliveryRepo) GetDeliveriesByCourierID(ctx context.Context, courierID int, page int) ([]*dto.DeliveryBriefResponse, error) {
	query := `
	SELECT deliveries.id, type_id, has_loader, status_id, price, geo.from_object, geo.to_object, geo.distance, created_at
	FROM deliveries INNER JOIN geo ON deliveries.geo_id = geo.id
	WHERE courier_id = $1
	ORDER BY deliveries.id DESC
	LIMIT 10 OFFSET $2
	`

	rows, err := dr.QueryContext(ctx, query, courierID, (page-1)*10)
	if err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}
	defer rows.Close()

	results := make([]*dto.DeliveryBriefResponse, 0)
	for rows.Next() {
		result := &dto.DeliveryBriefResponse{}
		err = rows.Scan(&result.ID, &result.TypeID, &result.HasLoader, &result.StatusID, &result.Price, &result.FromObject, &result.ToObject, &result.Distance, &result.Time)
		if err != nil {
			dr.appLogger.Error(err)
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		dr.appLogger.Error(err)
		return nil, err
	}

	if len(results) == 0 {
		err = fmt.Errorf("results not found")
		dr.appLogger.Error(err)
		return nil, err
	}

	return results, nil
}
