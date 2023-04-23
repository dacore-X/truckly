package middleware

import (
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/redis/go-redis/v9"

	"github.com/dacore-x/truckly/internal/usecase"
)

// Middlewares is a struct that provides
// all entities' middlewares
type Middlewares struct {
	userMiddlewares
	loggerMiddlewares
	redisMiddlewares
}

func New(u usecase.User, l *logger.Logger, rdb *redis.Client) *Middlewares {
	return &Middlewares{
		userMiddlewares{u},
		loggerMiddlewares{l},
		redisMiddlewares{rdb},
	}
}
