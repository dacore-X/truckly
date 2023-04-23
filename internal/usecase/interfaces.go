package usecase

import (
	"context"

	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/entity"
)

type (
	// User interface represents user's usecases
	User interface {
		CreateUser(context.Context, *dto.UserSignUpRequestBody) error
		BanUser(context.Context, int) error
		UnbanUser(context.Context, int) error
		GetUserByID(context.Context, int) (*dto.UserMeResponse, error)
		GetUserPrivateByID(context.Context, int) (*dto.UserInfoResponse, error)
		GetUserPrivateByEmail(context.Context, string) (*dto.UserInfoResponse, error)
		GetUserMeta(context.Context, int) (*dto.UserMetaResponse, error)
	}

	// UserRepo interface represents user's repository contract
	UserRepo interface {
		CreateUser(context.Context, *dto.UserSignUpRequestBody) error
		BanUser(context.Context, int) error
		UnbanUser(context.Context, int) error
		GetUserByID(context.Context, int) (*dto.UserMeResponse, error)
		GetUserPrivateByID(context.Context, int) (*dto.UserInfoResponse, error)
		GetUserPrivateByEmail(context.Context, string) (*dto.UserInfoResponse, error)
		GetUserMeta(context.Context, int) (*dto.UserMetaResponse, error)
	}

	// Delivery interface represents delivery usecases
	Delivery interface {
		CreateDelivery(context.Context, *entity.Delivery) error
		GetDeliveryByID(ctx context.Context, clientID int, deliveryID int) (*dto.DeliveryFullInfoResponse, error)
		GetDeliveriesByGeolocation(ctx context.Context, query *dto.DeliveryListGeolocationQuery) ([]*dto.DeliveryBriefResponse, error)
		GetDeliveriesByClientID(ctx context.Context, clientID int, page int) ([]*dto.DeliveryBriefResponse, error)
		GetDeliveriesByCourierID(ctx context.Context, courierID int, page int) ([]*dto.DeliveryBriefResponse, error)
		AcceptDelivery(ctx context.Context, courierID int, deliveryID int) error
		ChangeDeliveryStatus(ctx context.Context, courierID, deliveryID, statusID int) error
		CancelDelivery(ctx context.Context, clientID, deliveryID int) error
	}

	// DeliveryRepo interface represents delivery's repository contract
	DeliveryRepo interface {
		CreateDelivery(context.Context, *entity.Delivery) error
		GetDeliveryByID(ctx context.Context, clientID int, deliveryID int) (*dto.DeliveryFullInfoResponse, error)
		GetDeliveriesByGeolocation(context.Context, *dto.DeliveryListGeolocationQuery, float64) ([]*dto.DeliveryBriefResponse, error)
		GetDeliveriesByClientID(ctx context.Context, clientID int, page int) ([]*dto.DeliveryBriefResponse, error)
		GetDeliveriesByCourierID(ctx context.Context, courierID int, page int) ([]*dto.DeliveryBriefResponse, error)
		AcceptDelivery(ctx context.Context, courierID int, deliveryID int) error
		GetActiveDeliveryAmount(ctx context.Context, courierID int) (int, error)
		IsDeliveryPerformer(ctx context.Context, courierID, deliveryID int) (bool, error)
		IsDeliveryOwner(ctx context.Context, courierID, deliveryID int) (bool, error)
		ChangeDeliveryStatus(ctx context.Context, statusID, deliveryID int) error
	}

	// Metrics interface represents metrics usecases
	Metrics interface {
		GetMetrics(context.Context) (*dto.MetricsPerDayResponse, error)
		GetCurrentDeliveries(context.Context) (*dto.MetricsDeliveriesResponse, error)
	}

	// MetricsRepo interface represents metrics' repository contract
	MetricsRepo interface {
		GetDeliveriesCntPerDay(context.Context) (*dto.DeliveriesCntPerDay, error)
		GetRevenuePerDay(context.Context) (*dto.RevenuePerDay, error)
		GetNewClientsCntPerDay(context.Context) (*dto.NewClientsCntPerDay, error)
		GetDeliveryTypesPercentPerDay(context.Context) (*dto.DeliveryTypesPercentPerDay, error)
		GetCurrentDeliveries(context.Context) (*dto.MetricsDeliveriesResponse, error)
	}

	Geo interface {
		GetCoordsByObject(ctx context.Context, q string) (*dto.PointResponse, error)
		GetObjectByCoords(ctx context.Context, lat, lon float64) (string, error)
	}
	// GeoWebAPI interface represents Geo API contract
	GeoWebAPI interface {
		GetCoordsByObject(q string) (*dto.PointResponse, error)
		GetObjectByCoords(lat, lon float64) (string, error)
		GetDistanceBetweenPoints(latFrom, lonFrom, latTo, lonTo float64) (float64, error)
	}

	PriceEstimator interface {
		EstimateDeliveryPrice(ctx context.Context, body *dto.EstimatePriceRequestBody) (float64, error)
	}

	PriceEstimatorService interface {
		EstimateDeliveryPrice(*dto.EstimatePriceInternalRequestBody) (float64, error)
	}
)
