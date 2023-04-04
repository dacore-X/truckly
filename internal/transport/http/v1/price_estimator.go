package v1

import (
	"context"
	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type priceEstimatorHandlers struct {
	usecase.PriceEstimator
}

func newPriceEstimatorHandlers(superGroup *gin.RouterGroup, u usecase.PriceEstimator, m *middleware.Middlewares) {
	handler := &priceEstimatorHandlers{u}

	priceEstimatorGroup := superGroup.Group("/price")
	{
		priceEstimatorGroup.POST("/", m.RequireAuth, m.RequireNoBan, handler.estimatePrice)
	}
}

func (h *priceEstimatorHandlers) estimatePrice(c *gin.Context) {
	var body dto.EstimatePriceRequestBody
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}
	price, err := h.EstimateDeliveryPrice(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to estimate price",
		})
		return
	}

	resp := dto.EstimatePriceResponse{Price: price}
	c.JSON(http.StatusOK, resp)
}
