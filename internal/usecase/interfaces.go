package usecase

import (
	"context"

	"github.com/dacore-x/truckly/internal/dto"
)

type (
	// User interface represents user's usecases
	User interface {
		Create(context.Context, dto.UserSignUpRequestBody) error
		GetMe(context.Context, int64) (*dto.UserMeResponse, error)
		GetByID(context.Context, int64) (*dto.UserInfoResponse, error)
		GetByEmail(context.Context, string) (*dto.UserInfoResponse, error)
	}

	// UserRepo interface represents user's repository contract
	UserRepo interface {
		Create(context.Context, dto.UserSignUpRequestBody) error
		GetMe(context.Context, int64) (*dto.UserMeResponse, error)
		GetByID(context.Context, int64) (*dto.UserInfoResponse, error)
		GetByEmail(context.Context, string) (*dto.UserInfoResponse, error)
	}

	// GeoWebAPI interface represents Geo API contract
	GeoWebAPI interface {
		GetCoordsByObject(q string) (*dto.PointResponse, error)
		GetObjectByCoords(lat, lon float64) (string, error)
		GetDistanceBetweenPoints(latFrom, lonFrom, latTo, lonTo float64) (float64, error)
	}
)
