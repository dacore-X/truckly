package usecase

import (
	"context"

	"github.com/dacore-x/truckly/internal/dto"
)

// UserUseCase is a struct that provides
// all user's usecases
type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

// Create usecase creates new user account
func (uc *UserUseCase) CreateTx(ctx context.Context, req *dto.UserSignUpRequestBody) error {
	err := uc.repo.CreateTx(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// GetMe usecase gets user's account data from storage based on user's id
func (uc *UserUseCase) GetMe(ctx context.Context, id int) (*dto.UserMeResponse, error) {
	resp, err := uc.repo.GetMe(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetByID usecase gets private user's data from storage based on user's id
func (uc *UserUseCase) GetByID(ctx context.Context, id int) (*dto.UserInfoResponse, error) {
	resp, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetByEmail usecase gets private user's data from storage based on user's email
func (uc *UserUseCase) GetByEmail(ctx context.Context, email string) (*dto.UserInfoResponse, error) {
	resp, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
