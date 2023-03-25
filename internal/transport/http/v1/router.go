package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

// Handlers is a struct that provides
// all entities' handlers and middlewares
type Handlers struct {
	userHandlers
	*middleware.Middlewares
}

func NewHandlers(u usecase.User) *Handlers {
	return &Handlers{
		userHandlers{u},
		middleware.New(u),
	}
}

// NewRouter initializes a group of all entities' routes
func (h *Handlers) NewRouter(r *gin.Engine) {
	superGroup := r.Group("/api")
	{
		newUserHandlers(superGroup, h.userHandlers, h.Middlewares)
	}
}
