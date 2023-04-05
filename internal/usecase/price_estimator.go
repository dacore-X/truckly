package usecase

import (
	"context"
	"errors"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
)

type PriceEstimatorUseCase struct {
	service   PriceEstimatorService
	appLogger *logger.Logger
}

func NewPriceEstimatorUseCase(s PriceEstimatorService, l *logger.Logger) *PriceEstimatorUseCase {
	return &PriceEstimatorUseCase{
		service:   s,
		appLogger: l,
	}
}

func (uc *PriceEstimatorUseCase) EstimateDeliveryPrice(ctx context.Context, body *dto.EstimatePriceRequestBody) (float64, error) {
	if body.TypeID < 1 || body.TypeID > 5 {
		err := errors.New("incorrect type id")
		uc.appLogger.Error(err)
		return 0, err
	}

	price, err := uc.service.EstimateDeliveryPrice(body)
	if err != nil {
		uc.appLogger.Error(err)
		return 0, err
	}

	return price, nil
}
