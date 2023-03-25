package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
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

// Config is a struct for storing all required configuration parameters
type Config struct {
	*PG
  *GEO
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
	}, nil
}
