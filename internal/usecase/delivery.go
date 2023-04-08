package usecase

import (
	"context"
	"fmt"
	"github.com/dacore-x/truckly/internal/dto"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/entity"
	"sync"
	"time"
)

// DeliveryUseCase is a struct that provides all use cases of the delivery entity
type DeliveryUseCase struct {
	repo      DeliveryRepo
	geo       GeoWebAPI
	service   PriceEstimatorService
	appLogger *logger.Logger
}

// ObjectResponse is an internal struct for syncing results of goroutines
type ObjectResponse struct {
	Object string
	Error  error
}

// DistanceResponse is an internal struct for syncing results of goroutines
type DistanceResponse struct {
	Distance float64
	Error    error
}

func NewDeliveryUseCase(r DeliveryRepo, g GeoWebAPI, s PriceEstimatorService, l *logger.Logger) *DeliveryUseCase {
	return &DeliveryUseCase{repo: r, geo: g, service: s, appLogger: l}
}

// CreateDelivery creates new user's delivery
func (uc *DeliveryUseCase) CreateDelivery(ctx context.Context, delivery *entity.Delivery) error {
	fromObj := make(chan ObjectResponse, 2)
	toObj := make(chan ObjectResponse, 2)
	distCh := make(chan DistanceResponse, 2)

	defer close(fromObj)
	defer close(toObj)
	defer close(distCh)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		fromObject, err := uc.geo.GetObjectByCoords(delivery.Geo.FromLatitude, delivery.Geo.FromLongitude)
		fromObj <- ObjectResponse{Object: fromObject, Error: err}
		wg.Done()
	}()
	go func() {
		toObject, err := uc.geo.GetObjectByCoords(delivery.Geo.ToLatitude, delivery.Geo.ToLongitude)
		toObj <- ObjectResponse{Object: toObject, Error: err}
		wg.Done()
	}()
	go func() {
		distance, err := uc.geo.GetDistanceBetweenPoints(delivery.Geo.FromLatitude, delivery.Geo.FromLongitude, delivery.Geo.ToLatitude, delivery.Geo.ToLongitude)
		distCh <- DistanceResponse{Distance: distance, Error: err}
		wg.Done()
	}()
	wg.Wait()

	fromObjResponse, _ := <-fromObj
	toObjResponse, _ := <-toObj
	distResponse, _ := <-distCh

	if fromObjResponse.Error != nil {
		err := fmt.Errorf("error getting from geo object")
		uc.appLogger.Error(err)
		return err
	}

	if toObjResponse.Error != nil {
		err := fmt.Errorf("error getting to geo object")
		uc.appLogger.Error(err)
		return err
	}

	if distResponse.Error != nil {
		err := fmt.Errorf("error finding distance between points")
		uc.appLogger.Error(err)
		return err
	}

	body := &dto.EstimatePriceInternalRequestBody{
		TypeID:    delivery.TypeID,
		HasLoader: delivery.HasLoader,
		Time:      time.Now(),
		Distance:  distResponse.Distance / 1000, // in km
	}
	price, err := uc.service.EstimateDeliveryPrice(body)
	if err != nil {
		err := fmt.Errorf("error finding object from")
		uc.appLogger.Error(err)
		return err
	}

	delivery.Geo.FromObject = fromObjResponse.Object
	delivery.Geo.ToObject = toObjResponse.Object
	delivery.Geo.Distance = distResponse.Distance
	delivery.Price = price

	err = uc.repo.CreateDelivery(ctx, delivery)
	if err != nil {
		uc.appLogger.Error(err)
		return err
	}
	return nil
}
