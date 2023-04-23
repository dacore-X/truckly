package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/go-test/deep"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/dacore-x/truckly/internal/dto"
)

func TestPostgres_CreateUserTxNoRollback(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		id     int
		body   *dto.UserSignUpRequestBody
		rows   *sqlmock.Rows
		result driver.Result
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "user usual case",
			args: args{
				id: 1,
				body: &dto.UserSignUpRequestBody{
					Surname:     "Николаев",
					Name:        "Николай",
					Email:       "nikolaev@mail.ru",
					PhoneNumber: "89157674599",
					Password:    "password123",
					IsCourier:   false,
				},
				rows:   sqlmock.NewRows([]string{"id"}).AddRow(1),
				result: sqlmock.NewResult(1, 1),
			},
		},
		{
			name: "courier usual case",
			args: args{
				id: 2,
				body: &dto.UserSignUpRequestBody{
					Surname:     "Дмитриев",
					Name:        "Дмитрий",
					Email:       "dmitriev@bk.ru",
					PhoneNumber: "89850357751",
					Password:    "mypassword",
					IsCourier:   true,
				},
				rows:   sqlmock.NewRows([]string{"id"}).AddRow(2),
				result: sqlmock.NewResult(2, 1),
			},
		},
		{
			name: "user usual case",
			args: args{
				id: 3,
				body: &dto.UserSignUpRequestBody{
					Surname:     "Романов",
					Name:        "Роман",
					Email:       "romanov@yandex.ru",
					PhoneNumber: "89456041022",
					Password:    "pswd1414",
					IsCourier:   false,
				},
				rows:   sqlmock.NewRows([]string{"id"}).AddRow(3),
				result: sqlmock.NewResult(3, 1),
			},
		},
		{
			name: "user usual case",
			args: args{
				id: 4,
				body: &dto.UserSignUpRequestBody{
					Surname:     "Эдуардов",
					Name:        "Эдуард",
					Email:       "eduardov@yandex.ru",
					PhoneNumber: "89406432826",
					Password:    "eduardpassword",
					IsCourier:   false,
				},
				rows:   sqlmock.NewRows([]string{"id"}).AddRow(4),
				result: sqlmock.NewResult(4, 1),
			},
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect transaction begin
			mock.ExpectBegin()

			// Expect query to create a new user record, return last inserted id
			// and either return error or not, match it with regexp
			mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO users(surname, name, email, phone_number, hash_password)
					VALUES($1, $2, $3, $4, $5) 
					RETURNING id
				`)).
				WithArgs(
					tc.args.body.Surname,
					tc.args.body.Name,
					tc.args.body.Email,
					tc.args.body.PhoneNumber,
					tc.args.body.Password).
				WillReturnRows(tc.args.rows).
				WillReturnError(tc.wantErr)

			// Expect query to create a new user metadata record
			// and either return error or not, match it with regexp
			mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO meta(user_id, is_courier)
					VALUES($1, $2) 
				`)).
				WithArgs(tc.args.id, tc.args.body.IsCourier).
				WillReturnResult(tc.args.result)

			// Expect transaction commit
			mock.ExpectCommit()

			// Run the create transaction function
			err := repo.CreateUser(context.Background(), tc.args.body)
			require.Nil(t, deep.Equal(tc.wantErr, err))
		})
	}
}

func TestPostgres_CreateUserTxWithRollback(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		id   int
		body *dto.UserSignUpRequestBody
		rows *sqlmock.Rows
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "no return id value in first query",
			args: args{
				id: 1,
				body: &dto.UserSignUpRequestBody{
					Surname:     "Карпов",
					Name:        "Александр",
					Email:       "karpov@inbox.ru",
					PhoneNumber: "89701010166",
					Password:    "karpovpassword",
					IsCourier:   false,
				},
				rows: sqlmock.NewRows([]string{"id"}).AddRow(nil),
			},
			wantErr: sql.ErrNoRows,
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect transaction begin
			mock.ExpectBegin()

			// Expect query to create a new user record, return last inserted id
			// and either return error or not, match it with regexp
			mock.ExpectQuery(regexp.QuoteMeta(`
				INSERT INTO users(surname, name, email, phone_number, hash_password)
				VALUES($1, $2, $3, $4, $5) 
				RETURNING id
			`)).
				WithArgs(
					tc.args.body.Surname,
					tc.args.body.Name,
					tc.args.body.Email,
					tc.args.body.PhoneNumber,
					tc.args.body.Password).
				WillReturnRows(tc.args.rows).
				WillReturnError(tc.wantErr)

			// Expect transaction rollback
			mock.ExpectRollback()

			// Run the create transaction function
			err := repo.CreateUser(context.Background(), tc.args.body)
			require.Nil(t, deep.Equal(tc.wantErr, err))
		})
	}
}

func TestPostgres_BanUser(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		id     int
		result driver.Result
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "ban existing user",
			args: args{
				id:     1,
				result: sqlmock.NewResult(0, 1),
			},
		},
		{
			name: "ban non-existent user",
			args: args{
				id:     3,
				result: sqlmock.NewResult(0, 0),
			},
			wantErr: fmt.Errorf("user is not found"),
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query to update users's ban status and
			// either return error or not, match it with regexp
			mock.ExpectExec(regexp.QuoteMeta(`
				UPDATE meta
				SET is_banned=true
				WHERE user_id=$1
			`)).
				WithArgs(tc.args.id).
				WillReturnResult(tc.args.result).
				WillReturnError(tc.wantErr)

			// Run the ban function
			err := repo.BanUser(context.Background(), tc.args.id)
			require.Nil(t, deep.Equal(tc.wantErr, err))
		})
	}
}

func TestPostgres_UnbanUser(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		id     int
		result driver.Result
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "unban existing user",
			args: args{
				id:     1,
				result: sqlmock.NewResult(0, 1),
			},
		},
		{
			name: "unban non-existent user",
			args: args{
				id:     3,
				result: sqlmock.NewResult(0, 0),
			},
			wantErr: fmt.Errorf("user is not found"),
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query to update users's ban status and
			// either return error or not, match it with regexp
			mock.ExpectExec(regexp.QuoteMeta(`
				UPDATE meta
				SET is_banned=false
				WHERE user_id=$1
			`)).
				WithArgs(tc.args.id).
				WillReturnResult(tc.args.result).
				WillReturnError(tc.wantErr)

			// Run the ban function
			err := repo.UnbanUser(context.Background(), tc.args.id)
			require.Nil(t, deep.Equal(tc.wantErr, err))
		})
	}
}

func TestPostgres_GetUserByID(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// time.Time variable for tests
	now := time.Now()

	// Required args for tests
	type args struct {
		id   int
		rows *sqlmock.Rows
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		want    *dto.UserMeResponse
		wantErr error
	}{
		{
			name: "user is found",
			args: args{
				id: 1,
				rows: sqlmock.NewRows([]string{"id", "surname", "name", "email", "phone_number", "created_at", "is_admin", "is_courier", "is_banned"}).
					AddRow(1, "Иванов", "Иван", "ivanov@yandex.ru", "89157650030", now, false, false, false),
			},
			want: &dto.UserMeResponse{
				ID:          1,
				Surname:     "Иванов",
				Name:        "Иван",
				Email:       "ivanov@yandex.ru",
				PhoneNumber: "89157650030",
				CreatedAt:   now,
				Meta: &dto.RoleMeta{
					IsAdmin:   false,
					IsCourier: false,
					IsBanned:  false,
				},
			},
		},
		{
			name: "user is not found",
			args: args{
				id: 2,
				rows: sqlmock.NewRows([]string{"id", "surname", "name", "email", "phone_number", "created_at", "is_admin", "is_courier", "is_banned"}).
					AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil),
			},
			wantErr: sql.ErrNoRows,
		},
	}
	
	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query to fetch user's account data and
			// either return error or not, match it with regexp
			mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT users.id, surname, name, email, phone_number, created_at, is_admin, is_courier, is_banned
				FROM users INNER JOIN meta ON users.id = meta.user_id
				WHERE users.id=$1
			`)).
				WithArgs(tc.args.id).
				WillReturnRows(tc.args.rows).
				WillReturnError(tc.wantErr)

			// Run the get me function
			result, err := repo.GetUserByID(context.Background(), tc.args.id)
			require.Nil(t, deep.Equal(tc.wantErr, err))
			require.Nil(t, deep.Equal(tc.want, result))
		})
	}
}

func TestPostgres_GetUserPrivateByID(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		id   int
		rows *sqlmock.Rows
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		want    *dto.UserInfoResponse
		wantErr error
	}{
		{
			name: "user is found",
			args: args{
				id: 1,
				rows: sqlmock.NewRows([]string{"id", "email", "hash_password"}).
					AddRow(1, "sergeev@yahoo.com", "123456789"),
			},
			want: &dto.UserInfoResponse{
				ID:       1,
				Email:    "sergeev@yahoo.com",
				Password: "123456789",
			},
		},
		{
			name: "user is not found",
			args: args{
				id: 7,
				rows: sqlmock.NewRows([]string{"id", "email", "hash_password"}).
					AddRow(nil, nil, nil),
			},
			wantErr: sql.ErrNoRows,
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query to fetch private user's data by id and
			// either return error or not, match it with regexp
			mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT id, email, hash_password
				FROM users
				WHERE id=$1
			`)).
				WithArgs(tc.args.id).
				WillReturnRows(tc.args.rows).
				WillReturnError(tc.wantErr)

			// Run the get by id function
			result, err := repo.GetUserPrivateByID(context.Background(), tc.args.id)
			require.Nil(t, deep.Equal(tc.wantErr, err))
			require.Nil(t, deep.Equal(tc.want, result))
		})
	}
}

func TestPostgres_GetUserPrivateByEmail(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		email string
		rows  *sqlmock.Rows
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		want    *dto.UserInfoResponse
		wantErr error
	}{
		{
			name: "user is found",
			args: args{
				email: "ivanov@gmail.com",
				rows: sqlmock.NewRows([]string{"id", "email", "hash_password"}).
					AddRow(1, "ivanov@gmail.com", "password"),
			},
			want: &dto.UserInfoResponse{
				ID:       1,
				Email:    "ivanov@gmail.com",
				Password: "password",
			},
		},
		{
			name: "user is not found",
			args: args{
				email: "kuznetsov@outlook.com",
				rows: sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(nil, nil, nil),
			},
			wantErr: sql.ErrNoRows,
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query to fetch private user's data by email and
			// either return error or not, match it with regexp
			mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT id, email, hash_password
				FROM users
				WHERE email=$1
			`)).
				WithArgs(tc.args.email).
				WillReturnRows(tc.args.rows).
				WillReturnError(tc.wantErr)

			// Run the get by email function
			result, err := repo.GetUserPrivateByEmail(context.Background(), tc.args.email)
			require.Nil(t, deep.Equal(tc.wantErr, err))
			require.Nil(t, deep.Equal(tc.want, result))
		})
	}
}

func TestUser_GetUserMeta(t *testing.T) {
	// Open database stub
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repo
	testLogger := logrus.New()
	repo := NewUserRepo(db, logger.New(testLogger))

	// Required args for tests
	type args struct {
		id   int
		rows *sqlmock.Rows
	}

	// Slice of test cases
	cases := []struct {
		name    string
		args    args
		want    *dto.UserMetaResponse
		wantErr error
	}{
		{
			name: "default user",
			args: args{
				id: 1,
				rows: sqlmock.NewRows([]string{"user_id", "is_admin", "is_courier", "is_banned", "rating"}).
					AddRow(1, false, false, false, 4.00),
			},
			want: &dto.UserMetaResponse{
				UserID:    1,
				IsAdmin:   false,
				IsCourier: false,
				IsBanned:  false,
				Rating:    4.00,
			},
		},
		{
			name: "banned user",
			args: args{
				id: 2,
				rows: sqlmock.NewRows([]string{"user_id", "is_admin", "is_courier", "is_banned", "rating"}).
					AddRow(2, false, false, true, 3.00),
			},
			want: &dto.UserMetaResponse{
				UserID:    2,
				IsAdmin:   false,
				IsCourier: false,
				IsBanned:  true,
				Rating:    3.00,
			},
		},
		{
			name: "admin user",
			args: args{
				id: 3,
				rows: sqlmock.NewRows([]string{"user_id", "is_admin", "is_courier", "is_banned", "rating"}).
					AddRow(3, true, false, false, 2.00),
			},
			want: &dto.UserMetaResponse{
				UserID:    3,
				IsAdmin:   true,
				IsCourier: false,
				IsBanned:  false,
				Rating:    2.00,
			},
		},
		{
			name: "courier user",
			args: args{
				id: 4,
				rows: sqlmock.NewRows([]string{"user_id", "is_admin", "is_courier", "is_banned", "rating"}).
					AddRow(4, false, true, false, 5.00),
			},
			want: &dto.UserMetaResponse{
				UserID:    4,
				IsAdmin:   false,
				IsCourier: true,
				IsBanned:  false,
				Rating:    5.00,
			},
		},
		{
			name: "user is not found",
			args: args{
				id: 5,
				rows: sqlmock.NewRows([]string{"user_id", "is_admin", "is_courier", "is_banned", "rating"}).
					AddRow(nil, nil, nil, nil, nil),
			},
			wantErr: sql.ErrNoRows,
		},
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query to fetch private user's data by email and
			// either return error or not, match it with regexp
			mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT user_id, is_admin, is_courier, is_banned, rating
				FROM meta
				WHERE user_id=$1
			`)).
				WithArgs(tc.args.id).
				WillReturnRows(tc.args.rows).
				WillReturnError(tc.wantErr)

			// Run the get by email function
			result, err := repo.GetUserMeta(context.Background(), tc.args.id)
			require.Nil(t, deep.Equal(tc.wantErr, err))
			require.Nil(t, deep.Equal(tc.want, result))
		})
	}
}
