package usecase

import (
	"context"

	"github.com/dacore-x/truckly/internal/dto"
)

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) Create(ctx context.Context, req dto.UserRequestSignUpBody) error {
	err := uc.repo.Create(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) GetMe(ctx context.Context, id int64) (*dto.UserResponseMeBody, error) {
	resp, err := uc.repo.GetMe(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id int64) (*dto.UserResponseInfoBody, error) {
	resp, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *UserUseCase) GetByEmail(ctx context.Context, email string) (*dto.UserResponseInfoBody, error) {
	resp, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
