package usecase

import (
	"context"

	"github.com/dacore-x/truckly/internal/dto"
)

type (
	// UseCase User interface
	User interface {
		Create(context.Context, dto.UserRequestSignUpBody) error
		GetMe(context.Context, int64) (*dto.UserResponseMeBody, error)
		GetByID(context.Context, int64) (*dto.UserResponseInfoBody, error)
		GetByEmail(context.Context, string) (*dto.UserResponseInfoBody, error)
	}

	// Repository User interface
	UserRepo interface {
		Create(context.Context, dto.UserRequestSignUpBody) error
		GetMe(context.Context, int64) (*dto.UserResponseMeBody, error)
		GetByID(context.Context, int64) (*dto.UserResponseInfoBody, error)
		GetByEmail(context.Context, string) (*dto.UserResponseInfoBody, error)
	}
)
