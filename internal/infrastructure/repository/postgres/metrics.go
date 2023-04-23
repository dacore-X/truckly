package postgres

import (
	"context"
	"database/sql"
	"math"

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
// and differences in percents between previous and current day for corresponding values
// from the database and returns it
func (mr *MetricsRepo) GetDeliveriesCntPerDay(ctx context.Context) (*dto.DeliveriesCntPerDay, error) {
	tx, err := mr.Begin()
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	defer tx.Rollback()

	resp := &dto.DeliveriesCntPerDay{}

	queryNewToday := `
		SELECT COUNT(*) as cnt
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400
	`
	// Count new deliveries per last 24 hours
	rowNewToday := mr.QueryRowContext(ctx, queryNewToday)

	err = rowNewToday.Scan(&resp.NewCnt)
	if err == sql.ErrNoRows {
		resp.NewCnt = 0
	} else if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	queryCompletedToday := `
		SELECT COUNT(*) as cnt
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400
			AND status_id=3
		GROUP BY status_id
	`

	// Count completed deliveries per last 24 hours
	rowCompletedToday := mr.QueryRowContext(ctx, queryCompletedToday)

	// Count completed deliveries per last 24 hours
	err = rowCompletedToday.Scan(&resp.CompletedCnt)
	if err == sql.ErrNoRows {
		resp.CompletedCnt = 0
	} else if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Variables to store count per previous 24 hours
	var newYesterdayCnt, completedYesterdayCnt int

	queryNewYesterday := `
		SELECT COUNT(*) as cnt
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) >= 86400
			AND EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 * 2
	`

	// Count new deliveries per previous 24 hours
	rowNewYesterday := mr.QueryRowContext(ctx, queryNewYesterday)

	err = rowNewYesterday.Scan(&newYesterdayCnt)
	if err == sql.ErrNoRows {
		newYesterdayCnt = 0
	} else if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Calculate differences between today and yesterday for new deliveries
	// Check if count is zero to prevent zero division
	if newYesterdayCnt != 0 {
		newCntDiff := float64(resp.NewCnt-newYesterdayCnt) * 100 / float64(newYesterdayCnt)
		resp.NewCntDiff = math.Round(newCntDiff*10) / 10
	} else if newYesterdayCnt+resp.NewCnt == 0 {
		resp.NewCntDiff = 0.0 // If both days' count is zero then difference is 0%
	} else if newYesterdayCnt == 0 {
		resp.NewCntDiff = 100.0 // If yesterday's count is zero then difference is 100%
	} else if resp.NewCnt == 0 {
		resp.NewCntDiff = -100.0 // If today's count is zero then difference is -100%
	}

	queryCompletedYesterday := `
		SELECT COUNT(*) as cnt
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) >= 86400
			AND EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 * 2
			AND status_id=3
		GROUP BY status_id
	`

	// Count completed deliveries per previous 24 hours
	rowCompletedYesterday := mr.QueryRowContext(ctx, queryCompletedYesterday)

	err = rowCompletedYesterday.Scan(&completedYesterdayCnt)
	if err == sql.ErrNoRows {
		completedYesterdayCnt = 0
	} else if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Calculate differences between today and yesterday for completed deliveries
	// Check if count is zero to prevent zero division
	if completedYesterdayCnt != 0 {
		completedCntDiff := float64(resp.CompletedCnt-completedYesterdayCnt) * 100 / float64(completedYesterdayCnt)
		resp.CompletedCntDiff = math.Round(completedCntDiff*10) / 10
	} else if completedYesterdayCnt+resp.CompletedCnt == 0 {
		resp.CompletedCntDiff = 0.0 // If both days' count is zero then difference is 0%
	} else if completedYesterdayCnt == 0 { //
		resp.CompletedCntDiff = 100.0 // If yesterday's count is zero then difference is 100%
	} else if resp.CompletedCnt == 0 {
		resp.CompletedCntDiff = -100.0 // If today's count is zero then difference is -100%
	}

	if err = tx.Commit(); err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetRevenuePerDay fetches revenue sum per last 24 hours and difference in percents
// between previous and current day for revenue from the database and returns it
func (mr *MetricsRepo) GetRevenuePerDay(ctx context.Context) (*dto.RevenuePerDay, error) {
	tx, err := mr.Begin()
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	defer tx.Rollback()

	resp := &dto.RevenuePerDay{}

	queryRevenueToday := `
		SELECT COALESCE(SUM(price), 0) AS sum
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400
			AND status_id = 3
	`

	// Revenue sum per last 24 hours
	rowRevenueToday := mr.QueryRowContext(ctx, queryRevenueToday)

	err = rowRevenueToday.Scan(&resp.Revenue)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Variable to store revenue sum per previous 24 hours
	var revenueYesterday int

	queryRevenueYesterday := `
		SELECT COALESCE(SUM(price), 0) AS sum
		FROM deliveries
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) >= 86400
			AND EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 * 2
			AND status_id = 3
	`

	// Revenue sum per previous 24 hours
	rowRevenueYesterday := mr.QueryRowContext(ctx, queryRevenueYesterday)

	err = rowRevenueYesterday.Scan(&revenueYesterday)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Calculate differences between today and yesterday for revenue sum
	// Check if sum is zero to prevent zero division
	if revenueYesterday != 0 {
		revenueDiff := float64(resp.Revenue-revenueYesterday) * 100 / float64(revenueYesterday)
		resp.RevenueDiff = math.Round(revenueDiff*10) / 10
	} else if revenueYesterday+resp.Revenue == 0 {
		resp.RevenueDiff = 0.0 // If both days' revenue sum is zero then difference is 0%
	} else if revenueYesterday == 0 { //
		resp.RevenueDiff = 100.0 // If yesterday's revenue sum is zero then difference is 100%
	} else if resp.Revenue == 0 {
		resp.RevenueDiff = -100.0 // If today's revenue sum is zero then difference is -100%
	}

	if err = tx.Commit(); err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetNewClientsCntPerDay fetches new registered clients' count per last 24 hours
// and difference in percents between previous and current day for new registered clients' count
// from the database and returns itfrom the database and returns it
func (mr *MetricsRepo) GetNewClientsCntPerDay(ctx context.Context) (*dto.NewClientsCntPerDay, error) {
	tx, err := mr.Begin()
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	defer tx.Rollback()

	resp := &dto.NewClientsCntPerDay{}

	queryCntToday := `
		SELECT COUNT(*)
		FROM users INNER JOIN meta ON users.id = meta.user_id 
		WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 
			AND meta.is_courier = FALSE
			AND meta.is_admin = FALSE
	`

	// New registered clients' count per last 24 hours
	rowCntToday := mr.QueryRowContext(ctx, queryCntToday)

	err = rowCntToday.Scan(&resp.NewClientsCnt)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Variable to store new registered clients' count per previous 24 hours
	var cntYesterday int

	queryCntYesterday := `
	SELECT COUNT(*)
	FROM users INNER JOIN meta ON users.id = meta.user_id 
	WHERE EXTRACT(EPOCH FROM (NOW() - created_at)) >= 86400
		AND EXTRACT(EPOCH FROM (NOW() - created_at)) < 86400 * 2
		AND meta.is_courier = FALSE
		AND meta.is_admin = FALSE
`

	// New registered clients' count per previous 24 hours
	rowCntYesterday := mr.QueryRowContext(ctx, queryCntYesterday)

	err = rowCntYesterday.Scan(&cntYesterday)
	if err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}

	// Calculate differences between today and yesterday for new registered clients' count
	// Check if count is zero to prevent zero division
	if cntYesterday != 0 {
		cntDiff := float64(resp.NewClientsCnt-cntYesterday) * 100 / float64(cntYesterday)
		resp.NewClientsCntDiff = math.Round(cntDiff*10) / 10
	} else if cntYesterday+resp.NewClientsCnt == 0 {
		resp.NewClientsCntDiff = 0.0 // If both days' new registered clients' count is zero then difference is 0%
	} else if cntYesterday == 0 { //
		resp.NewClientsCntDiff = 100.0 // If yesterday's new registered clients' count is zero then difference is 100%
	} else if resp.NewClientsCnt == 0 {
		resp.NewClientsCntDiff = -100.0 // If today's new registered clients' count is zero then difference is -100%
	}

	if err = tx.Commit(); err != nil {
		mr.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetDeliveryTypesPercentPerDay fetches different delivery types' percentages per last 24 hours
// from the database and returns it
func (mr *MetricsRepo) GetDeliveryTypesPercentPerDay(ctx context.Context) (*dto.DeliveryTypesPercentPerDay, error) {
	query := `
		SELECT type_id, ROUND(COUNT(type_id) / SUM(COUNT(type_id)) OVER() * 100, 3)
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
	// Count is 0 by default if there is no deliveries with such delivery type ID per last 24 hours
	// Each index stands for n-1 delivery type ID
	percents := make([]float64, 5)
	for rows.Next() {
		var (
			typeID  int
			percent float64
		)
		if err := rows.Scan(&typeID, &percent); err != nil {
			mr.appLogger.Error(err)
			return nil, err
		}

		// Change value if such delivery type was presented in the database per last 24 hours
		percents[typeID-1] = percent
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
		SELECT from_object, from_longitude, from_latitude, to_object, price
		FROM deliveries INNER JOIN geo ON deliveries.geo_id = geo.id
		WHERE status_id = 1
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
			FromObject: &dto.GeoObjectResponse{},
		}
		if err := rows.Scan(&info.FromObject.Object, &info.FromObject.Longitude, &info.FromObject.Latitude, &info.ToObject, &info.Price); err != nil {
			mr.appLogger.Error(err)
			return nil, err
		}
		currentDeliveries = append(currentDeliveries, info)
	}
	list := &dto.MetricsDeliveriesResponse{Deliveries: currentDeliveries}

	return list, nil
}
