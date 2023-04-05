package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/dacore-x/truckly/internal/entity"
)

// DeliveryUseCase is a struct that provides all use cases of the delivery entity
type DeliveryUseCase struct {
	repo DeliveryRepo
	geo  GeoWebAPI
}

func NewDeliveryUseCase(r DeliveryRepo, g GeoWebAPI) *DeliveryUseCase {
	return &DeliveryUseCase{repo: r, geo: g}
}

func (uc *DeliveryUseCase) CreateDelivery(ctx context.Context, delivery *entity.Delivery) error {
	objectFrom, err := uc.geo.GetObjectByCoords(delivery.Geo.FromLatitude, delivery.Geo.FromLongitude)
	if err != nil {
		return fmt.Errorf("error finding object from")
	}

	objectTo, err := uc.geo.GetObjectByCoords(delivery.Geo.ToLatitude, delivery.Geo.ToLongitude)
	if err != nil {
		return fmt.Errorf("error finding object from")
	}

	log.Println(objectFrom)
	log.Println(objectTo)

	err = uc.repo.CreateDelivery(ctx, delivery)
	if err != nil {
		return err
	}
	return nil
}
