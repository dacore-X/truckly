package v1

import (
	"context"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error finding geo object",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coords": coords,
	})
	return

}

func (h *geoHandlers) getObjectByCoords(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")

	latConv, err := strconv.ParseFloat(lat, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error converting latitude",
		})
		return
	}

	lonConv, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error converting longitude",
		})
		return
	}
	object, err := h.GetObjectByCoords(context.Background(), latConv, lonConv)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error finding geo object",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coords": object,
	})
	return
}
