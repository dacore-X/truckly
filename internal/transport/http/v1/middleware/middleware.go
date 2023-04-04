package middleware

import (
	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/usecase"
)

// Middlewares is a struct that provides
// all entities' middlewares
type Middlewares struct {
	userMiddlewares
	loggerMiddlewares
}

func New(u usecase.User, l *logger.Logger) *Middlewares {
	return &Middlewares{
		userMiddlewares{u},
		loggerMiddlewares{l},
	}
}
