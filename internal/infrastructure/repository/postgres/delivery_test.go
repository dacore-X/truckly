package postgres

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dacore-x/truckly/internal/entity"
	"github.com/stretchr/testify/require"
	"log"
	"regexp"
	"testing"
)

func TestDeliveryRepo_CreateDeliverySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	type args struct {
		ctx      context.Context
		delivery *entity.Delivery
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "usual case with type_id",
			args: args{
				ctx: context.Background(),
				delivery: &entity.Delivery{
					ClientID:      1,
					TypeID:        1,
					StatusID:      1,
					FromLatitude:  1,
					FromLongitude: 1,
					ToLatitude:    1,
					ToLongitude:   1,
					Distance:      100,
					Price:         100,
					HasLoader:     true,
				},
			},
			wantErr: false,
		},
		{
			name: "usual case without type_id",
			args: args{
				ctx: context.Background(),
				delivery: &entity.Delivery{
					ClientID:      1,
					StatusID:      1,
					FromLatitude:  1,
					FromLongitude: 1,
					ToLatitude:    1,
					ToLongitude:   1,
					Distance:      100,
					Price:         100,
					HasLoader:     true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(`
				INSERT INTO deliveries(client_id, status_id, type_id, from_longitude, to_longitude, from_latitude, to_latitude, distance, price, has_loader)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`)).
				WithArgs(
					tt.args.delivery.ClientID,
					tt.args.delivery.StatusID,
					tt.args.delivery.TypeID,
					tt.args.delivery.FromLongitude,
					tt.args.delivery.ToLongitude,
					tt.args.delivery.FromLatitude,
					tt.args.delivery.ToLatitude,
					tt.args.delivery.Distance,
					tt.args.delivery.Price,
					tt.args.delivery.HasLoader).
				WillReturnResult(sqlmock.NewResult(1, 1))

			dr := &DeliveryRepo{
				DB: db,
			}
			err := dr.CreateDelivery(tt.args.ctx, tt.args.delivery)
			log.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDelivery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeliveryRepo_CreateDeliveryFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	type args struct {
		ctx      context.Context
		delivery *entity.Delivery
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		error   error
	}{
		{
			name: "error with no client_id",
			args: args{
				ctx: context.Background(),
				delivery: &entity.Delivery{
					//ClientID:
					TypeID:        1,
					StatusID:      1,
					FromLatitude:  1,
					FromLongitude: 1,
					ToLatitude:    1,
					ToLongitude:   1,
					Distance:      100,
					Price:         100,
					HasLoader:     true,
				},
			},
			wantErr: true,
			error:   fmt.Errorf("SQL Error [23502] null value in column with relation with not null constraint"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(`
				INSERT INTO deliveries(client_id, status_id, type_id, from_longitude, to_longitude, from_latitude, to_latitude, distance, price, has_loader)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`)).
				WithArgs(
					tt.args.delivery.ClientID,
					tt.args.delivery.StatusID,
					tt.args.delivery.TypeID,
					tt.args.delivery.FromLongitude,
					tt.args.delivery.ToLongitude,
					tt.args.delivery.FromLatitude,
					tt.args.delivery.ToLatitude,
					tt.args.delivery.Distance,
					tt.args.delivery.Price,
					tt.args.delivery.HasLoader).
				WillReturnError(fmt.Errorf("SQL Error [23502] null value in column with relation with not null constraint"))

			dr := &DeliveryRepo{
				DB: db,
			}
			err := dr.CreateDelivery(tt.args.ctx, tt.args.delivery)
			log.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDelivery() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err.Error() != tt.error.Error() {
				t.Errorf("CreateDelivery() error = %v, wantError %v", err, tt.error)
			}
		})
	}
}
