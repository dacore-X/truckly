package v1

import (
	"context"
	"log"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}
	log.Println(body)
	clientID := c.GetInt("user")
	geo := &entity.Geo{
		FromLongitude: body.FromLongitude,
		FromLatitude:  body.FromLatitude,
		ToLongitude:   body.ToLongitude,
		ToLatitude:    body.ToLatitude,
	}
	delivery := &entity.Delivery{
		ClientID:  clientID,
		TypeID:    body.TypeID,
		Geo:       geo,
		HasLoader: body.HasLoader,
	}

	err := h.CreateDelivery(context.Background(), delivery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "delivery created successfully",
	})
}
