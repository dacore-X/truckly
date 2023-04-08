package webapi

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/sirupsen/logrus"

	"github.com/dacore-x/truckly/internal/dto"
)

func isCloseCoordinate(test, response float64) bool {
	var tolerance = 0.002
	return math.Abs(test-response) < tolerance
}

func isCloseDistance(test, response float64) bool {
	var tolerance = 20.
	return math.Abs(test-response)/response*100 < tolerance
}

func TestGeo_GetCoordsByObject(t *testing.T) {
	APIKeyCatalog := os.Getenv("API_KEY_CATALOG")
	BaseURLCatalog := os.Getenv("BASE_URL_CATALOG")

	// Create instance of GeoWebAPI
	testLogger := logrus.New()
	g := New(
		&config.GEO{
			APIKeyCatalog:  APIKeyCatalog,
			BaseURLCatalog: BaseURLCatalog,
		},
		logger.New(testLogger),
	)
	// Query string for finding geo objects
	type args struct {
		q string
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.PointResponse
		wantErr bool
		error   string
	}{
		{
			name: "empty query string",
			args: args{
				q: "",
			},
			wantErr: true,
			error:   "query is empty",
		},
		{
			name: "results not found",
			args: args{
				q: "уkfmlwemfpe 3333",
			},
			wantErr: true,
			error:   "bad status code from geo",
		},
		{
			name: "usual case",
			args: args{
				q: "Мичуринский проспект 38 москва",
			},
			want: &dto.PointResponse{
				Lat: 55.696392,
				Lon: 37.494836,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.GetCoordsByObject(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCoordsByObject() error = %v, wantErr %v", err != nil, tt.wantErr)
				return
			}
			if err != nil {
				if err.Error() != tt.error {
					t.Errorf("GetCoordsByObject() error = %v, expected error = %v", err, tt.error)
					return
				}
			}
			// check for the similarity of test point and response point
			if err == nil && !(isCloseCoordinate(tt.want.Lat, got.Lat) && isCloseCoordinate(tt.want.Lon, got.Lon)) {
				t.Errorf("GetCoordsByObject() result = %v, expected result = %v", got, tt.want)
				return
			}
		})
	}
}

func TestGeo_GetObjectByCoords(t *testing.T) {
	APIKeyCatalog := os.Getenv("API_KEY_CATALOG")
	BaseURLCatalog := os.Getenv("BASE_URL_CATALOG")

	// Create instance of GeoWebAPI
	testLogger := logrus.New()
	g := New(
		&config.GEO{
			APIKeyCatalog:  APIKeyCatalog,
			BaseURLCatalog: BaseURLCatalog,
		},
		logger.New(testLogger),
	)

	type args struct {
		lat float64
		lon float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		error   string
	}{
		{
			name: "null coordinate error",
			args: args{
				lat: 56.34555,
				lon: 0,
			},
			want:    "",
			wantErr: true,
			error:   "coordinate couldn't be zero",
		},
		{
			name: "рту мирэа usual case",
			args: args{
				lat: 55.66214,
				lon: 37.47803,
			},
			want: "проспект Вернадского, 86 ст8",
		},
		{
			name: "bad request from 2gis",
			args: args{
				lat: 1234,
				lon: 234,
			},
			wantErr: true,
			error:   "bad status code from geo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.GetObjectByCoords(tt.args.lat, tt.args.lon)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObjectByCoords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				if err.Error() != tt.error {
					t.Errorf("GetCoordsByObject() error = %v, expected error = %v", err, tt.error)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetObjectByCoords() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeo_GetDistanceBetweenPoints(t *testing.T) {
	APIKeyRouting := os.Getenv("API_KEY_NAVIGATION")
	BaseURLRouting := os.Getenv("BASE_URL_ROUTING")

	// Create instance of GeoWebAPI
	testLogger := logrus.New()
	g := New(
		&config.GEO{
			APIKeyRouting:  APIKeyRouting,
			BaseURLRouting: BaseURLRouting,
		},
		logger.New(testLogger),
	)

	type args struct {
		latFrom float64
		lonFrom float64
		latTo   float64
		lonTo   float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
		error   string
	}{
		{
			name: "null coordinate error",
			args: args{
				latFrom: 0,
				lonFrom: 37.222,
				latTo:   56.444,
				lonTo:   37.444,
			},
			wantErr: true,
			error:   "coordinate couldn't be zero",
		},
		{
			name: "usual case",
			args: args{
				latFrom: 55.680683,
				lonFrom: 37.484534,
				latTo:   55.669856,
				lonTo:   37.481003,
			},
			want: 3500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.GetDistanceBetweenPoints(tt.args.latFrom, tt.args.lonFrom, tt.args.latTo, tt.args.lonTo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDistanceBetweenPoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				if err.Error() != tt.error {
					t.Errorf("GetCoordsByObject() error = %v, expected error = %v", err, tt.error)
					return
				}
			} else {
				if !isCloseDistance(tt.want, got) {
					t.Errorf("GetDistanceBetweenPoints() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
