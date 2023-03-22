package config

type PG struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_NAME     string
	POSTGRES_PORT     string
}

type Config struct {
	*PG
}