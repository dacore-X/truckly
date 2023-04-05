package config

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// PG is a struct for storing Postgres connection settings
type PG struct {
	PostgresUser     string
	PostgresPassword string
	PostgresName     string
	PostgresPort     string
}

// GEO is a struct for storing API Keys and Base URLS for 2GIS
type GEO struct {
	// API Keys for 2GIS
	APIKeyCatalog string
	APIKeyRouting string

	// Base URLS for 2GIS
	BaseURLCatalog string
	BaseURLRouting string
}

// LOG is a struct for storing Logrus configatrion settings
type LOG struct {
	LogrusFormatter *logrus.TextFormatter
}

type SERVICES struct {
	Ports map[string]int
}

// Config is a struct for storing all required configuration parameters
type Config struct {
	*PG
	*GEO
	*SERVICES
	*LOG
}

// New returns application config
func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, errors.New("POSTGRES_USER is not set")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}

	name, ok := os.LookupEnv("POSTGRES_NAME")
	if !ok {
		return nil, errors.New("POSTGRES_NAME is not set")
	}

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, errors.New("POSTGRES_PORT is not set")
	}

	apiKeyCatalog, ok := os.LookupEnv("API_KEY_CATALOG")
	if !ok {
		return nil, errors.New("API_KEY_CATALOG is not set")
	}

	apiKeyRouting, ok := os.LookupEnv("API_KEY_NAVIGATION")
	if !ok {
		return nil, errors.New("API_KEY_NAVIGATION is not set")
	}

	baseURLCatalog, ok := os.LookupEnv("BASE_URL_CATALOG")
	if !ok {
		return nil, errors.New("BASE_URL_CATALOG is not set")
	}

	baseURLRouting, ok := os.LookupEnv("BASE_URL_ROUTING")
	if !ok {
		return nil, errors.New("BASE_URL_ROUTING is not set")
	}

	var mainPort int

	port1 := os.Getenv("PORT")
	if port1 != "" {
		mainPort, _ = strconv.Atoi(port1)
	} else {
		mainPort = 8080
	}

	port2, ok := os.LookupEnv("PRICE_ESTIMATOR_PORT")
	if !ok {
		return nil, errors.New("PRICE_ESTIMATOR_PORT is not set")
	}
	priceEstimatorPort, _ := strconv.Atoi(port2)

	return &Config{
		PG: &PG{
			PostgresUser:     user,
			PostgresPassword: password,
			PostgresName:     name,
			PostgresPort:     port,
		},
		GEO: &GEO{
			APIKeyCatalog: apiKeyCatalog,
			APIKeyRouting: apiKeyRouting,

			BaseURLCatalog: baseURLCatalog,
			BaseURLRouting: baseURLRouting,
		},
		SERVICES: &SERVICES{
			Ports: map[string]int{
				"Main Application": mainPort,
				"PriceEstimator":   priceEstimatorPort,
			},
		},
		LOG: &LOG{
			LogrusFormatter: &logrus.TextFormatter{
				TimestampFormat:        "02-01-2006 15:04:05",
				FullTimestamp:          true,
				ForceColors:            true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					// Format string to get layer name
					i := strings.Index(f.File, "truckly/")
					layerPath, _ := strings.CutPrefix(f.File[i:], "truckly/")
					layerArr := strings.Split(layerPath, "/")
					layerMsg := fmt.Sprintf("level:%s/%s", layerArr[0], layerArr[1])

					// Split string to get file name
					pathArr := strings.Split(f.File, "/")
					fileName := pathArr[len(pathArr)-1]

					// Split string to get func name
					funcArr := strings.Split(f.Function, ".")
					funcName := funcArr[len(funcArr)-1]

					// Logger message
					var msg string
					if layerMsg != "level:internal/transport" && fileName != "logger.go" {
						msg = fmt.Sprintf("%30s | %s:%d | func:%s |", layerMsg, fileName, f.Line, funcName)
					} else {
						msg = fmt.Sprintf("%30s", layerMsg)
					}

					// Return info
					return "", msg
				},
			},
		},
	}, nil
}
