package v1

import (
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handlers is a struct that provides
// all entities' handlers and middlewares
type Handlers struct {
	userHandlers
	deliveryHandlers
	metricsHandlers
	geoHandlers
	priceEstimatorHandlers
	*middleware.Middlewares
}

func NewHandlers(
	u usecase.User,
	d usecase.Delivery,
	m usecase.Metrics,
	g usecase.Geo,
	p usecase.PriceEstimator,
	l *logger.Logger,
) *Handlers {
	return &Handlers{
		userHandlers{u},
		deliveryHandlers{d},
		metricsHandlers{m},
		geoHandlers{g},
		priceEstimatorHandlers{p},
		middleware.New(u, l),
	}
}

// NewRouter initializes a group of all entities' routes
func (h *Handlers) NewRouter(r *gin.Engine) {
	r.Use(h.DefaultLogger())
	r.Use(gin.Recovery())
	config := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET, POST, PUT, DELETE, HEAD, OPTIONS"},
		AllowHeaders:     []string{"Origin, X-Requested-With, Content-Type, Accept, Authorization"},
		AllowCredentials: true,
	})
	r.Use(config)
	superGroup := r.Group("/api")
	{
		newUserHandlers(superGroup, h.userHandlers, h.Middlewares)
		newDeliveryHandlers(superGroup, h.deliveryHandlers, h.Middlewares)
		newMetricsHandlers(superGroup, h.metricsHandlers, h.Middlewares)
		newGeoHandlers(superGroup, h.geoHandlers, h.Middlewares)
		newPriceEstimatorHandlers(superGroup, h.priceEstimatorHandlers, h.Middlewares)
	}
}
