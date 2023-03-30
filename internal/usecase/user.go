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

// CreateUserTx usecase creates new user account
func (uc *UserUseCase) CreateUserTx(ctx context.Context, req *dto.UserSignUpRequestBody) error {
	err := uc.repo.CreateUserTx(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// BanUser usecase changes user's ban status to banned
func (uc *UserUseCase) BanUser(ctx context.Context, id int) error {
	err := uc.repo.BanUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// UnbanUser usecase changes user's ban status to unbanned
func (uc *UserUseCase) UnbanUser(ctx context.Context, id int) error {
	err := uc.repo.UnbanUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID usecase gets user's account data from storage based on user's id
func (uc *UserUseCase) GetUserByID(ctx context.Context, id int) (*dto.UserMeResponse, error) {
	resp, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserPrivateByID usecase gets private user's data from storage based on user's id
func (uc *UserUseCase) GetUserPrivateByID(ctx context.Context, id int) (*dto.UserInfoResponse, error) {
	resp, err := uc.repo.GetUserPrivateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserPrivateByEmail usecase gets private user's data from storage based on user's email
func (uc *UserUseCase) GetUserPrivateByEmail(ctx context.Context, email string) (*dto.UserInfoResponse, error) {
	resp, err := uc.repo.GetUserPrivateByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserMeta usecase gets user's metadata from storage based on user's id
func (uc *UserUseCase) GetUserMeta(ctx context.Context, id int) (*dto.UserMetaResponse, error) {
	resp, err := uc.repo.GetUserMeta(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
