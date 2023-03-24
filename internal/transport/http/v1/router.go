package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

type Handlers struct {
	userHandlers
	middleware.Middlewares
}

func NewHandlers(u usecase.User) *Handlers {
	return &Handlers{
		userHandlers{u},
		middleware.Middlewares{},
	}
}

func (h *Handlers) NewRouter(r *gin.Engine) {
	// All entities' routes
	superGroup := r.Group("/api")
	{
		newUserHandlers(superGroup, h.userHandlers, h.Middlewares)
	}
}
