package usecase

import (
	"context"
	"github.com/dacore-x/truckly/internal/entity"

	"github.com/dacore-x/truckly/internal/dto"
)

type (
	// User interface represents user's usecases
	User interface {
		CreateUserTx(context.Context, *dto.UserSignUpRequestBody) error
		BanUser(context.Context, int) error
		UnbanUser(context.Context, int) error
		GetUserByID(context.Context, int) (*dto.UserMeResponse, error)
		GetUserPrivateByID(context.Context, int) (*dto.UserInfoResponse, error)
		GetUserPrivateByEmail(context.Context, string) (*dto.UserInfoResponse, error)
		GetUserMeta(context.Context, int) (*dto.UserMetaResponse, error)
	}

	// UserRepo interface represents user's repository contract
	UserRepo interface {
		CreateUserTx(context.Context, *dto.UserSignUpRequestBody) error
		BanUser(context.Context, int) error
		UnbanUser(context.Context, int) error
		GetUserByID(context.Context, int) (*dto.UserMeResponse, error)
		GetUserPrivateByID(context.Context, int) (*dto.UserInfoResponse, error)
		GetUserPrivateByEmail(context.Context, string) (*dto.UserInfoResponse, error)
		GetUserMeta(context.Context, int) (*dto.UserMetaResponse, error)
	}

	Delivery interface {
		CreateDelivery(context.Context, *entity.Delivery) error
	}

	// DeliveryRepo interface represents delivery's repository contract
	DeliveryRepo interface {
		CreateDelivery(context.Context, *entity.Delivery) error
	}

	// GeoWebAPI interface represents Geo API contract
	GeoWebAPI interface {
		GetCoordsByObject(q string) (*dto.PointResponse, error)
		GetObjectByCoords(lat, lon float64) (string, error)
		GetDistanceBetweenPoints(latFrom, lonFrom, latTo, lonTo float64) (float64, error)
	}
)
