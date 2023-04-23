package app

import (
	"database/sql"
	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/dacore-x/truckly/pkg/pghelper"
	"github.com/dacore-x/truckly/pkg/redishelper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

	// Redis
	options := redishelper.GetOptions(cfg.REDIS)

	rdb := redis.NewClient(options)
	defer rdb.Close()

	// Use cases
	userUseCase := usecase.NewUserUseCase(
		postgres.NewUserRepo(conn, appLogger),
		appLogger,
	)

	metricsUseCase := usecase.NewMetricsUseCase(
		postgres.NewMetricsRepo(conn, appLogger),
		appLogger,
	)

	geoWebAPI := webapi.New(cfg.GEO, appLogger)
	priceEstimatorService := microservice.New(cfg.SERVICES, appLogger)

	deliveryUseCase := usecase.NewDeliveryUseCase(
		postgres.NewDeliveryRepo(conn, appLogger),
		geoWebAPI,
		priceEstimatorService,
		appLogger,
	)

	geoUseCase := usecase.NewGeoUseCase(geoWebAPI, appLogger)
	priceEstimatorUseCase := usecase.NewPriceEstimatorUseCase(priceEstimatorService, geoWebAPI, appLogger)

	// Create HTTP server using Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Setting cors settings
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	h := v1.NewHandlers(
		userUseCase,
		deliveryUseCase,
		metricsUseCase,
		geoUseCase,
		priceEstimatorUseCase,
		appLogger,
		rdb,
	)
	h.NewRouter(r)

	// Log all running services ports
	for k, v := range cfg.SERVICES.Ports {
		appLogger.Infof("Running %v on :%v", k, v)
	}
	r.Run()
}
