package pghelper

import (
	"fmt"

	"github.com/dacore-x/truckly/config"
)

// GetConnURL returns the Postgres connection URL using config
func GetConnURL(cfg *config.PG) string {
	return fmt.Sprintf("postgresql://%v:%v@localhost:%v/%v?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresPort, cfg.PostgresName)
}
