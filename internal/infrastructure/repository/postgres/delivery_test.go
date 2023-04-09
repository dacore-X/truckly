package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dacore-x/truckly/internal/entity"
	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestDeliveryRepo_CreateDeliverySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	testLogger := logrus.New()
	repo := NewDeliveryRepo(db, logger.New(testLogger))

	type args struct {
		ctx      context.Context
		delivery *entity.Delivery
	}
	tests := []struct {
		name   string
		args   args
		rows   *sqlmock.Rows
		result driver.Result
		error  error
	}{
		{
			name: "delivery usual case",
			args: args{
				ctx: context.Background(),
				delivery: &entity.Delivery{
					ID:       1,
					ClientID: 1,
					StatusID: 1,
					TypeID:   1,
					Geo: &entity.Geo{
						FromLongitude: 37.22,
						FromLatitude:  55.77,
						FromObject:    "улица веселая д.1",
						ToLongitude:   37.22,
						ToLatitude:    55.77,
						ToObject:      "улица веселая д.10",
						Distance:      1200,
					},
					Price:     1220,
					HasLoader: true,
				},
			},
			rows:   sqlmock.NewRows([]string{"id"}).AddRow(1),
			result: sqlmock.NewResult(1, 1),
		},
		{
			name: "delivery usual case without type id",
			args: args{
				ctx: context.Background(),
				delivery: &entity.Delivery{
					ID:       2,
					ClientID: 1,
					StatusID: 1,
					//TypeID:   1,
					Geo: &entity.Geo{
						FromLongitude: 37.22,
						FromLatitude:  55.77,
						FromObject:    "улица веселая д.2",
						ToLongitude:   37.22,
						ToLatitude:    55.77,
						ToObject:      "улица веселая д.10",
						Distance:      1200,
					},
					Price:     1320,
					HasLoader: false,
				},
			},
			rows:   sqlmock.NewRows([]string{"id"}).AddRow(2),
			result: sqlmock.NewResult(2, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()

			mock.ExpectQuery(regexp.QuoteMeta(`
				INSERT INTO geo(from_longitude, from_latitude, from_object, to_longitude, to_latitude, to_object, distance)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
			`)).
				WithArgs(
					tt.args.delivery.Geo.FromLongitude,
					tt.args.delivery.Geo.FromLatitude,
					tt.args.delivery.Geo.FromObject,
					tt.args.delivery.Geo.ToLongitude,
					tt.args.delivery.Geo.ToLatitude,
					tt.args.delivery.Geo.ToObject,
					tt.args.delivery.Geo.Distance,
				).
				WillReturnRows(tt.rows).
				WillReturnError(tt.error)

			mock.ExpectExec(regexp.QuoteMeta(`
				INSERT INTO deliveries(client_id, status_id, type_id, geo_id, price, has_loader)
				VALUES ($1, $2, $3, $4, $5, $6)
			`)).
				WithArgs(tt.args.delivery.ClientID, 1, tt.args.delivery.TypeID, tt.args.delivery.ID, tt.args.delivery.Price, tt.args.delivery.HasLoader).
				WillReturnResult(tt.result)

			mock.ExpectCommit()

			err := repo.CreateDelivery(context.Background(), tt.args.delivery)
			require.Nil(t, deep.Equal(tt.error, err))

		})
	}
}

func TestDeliveryRepo_CreateDeliveryRollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	testLogger := logrus.New()
	repo := NewDeliveryRepo(db, logger.New(testLogger))

	type args struct {
		ctx      context.Context
		delivery *entity.Delivery
	}
	tests := []struct {
		name   string
		args   args
		rows   *sqlmock.Rows
		result driver.Result
		error  error
	}{
		{
			name: "no return id value in first query",
			args: args{
				ctx: context.Background(),
				delivery: &entity.Delivery{
					ID:       1,
					ClientID: 1,
					StatusID: 1,
					TypeID:   1,
					Geo: &entity.Geo{
						FromLongitude: 37.22,
						FromLatitude:  55.77,
						FromObject:    "улица веселая д.1",
						ToLongitude:   37.22,
						ToLatitude:    55.77,
						ToObject:      "улица веселая д.10",
						Distance:      1200,
					},
					Price:     1220,
					HasLoader: true,
				},
			},
			rows:  sqlmock.NewRows([]string{"id"}).AddRow(nil),
			error: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()

			mock.ExpectQuery(regexp.QuoteMeta(`
				INSERT INTO geo(from_longitude, from_latitude, from_object, to_longitude, to_latitude, to_object, distance)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
			`)).
				WithArgs(
					tt.args.delivery.Geo.FromLongitude,
					tt.args.delivery.Geo.FromLatitude,
					tt.args.delivery.Geo.FromObject,
					tt.args.delivery.Geo.ToLongitude,
					tt.args.delivery.Geo.ToLatitude,
					tt.args.delivery.Geo.ToObject,
					tt.args.delivery.Geo.Distance,
				).
				WillReturnRows(tt.rows).
				WillReturnError(tt.error)

			mock.ExpectRollback()

			err := repo.CreateDelivery(context.Background(), tt.args.delivery)
			require.Nil(t, deep.Equal(tt.error, err))

		})
	}
}
