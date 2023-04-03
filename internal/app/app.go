package app

import (
	"database/sql"
	"github.com/dacore-x/truckly/internal/infrastructure/webapi"
	"log"

	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/pkg/pghelper"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/dacore-x/truckly/internal/infrastructure/repository/postgres"
	v1 "github.com/dacore-x/truckly/internal/transport/http/v1"
	"github.com/dacore-x/truckly/internal/usecase"
)

func Run(cfg *config.Config) {
	// Repository
	connURL := pghelper.GetConnURL(cfg.PG)

	conn, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Use cases
	userUseCase := usecase.NewUserUseCase(
		postgres.NewUserRepo(conn),
	)
	deliveryUseCase := usecase.NewDeliveryUseCase(
		postgres.NewDeliveryRepo(conn),
		webapi.New(cfg.GEO),
	)

	// HTTP server
	r := gin.Default()
	h := v1.NewHandlers(userUseCase, deliveryUseCase)
	h.NewRouter(r)
	r.Run()
}
