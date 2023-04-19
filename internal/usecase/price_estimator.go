package usecase

import (
	"context"
	"errors"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
	"time"
)

// PriceEstimatorUseCase is a struct that provides all use cases for estimating delivery prices
type PriceEstimatorUseCase struct {
	service   PriceEstimatorService
	geo       GeoWebAPI
	appLogger *logger.Logger
}

func NewPriceEstimatorUseCase(s PriceEstimatorService, g GeoWebAPI, l *logger.Logger) *PriceEstimatorUseCase {
	return &PriceEstimatorUseCase{
		geo:       g,
		service:   s,
		appLogger: l,
	}
}

// EstimateDeliveryPrice usecase estimates delivery price
func (uc *PriceEstimatorUseCase) EstimateDeliveryPrice(ctx context.Context, req *dto.EstimatePriceRequestBody) (float64, error) {
	if req.TypeID < 1 || req.TypeID > 5 {
		err := errors.New("incorrect type id")
		uc.appLogger.Error(err)
		return 0, err
	}

	distance, err := uc.geo.GetDistanceBetweenPoints(req.FromPoint.Lat, req.FromPoint.Lon, req.ToPoint.Lat, req.ToPoint.Lon)
	if err != nil {
		return 0, err
	}

	body := &dto.EstimatePriceInternalRequestBody{
		TypeID:    req.TypeID,
		HasLoader: req.HasLoader,
		Time:      time.Now(),
		Distance:  distance, // in m
	}
	price, err := uc.service.EstimateDeliveryPrice(body)
	if err != nil {
		uc.appLogger.Error(err)
		return 0, err
	}

	return price, nil
}
