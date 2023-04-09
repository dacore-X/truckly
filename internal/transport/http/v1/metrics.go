package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

// metricsHandlers is a non-exportable struct
// that provides metrics-related handlers
type metricsHandlers struct {
	usecase.Metrics
}

// newMetricsHandlers initializes a group of metrics' routes
func newMetricsHandlers(superGroup *gin.RouterGroup, u usecase.Metrics, m *middleware.Middlewares) {
	handler := &metricsHandlers{u}

	metricsGroup := superGroup.Group("/metrics")
	{
		metricsGroup.GET("/", m.RequireAuth, m.RequireNoBan, m.RequireAdmin, handler.metricsPerDay)
	}
}

// metricsPerDay handler gets all metrics per last 24 hours
func (h *metricsHandlers) metricsPerDay(c *gin.Context) {
	metrics, err := h.GetMetrics(context.Background())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}
