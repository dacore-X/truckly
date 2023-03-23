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
		INSERT INTO users(name, surname, patronymic, email, phone_number, hash_password)
		VALUES($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`
	lastInsertID := 0
	err := ur.QueryRowContext(ctx, query1, req.Name, req.Surname, req.Patronymic, req.Email, req.PhoneNumber, req.Password).Scan(&lastInsertID)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO meta(userr_id, is_courier)
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
