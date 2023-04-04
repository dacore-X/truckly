package usecase

import (
	"context"
	"errors"
	"github.com/dacore-x/truckly/internal/dto"
)

type PriceEstimatorUseCase struct {
	service PriceEstimatorService
}

func NewPriceEstimatorUseCase(s PriceEstimatorService) *PriceEstimatorUseCase {
	return &PriceEstimatorUseCase{
		service: s,
	}
}

func (uc *PriceEstimatorUseCase) EstimateDeliveryPrice(ctx context.Context, body *dto.EstimatePriceRequestBody) (float64, error) {
	if body.TypeID < 1 || body.TypeID > 5 {
		return 0, errors.New("incorrect type id")
	}

	price, err := uc.service.EstimateDeliveryPrice(body)
	if err != nil {
		return 0, err
	}

	return price, nil
}
