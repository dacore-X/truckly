package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

type priceEstimatorHandlers struct {
	usecase.PriceEstimator
}

func newPriceEstimatorHandlers(superGroup *gin.RouterGroup, u usecase.PriceEstimator, m *middleware.Middlewares) {
	handler := &priceEstimatorHandlers{u}

	priceEstimatorGroup := superGroup.Group("/price")
	{
		priceEstimatorGroup.POST("/predict", m.RequireAuth, m.RequireNoBan, handler.estimatePrice)
	}
}

func (h *priceEstimatorHandlers) estimatePrice(c *gin.Context) {
	var body dto.EstimatePriceRequestBody
	if c.BindJSON(&body) != nil {
		err := fmt.Errorf("failed to read body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	price, err := h.EstimateDeliveryPrice(context.Background(), &body)
	if err != nil {
		err := fmt.Errorf("failed to estimate price")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := dto.EstimatePriceResponse{Price: price}
	c.JSON(http.StatusOK, resp)
}
