package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alserov/rently/car/internal/db"
	"github.com/alserov/rently/car/internal/db/models"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRepo(db *sqlx.DB) db.Repository {
	return &repository{db: db}
}

type repository struct {
	db *sqlx.DB
}

const (
	ERR_NO_ROWS = "nothing found"
)

func (r *repository) CancelRent(ctx context.Context, rentUUID string) (rentInfo models.CancelRentInfo, err error) {
	query := `DELETE FROM rents * WHERE rent_uuid = $1 RETURNING card_credentials, rent_price, car_uuid, rent_start,rent_end`

	err = r.db.QueryRow(query, rentUUID).Scan(&rentInfo)
	if errors.Is(err, sql.ErrNoRows) {
		return models.CancelRentInfo{}, status.Error(codes.NotFound, fmt.Sprintf("%s by rent uuid: %s", ERR_NO_ROWS, rentUUID))
	}
	if err != nil {
		return models.CancelRentInfo{}, status.Error(codes.Internal, err.Error())
	}

	return rentInfo, nil
}

func (r *repository) CheckRent(_ context.Context, rentUUID string) (rent models.Rent, err error) {
	query := "SELECT * FROM rents WHERE rent_uuid = $1"

	err = r.db.Get(&rent, query, rentUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Rent{}, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return models.Rent{}, status.Error(codes.Internal, err.Error())
	}

	return rent, nil
}

func (r *repository) CreateRent(_ context.Context, req models.CreateRentReq) (err error) {
	query := "INSERT INTO rents(rent_uuid,rent_price, car_uuid, phone_number,passport_number,card_credentials, rent_start,rent_end)" +
		"VALUES ($1,$2,$3,$4,$5,$6,$7, $8)"

	_, err = r.db.
		Exec(query, req.RentUUID, req.RentPrice, req.CarUUID, req.PhoneNumber, req.PassportNumber, req.CardCredentials, req.RentStart, req.RentEnd)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
