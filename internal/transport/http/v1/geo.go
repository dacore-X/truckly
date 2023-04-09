package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

type geoHandlers struct {
	usecase.Geo
}

func newGeoHandlers(superGroup *gin.RouterGroup, u usecase.Geo, m *middleware.Middlewares) {
	handler := &geoHandlers{u}

	geoGroup := superGroup.Group("/geo")
	{
		geoGroup.GET("/coords", m.RequireAuth, m.RequireNoBan, handler.getCoordsByObject)
		geoGroup.GET("/object", m.RequireAuth, m.RequireNoBan, handler.getObjectByCoords)
	}
}

func (h *geoHandlers) getCoordsByObject(c *gin.Context) {
	q := c.Query("q")
	coords, err := h.GetCoordsByObject(context.Background(), q)
	if err != nil {
		err := fmt.Errorf("error finding geo object")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coords": coords,
	})
}

func (h *geoHandlers) getObjectByCoords(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")

	latConv, err := strconv.ParseFloat(lat, 64)

	if err != nil {
		err := fmt.Errorf("error converting latitude")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	lonConv, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		err := fmt.Errorf("error converting longitude")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	object, err := h.GetObjectByCoords(context.Background(), latConv, lonConv)
	if err != nil {
		err := fmt.Errorf("error finding geo object")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coords": object,
	})
}
