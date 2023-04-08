package usecase

import (
	"context"
	"errors"
	"github.com/dacore-x/truckly/internal/dto"
	"time"
)

// PriceEstimatorUseCase is a struct that provides all use cases for estimating delivery prices
type PriceEstimatorUseCase struct {
	service PriceEstimatorService
	geo     GeoWebAPI
}

func NewPriceEstimatorUseCase(s PriceEstimatorService, g GeoWebAPI) *PriceEstimatorUseCase {
	return &PriceEstimatorUseCase{
		service: s,
		geo:     g,
	}
}

// EstimateDeliveryPrice usecase estimates delivery price
func (uc *PriceEstimatorUseCase) EstimateDeliveryPrice(ctx context.Context, req *dto.EstimatePriceRequestBody) (float64, error) {
	if req.TypeID < 1 || req.TypeID > 5 {
		return 0, errors.New("incorrect type id")
	}

	distance, err := uc.geo.GetDistanceBetweenPoints(req.FromPoint.Lat, req.FromPoint.Lon, req.ToPoint.Lat, req.ToPoint.Lon)
	if err != nil {
		return 0, err
	}

	body := &dto.EstimatePriceInternalRequestBody{
		TypeID:    req.TypeID,
		HasLoader: req.HasLoader,
		Time:      time.Now(),
		Distance:  distance / 1000, // in km
	}
	price, err := uc.service.EstimateDeliveryPrice(body)
	if err != nil {
		return 0, err
	}

	return price, nil
}
