package postgres

import (
	"context"
	"database/sql"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
)

type MetricsRepo struct {
	*sql.DB
	appLogger *logger.Logger
}

func NewMetricsRepo(db *sql.DB, l *logger.Logger) *MetricsRepo {
	return &MetricsRepo{db, l}
}

// GetDeliveriesCntPerDay fetches new and completed deliveries' counts per last 24 hours
// from the database and returns it
func (mr *MetricsRepo) GetDeliveriesCntPerDay(ctx context.Context) (*dto.DeliveriesCntPerDay, error) {
	query := `
		SELECT COUNT(*) as cnt
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 AND status_id IN(1, 3)
		group by status_id
		order by status_id
	`
	rows, err := mr.QueryContext(ctx, query)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	defer rows.Close()

	// Slice with calculated counts from sql query
	deliveryCounts := []int{}

	for rows.Next() {
		var cnt int
		if err := rows.Scan(&cnt); err != nil {
			mr.appLogger.Error(err)
			return nil, err
		}
		deliveryCounts = append(deliveryCounts, cnt)
	}

	// Attach calculated counts to response body's fields respectively
	resp := &dto.DeliveriesCntPerDay{
		NewCnt:       deliveryCounts[0],
		CompletedCnt: deliveryCounts[1],
	}

	return resp, nil
}

// GetRevenuePerDay fetches revenue sum per last 24 hours from the database and returns it
func (mr *MetricsRepo) GetRevenuePerDay(ctx context.Context) (*dto.RevenuePerDay, error) {
	query := `
		SELECT SUM(price)
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400
	`
	row := mr.QueryRowContext(ctx, query)

	resp := &dto.RevenuePerDay{}
	err := row.Scan(&resp.Revenue)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetNewClientsCntPerDay fetches new registered clients' count per last 24 hours
// from the database and returns it
func (mr *MetricsRepo) GetNewClientsCntPerDay(ctx context.Context) (*dto.NewClientsCntPerDay, error) {
	query := `
		SELECT COUNT(*)
		FROM users INNER JOIN meta ON users.id = meta.user_id 
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 AND 
		meta.is_courier = FALSE AND 
		meta.is_admin = FALSE
	`
	row := mr.QueryRowContext(ctx, query)

	resp := &dto.NewClientsCntPerDay{}
	err := row.Scan(&resp.NewClientsCnt)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetDeliveryTypesPercentPerDay fetches different delivery types' percentages per last 24 hours
// from the database and returns it
func (mr *MetricsRepo) GetDeliveryTypesPercentPerDay(ctx context.Context) (*dto.DeliveryTypesPercentPerDay, error) {
	query := `
		SELECT ROUND(COUNT(type_id) / SUM(COUNT(type_id)) OVER() * 100, 3)
		FROM deliveries 
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400
		GROUP BY type_id
		ORDER BY type_id 
    `
	rows, err := mr.QueryContext(ctx, query)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	defer rows.Close()

	// Slice with calculated percents from sql query
	percents := []float64{}

	for rows.Next() {
		var percent float64
		if err := rows.Scan(&percent); err != nil {
			mr.appLogger.Error(err)
			return nil, err
		}
		percents = append(percents, percent)
	}

	// Attach calculated percents to response body's fields respectively
	resp := &dto.DeliveryTypesPercentPerDay{
		FootPercent:      percents[0],
		CarPercent:       percents[1],
		MinivanPercent:   percents[2],
		TruckPercent:     percents[3],
		LongTruckPercent: percents[4],
	}

	return resp, nil
}

// GetCurrentDeliveries fetches list of brief information about current deliveries
// from the database and returns it
func (mr *MetricsRepo) GetCurrentDeliveries(context.Context) (*dto.MetricsDeliveriesResponse, error) {
	query := `
		SELECT from_longitude, from_latitude, price
		FROM deliveries INNER JOIN geo ON deliveries.geo_id = geo.id
		WHERE status_id = 2
	`
	rows, err := mr.QueryContext(context.Background(), query)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	defer rows.Close()

	currentDeliveries := []*dto.CurrentDelivery{}

	for rows.Next() {
		info := &dto.CurrentDelivery{
			FromPoint: &dto.PointResponse{},
		}
		if err := rows.Scan(&info.FromPoint.Lat, &info.FromPoint.Lon, &info.Price); err != nil {
			mr.appLogger.Error(err)
			return nil, err
		}
		currentDeliveries = append(currentDeliveries, info)
	}
	list := &dto.MetricsDeliveriesResponse{Deliveries: currentDeliveries}

	return list, nil
}
