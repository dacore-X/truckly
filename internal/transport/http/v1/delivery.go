package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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
		deliveryGroup.GET("/search", m.RequireAuth, m.RequireNoBan, m.RequireCourier, handler.getDeliveriesByGeolocation)
		deliveryGroup.GET("/", m.RequireAuth, m.RequireNoBan, m.RequireCourier, handler.getDeliveriesByCourierID)
		deliveryGroup.GET("/:id", m.RequireAuth, m.RequireNoBan, handler.getDeliveryByID)
		deliveryGroup.GET("/my", m.RequireAuth, m.RequireNoBan, handler.getDeliveriesByClientID)
		deliveryGroup.POST("/", m.RequireAuth, m.RequireNoBan, m.RateLimit, handler.createDelivery)
		deliveryGroup.POST("/:id/accept", m.RequireAuth, m.RequireNoBan, m.RequireCourier, handler.acceptDelivery)
		deliveryGroup.POST("/:id/cancel", m.RequireAuth, m.RequireNoBan, handler.cancelDelivery)
		deliveryGroup.POST("/:id/status", m.RequireAuth, m.RequireNoBan, m.RequireCourier, handler.changeDeliveryStatus)
	}
}

func (h *deliveryHandlers) createDelivery(c *gin.Context) {
	var body dto.DeliveryCreateBody
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

func (h *deliveryHandlers) getDeliveryByID(c *gin.Context) {
	// Get id of delivery from request
	var req dto.DeliveryIdURI
	if c.ShouldBindUri(&req) != nil {
		err := fmt.Errorf("failed to read uri")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clientID := c.GetInt("user")
	delivery, err := h.GetDeliveryByID(context.Background(), clientID, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, delivery)

}

func (h *deliveryHandlers) acceptDelivery(c *gin.Context) {
	// Get id of delivery from request
	var req dto.DeliveryIdURI
	if c.ShouldBindUri(&req) != nil {
		err := fmt.Errorf("failed to read uri")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	courierID := c.GetInt("user")
	err := h.AcceptDelivery(context.Background(), courierID, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "delivery is accepted",
	})
}

func (h *deliveryHandlers) changeDeliveryStatus(c *gin.Context) {
	// Get id of delivery from request
	var req dto.DeliveryIdURI
	if c.ShouldBindUri(&req) != nil {
		err := fmt.Errorf("failed to read uri")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body dto.DeliveryStatusChangeBody
	if c.BindJSON(&body) != nil {
		err := fmt.Errorf("failed to read body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	courierID := c.GetInt("user")
	err := h.ChangeDeliveryStatus(context.Background(), courierID, req.ID, body.StatusID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "delivery status is changed",
	})
}

func (h *deliveryHandlers) cancelDelivery(c *gin.Context) {
	// Get id of delivery from request
	var req dto.DeliveryIdURI
	if c.ShouldBindUri(&req) != nil {
		err := fmt.Errorf("failed to read uri")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clientID := c.GetInt("user")
	err := h.CancelDelivery(context.Background(), clientID, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "delivery status is changed",
	})
}

func (h *deliveryHandlers) getDeliveriesByGeolocation(c *gin.Context) {
	var q dto.DeliveryListGeolocationQuery
	if c.ShouldBindQuery(&q) != nil {
		err := fmt.Errorf("failed to read query")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	results, err := h.GetDeliveriesByGeolocation(context.Background(), &q)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *deliveryHandlers) getDeliveriesByClientID(c *gin.Context) {
	userID := c.GetInt("user")
	page := c.Query("page")
	if page == "" {
		err := fmt.Errorf("page is required")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		err = fmt.Errorf("bad request format")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	results, err := h.GetDeliveriesByClientID(context.Background(), userID, p)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *deliveryHandlers) getDeliveriesByCourierID(c *gin.Context) {
	courierID := c.GetInt("user")

	page := c.Query("page")
	if page == "" {
		err := fmt.Errorf("page is required")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		err = fmt.Errorf("bad request format")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	results, err := h.GetDeliveriesByCourierID(context.Background(), courierID, p)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}
