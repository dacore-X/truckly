package app

import (
	"database/sql"

	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/dacore-x/truckly/pkg/pghelper"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/dacore-x/truckly/internal/infrastructure/microservice"
	"github.com/dacore-x/truckly/internal/infrastructure/repository/postgres"
	"github.com/dacore-x/truckly/internal/infrastructure/webapi"
	v1 "github.com/dacore-x/truckly/internal/transport/http/v1"
	"github.com/dacore-x/truckly/internal/usecase"
)

func Run(cfg *config.Config) {
	// Create Logger using logrus
	logrusLogger := logrus.New()
	logrusLogger.SetReportCaller(true)
	logrusLogger.SetFormatter(cfg.LogrusFormatter)
	appLogger := logger.New(logrusLogger)

	// Repository
	connURL := pghelper.GetConnURL(cfg.PG)

	conn, err := sql.Open("postgres", connURL)
	if err != nil {
		appLogger.Fatal(err)
	}
	defer conn.Close()

	// Use cases
	userUseCase := usecase.NewUserUseCase(
		postgres.NewUserRepo(conn),
		appLogger,
	)

	geoWebAPI := webapi.New(cfg.GEO)
	priceEstimatorService := microservice.New(cfg.SERVICES)

	deliveryUseCase := usecase.NewDeliveryUseCase(
		postgres.NewDeliveryRepo(conn),
		geoWebAPI,
		appLogger,
	)

	geoUseCase := usecase.NewGeoUseCase(geoWebAPI, appLogger)
	priceEstimatorUseCase := usecase.NewPriceEstimatorUseCase(priceEstimatorService, appLogger)

	// Create HTTP server using Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	h := v1.NewHandlers(userUseCase, deliveryUseCase, geoUseCase, priceEstimatorUseCase, appLogger)
	h.NewRouter(r)

	// Log all running services ports
	for k, v := range cfg.SERVICES.Ports {
		appLogger.Infof("Running %v on :%v", k, v)
	}
	r.Run()
}
