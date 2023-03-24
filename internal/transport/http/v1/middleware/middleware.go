package middleware

import "github.com/dacore-x/truckly/internal/usecase"

type Middlewares struct {
	userMiddlewares
}

func New(u usecase.User) *Middlewares {
	return &Middlewares{
		userMiddlewares{u},
	}
}
