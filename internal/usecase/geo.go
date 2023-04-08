package usecase

import (
	"context"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
)

// GeoUseCase is a struct that provides all use cases connected with geo data
type GeoUseCase struct {
	webapi    GeoWebAPI
	appLogger *logger.Logger
}

func NewGeoUseCase(w GeoWebAPI, l *logger.Logger) *GeoUseCase {
	return &GeoUseCase{
		webapi:    w,
		appLogger: l,
	}
}

// GetCoordsByObject returning coordinates of geo object by query string
func (uc *GeoUseCase) GetCoordsByObject(ctx context.Context, q string) (*dto.PointResponse, error) {
	res, err := uc.webapi.GetCoordsByObject(q)
	if err != nil {
		uc.appLogger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetObjectByCoords returning geo object string by query coordinates
func (uc *GeoUseCase) GetObjectByCoords(ctx context.Context, lat, lon float64) (string, error) {
	res, err := uc.webapi.GetObjectByCoords(lat, lon)
	if err != nil {
		uc.appLogger.Error(err)
		return "", err
	}

	return res, nil
}
