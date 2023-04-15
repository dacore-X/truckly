package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dacore-x/truckly/pkg/logger"

	"github.com/dacore-x/truckly/internal/dto"
)

// UserRepo is a struct that provides
// all functions to execute SQL queries
// related to user's requests
type UserRepo struct {
	*sql.DB
	appLogger *logger.Logger
}

func NewUserRepo(db *sql.DB, l *logger.Logger) *UserRepo {
	return &UserRepo{db, l}
}

// CreateUser creates a new user record in the database with meta data attached to it
func (ur *UserRepo) CreateUser(ctx context.Context, req *dto.UserSignUpRequestBody) error {
	tx, err := ur.Begin()
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	defer tx.Rollback()

	query1 := `
		INSERT INTO users(surname, name, email, phone_number, hash_password)
		VALUES($1, $2, $3, $4, $5) 
		RETURNING id
	`
	lastInsertID := 0
	err = tx.QueryRowContext(
		ctx, query1, req.Surname, req.Name, req.Email, req.PhoneNumber, req.Password,
	).Scan(&lastInsertID)
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}

	query2 := `
		INSERT INTO meta(user_id, is_courier)
		VALUES($1, $2)
	`
	result, err := tx.ExecContext(ctx, query2, lastInsertID, req.IsCourier)
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	if rows != 1 {
		err := fmt.Errorf("expected to affect 1 row, affected %d", rows)
		ur.appLogger.Error(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		ur.appLogger.Error(err)
		return err
	}
	return nil
}

// BanUser updates user's is_banned field and sets its value to true
func (ur *UserRepo) BanUser(ctx context.Context, id int) error {
	query := `
		UPDATE meta
		SET is_banned=true
		WHERE user_id=$1
	`
	result, err := ur.ExecContext(ctx, query, id)
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	if rows != 1 {
		err := fmt.Errorf("expected to affect 1 row, affected %d", rows)
		ur.appLogger.Error(err)
		return err
	}
	return nil
}

// UnbanUser updates user's is_banned field and sets its value to false
func (ur *UserRepo) UnbanUser(ctx context.Context, id int) error {
	query := `
		UPDATE meta
		SET is_banned=false
		WHERE user_id=$1
	`
	result, err := ur.ExecContext(ctx, query, id)
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ur.appLogger.Error(err)
		return err
	}
	if rows != 1 {
		err := fmt.Errorf("expected to affect 1 row, affected %d", rows)
		ur.appLogger.Error(err)
		return err
	}
	return nil
}

// GetUserByID fetches user's account data from the database and returns it
func (ur *UserRepo) GetUserByID(ctx context.Context, id int) (*dto.UserMeResponse, error) {
	query := `
		SELECT users.id, surname, name, email, phone_number, created_at, is_admin, is_courier, is_banned
		FROM users INNER JOIN meta ON users.id = meta.user_id
		WHERE users.id=$1
	`
	row := ur.QueryRowContext(ctx, query, id)

	resp := &dto.UserMeResponse{
		Meta: &dto.RoleMeta{},
	}

	err := row.Scan(
		&resp.ID,
		&resp.Surname,
		&resp.Name,
		&resp.Email,
		&resp.PhoneNumber,
		&resp.CreatedAt,
		&resp.Meta.IsAdmin,
		&resp.Meta.IsCourier,
		&resp.Meta.IsBanned,
	)
	if err != nil {
		ur.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetUserPrivateByID fetches private user's data by id from the database and returns it
func (ur *UserRepo) GetUserPrivateByID(ctx context.Context, id int) (*dto.UserInfoResponse, error) {
	query := `
		SELECT id, email, hash_password
		FROM users
		WHERE id=$1
	`
	row := ur.QueryRowContext(ctx, query, id)

	resp := &dto.UserInfoResponse{}
	err := row.Scan(&resp.ID, &resp.Email, &resp.Password)
	if err != nil {
		ur.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetUserPrivateByEmail fetches private user's data by email from the database and returns it
func (ur *UserRepo) GetUserPrivateByEmail(ctx context.Context, email string) (*dto.UserInfoResponse, error) {
	query := `
		SELECT id, email, hash_password
		FROM users
		WHERE email=$1
	`
	row := ur.QueryRowContext(ctx, query, email)

	resp := &dto.UserInfoResponse{}
	err := row.Scan(&resp.ID, &resp.Email, &resp.Password)
	if err != nil {
		ur.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}

// GetUserMeta fetches user's metadata by id from the database and returns it
func (ur *UserRepo) GetUserMeta(ctx context.Context, id int) (*dto.UserMetaResponse, error) {
	query := `
		SELECT user_id, is_admin, is_courier, is_banned, rating
		FROM meta
		WHERE user_id=$1
	`
	row := ur.QueryRowContext(ctx, query, id)

	resp := &dto.UserMetaResponse{}
	err := row.Scan(&resp.UserID, &resp.IsAdmin, &resp.IsCourier, &resp.IsBanned, &resp.Rating)
	if err != nil {
		ur.appLogger.Error(err)
		return nil, err
	}
	return resp, nil
}
