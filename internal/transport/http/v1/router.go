package v1

import (
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
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
	superGroup := r.Group("/api")
	{
		newUserHandlers(superGroup, h.userHandlers, h.Middlewares)
		newDeliveryHandlers(superGroup, h.deliveryHandlers, h.Middlewares)
		newMetricsHandlers(superGroup, h.metricsHandlers, h.Middlewares)
		newGeoHandlers(superGroup, h.geoHandlers, h.Middlewares)
		newPriceEstimatorHandlers(superGroup, h.priceEstimatorHandlers, h.Middlewares)
	}
}
