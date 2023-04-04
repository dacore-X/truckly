package usecase

import (
	"context"
	"github.com/dacore-x/truckly/internal/dto"
)

type GeoUseCase struct {
	webapi GeoWebAPI
}

func NewGeoUseCase(w GeoWebAPI) *GeoUseCase {
	return &GeoUseCase{
		webapi: w,
	}
}

func (uc *GeoUseCase) GetCoordsByObject(ctx context.Context, q string) (*dto.PointResponse, error) {
	res, err := uc.webapi.GetCoordsByObject(q)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (uc *GeoUseCase) GetObjectByCoords(ctx context.Context, lat, lon float64) (string, error) {
	res, err := uc.webapi.GetObjectByCoords(lat, lon)
	if err != nil {
		return "", err
	}

	return res, nil
}
