package usecase

import (
	"context"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
)

// MetricsUseCase is a struct that provides
// all metrics' usecases
type MetricsUseCase struct {
	repo      MetricsRepo
	appLogger *logger.Logger
}

func NewMetricsUseCase(r MetricsRepo, l *logger.Logger) *MetricsUseCase {
	return &MetricsUseCase{
		repo:      r,
		appLogger: l,
	}
}

// GetMetrics usecase gets all metrics per last 24 hours from storage
func (uc *MetricsUseCase) GetMetrics(ctx context.Context) (*dto.MetricsPerDayResponse, error) {
	resp := &dto.MetricsPerDayResponse{}

	// Get new and completed deliveries' counts per last 24 hours
	firstMetric, err := uc.repo.GetDeliveriesCntPerDay(context.Background())
	if err != nil {
		uc.appLogger.Error(err)
		return nil, err
	}

	// Attach new and completed deliveries' count per last 24 hours metric to response
	resp.DeliveriesCnt = firstMetric

	// Get revenue sum per last 24 hours
	secondMetric, err := uc.repo.GetRevenuePerDay(context.Background())
	if err != nil {
		uc.appLogger.Error(err)
		return nil, err
	}

	// Attach revenue sum per last 24 hours metric to response
	resp.Revenue = secondMetric

	// Get new registered clients' count per last 24 hours
	thirdMetric, err := uc.repo.GetNewClientsCntPerDay(context.Background())
	if err != nil {
		uc.appLogger.Error(err)
		return nil, err
	}

	// Attach new registered clients' count per last 24 hours metric to response
	resp.NewClientsCnt = thirdMetric

	// Get different delivery types' percentages per last 24 hours
	fourthMetric, err := uc.repo.GetDeliveryTypesPercentPerDay(context.Background())
	if err != nil {
		uc.appLogger.Error(err)
		return nil, err
	}

	// Attach different delivery types' percentages per last 24 hours metric to response
	resp.DeliveryTypesPercent = fourthMetric

	return resp, nil
}

// GetCurrentDeliveries usecase gets list of brief information about current deliveries
func (uc *MetricsUseCase) GetCurrentDeliveries(context.Context) (*dto.MetricsDeliveriesResponse, error) {
	list, err := uc.repo.GetCurrentDeliveries(context.Background())
	if err != nil {
		return nil, err
	}
	return list, err
}
