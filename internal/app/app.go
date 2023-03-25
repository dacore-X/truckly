package app

import (
	"database/sql"
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

	// HTTP server
	r := gin.Default()
	h := v1.NewHandlers(userUseCase)
	h.NewRouter(r)
	r.Run()
}
