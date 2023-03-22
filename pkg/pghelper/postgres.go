package pghelper

import (
	"fmt"

	"github.com/dacore-x/truckly/config"
)

func GetConnURL(cfg *config.PG) string {
	return fmt.Sprintf("postgresql://%v:%v@localhost:%v/%v?sslmode=disable", cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_PORT, cfg.POSTGRES_NAME)
}
