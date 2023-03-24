package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
)

type UserUseCase interface {
	Create(context.Context, dto.UserRequestSignUpBody) error
	GetMe(context.Context, int64) (*dto.UserResponseMeBody, error)
	GetByID(context.Context, int64) (*dto.UserResponseInfoBody, error)
	GetByEmail(context.Context, string) (*dto.UserResponseInfoBody, error)
}

type userHandlers struct {
	UserUseCase
}

func newUserHandlers(superGroup *gin.RouterGroup, u UserUseCase, m middleware.Middlewares) {
}
