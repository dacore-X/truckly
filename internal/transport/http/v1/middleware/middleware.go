package middleware

import "github.com/dacore-x/truckly/internal/usecase"

// Middlewares is a struct that provides
// all entities' middlewares
type Middlewares struct {
	userMiddlewares
}

func New(u usecase.User) *Middlewares {
	return &Middlewares{
		userMiddlewares{u},
	}
}
