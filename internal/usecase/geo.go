package usecase

import (
	"context"
	"github.com/dacore-x/truckly/internal/dto"
)

// GeoUseCase is a struct that provides all use cases connected with geo data
type GeoUseCase struct {
	webapi GeoWebAPI
}

func NewGeoUseCase(w GeoWebAPI) *GeoUseCase {
	return &GeoUseCase{
		webapi: w,
	}
}

// GetCoordsByObject returning coordinates of geo object by query string
func (uc *GeoUseCase) GetCoordsByObject(ctx context.Context, q string) (*dto.PointResponse, error) {
	res, err := uc.webapi.GetCoordsByObject(q)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetObjectByCoords returning geo object string by query coordinates
func (uc *GeoUseCase) GetObjectByCoords(ctx context.Context, lat, lon float64) (string, error) {
	res, err := uc.webapi.GetObjectByCoords(lat, lon)
	if err != nil {
		return "", err
	}

	return res, nil
}
