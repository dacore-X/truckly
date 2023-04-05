package usecase

import (
	"context"
	"fmt"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/entity"
)

// DeliveryUseCase is a struct that provides all use cases of the delivery entity
type DeliveryUseCase struct {
	repo      DeliveryRepo
	geo       GeoWebAPI
	appLogger *logger.Logger
}

func NewDeliveryUseCase(r DeliveryRepo, g GeoWebAPI, l *logger.Logger) *DeliveryUseCase {
	return &DeliveryUseCase{
		repo:      r,
		geo:       g,
		appLogger: l,
	}
}

func (uc *DeliveryUseCase) CreateDelivery(ctx context.Context, delivery *entity.Delivery) error {
	objectFrom, err := uc.geo.GetObjectByCoords(delivery.Geo.FromLatitude, delivery.Geo.FromLongitude)
	if err != nil {
		err := fmt.Errorf("error finding object from")
		uc.appLogger.Error(err)
		return err
	}

	objectTo, err := uc.geo.GetObjectByCoords(delivery.Geo.ToLatitude, delivery.Geo.ToLongitude)
	if err != nil {
		err := fmt.Errorf("error finding object from")
		uc.appLogger.Error(err)
		return err
	}

	uc.appLogger.Info(objectFrom)
	uc.appLogger.Info(objectTo)

	err = uc.repo.CreateDelivery(ctx, delivery)
	if err != nil {
		uc.appLogger.Error(err)
		return err
	}
	return nil
}
