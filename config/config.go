package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// PG is a struct for storing Postgres connection settings
type PG struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_NAME     string
	POSTGRES_PORT     string
}

// Config is a struct for storing all required configuration parameters
type Config struct {
	*PG
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
		return nil, errors.New("POSTGRES_USER is not set")
	}

	name, ok := os.LookupEnv("POSTGRES_NAME")
	if !ok {
		return nil, errors.New("POSTGRES_NAME is not set")
	}

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, errors.New("POSTGRES_PORT is not set")
	}

	return &Config{
		PG: &PG{
			POSTGRES_USER:     user,
			POSTGRES_PASSWORD: password,
			POSTGRES_NAME:     name,
			POSTGRES_PORT:     port,
		},
	}, nil
}
