package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dacore-x/truckly/internal/dto"
)

type UserRepo struct {
	*sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) Create(ctx context.Context, req dto.UserRequestSignUpBody) error {
	query1 := `
		INSERT INTO users(surname, name, patronymic, email, phone_number, hash_password)
		VALUES($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`
	lastInsertID := 0
	err := ur.QueryRowContext(ctx, query1, req.Surname, req.Name, req.Patronymic, req.Email, req.PhoneNumber, req.Password).Scan(&lastInsertID)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO meta(user_id, is_courier)
		VALUES($1, $2)
	`
	result, err := ur.ExecContext(ctx, query2, lastInsertID, req.IsCourier)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (ur *UserRepo) GetMe(ctx context.Context, id int64) (*dto.UserResponseMeBody, error) {
	query := `
		SELECT id, surname, name, patronymic, email, phone_number, created_at
		FROM users
		WHERE id=$1
	`
	row := ur.QueryRowContext(ctx, query, id)

	resp := &dto.UserResponseMeBody{}
	err := row.Scan(&resp.ID, &resp.Surname, &resp.Name, &resp.Patronymic, &resp.Email, &resp.PhoneNumber, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id int64) (*dto.UserResponseInfoBody, error) {
	query := `
		SELECT id, email, hash_password
		FROM users
		WHERE id=$1
	`
	row := ur.QueryRowContext(ctx, query, id)

	resp := &dto.UserResponseInfoBody{}
	err := row.Scan(&resp.ID, &resp.Email, &resp.Password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*dto.UserResponseInfoBody, error) {
	query := `
		SELECT id, email, hash_password
		FROM users
		WHERE email=$1
	`
	row := ur.QueryRowContext(ctx, query, email)

	resp := &dto.UserResponseInfoBody{}
	err := row.Scan(&resp.ID, &resp.Email, &resp.Password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
