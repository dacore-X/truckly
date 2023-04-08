package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/entity"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

type deliveryHandlers struct {
	usecase.Delivery
}

func newDeliveryHandlers(superGroup *gin.RouterGroup, u usecase.Delivery, m *middleware.Middlewares) {
	handler := &deliveryHandlers{u}

	deliveryGroup := superGroup.Group("/delivery")
	{
		deliveryGroup.POST("/", m.RequireAuth, m.RequireNoBan, handler.createDelivery)
	}
}

func (h *deliveryHandlers) createDelivery(c *gin.Context) {
	var body dto.DeliveryCreateRequestBody
	if c.BindJSON(&body) != nil {
		err := fmt.Errorf("failed to read body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clientID := c.GetInt("user")
	geo := &entity.Geo{
		FromLongitude: body.FromPoint.Lon,
		FromLatitude:  body.FromPoint.Lat,
		ToLongitude:   body.ToPoint.Lon,
		ToLatitude:    body.ToPoint.Lat,
	}

	delivery := &entity.Delivery{
		ClientID:  clientID,
		TypeID:    body.TypeID,
		Geo:       geo,
		HasLoader: body.HasLoader,
	}

	err := h.CreateDelivery(context.Background(), delivery)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "delivery created successfully",
	})
}
