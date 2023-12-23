package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *repository) CancelRent(ctx context.Context, rentUUID string) (err error) {
	r.db.Exec("")
	panic("")
}

const (
	ERR_NOROWS = "nothing found"
)

func (r *repository) CheckRent(_ context.Context, rentUUID string) (rent models.Rent, err error) {
	query := "SELECT * FROM rents WHERE rent_uuid = $1"

	err = r.db.Get(&rent, query, rentUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Rent{}, status.Error(codes.NotFound, ERR_NOROWS)
	}
	if err != nil {
		return models.Rent{}, status.Error(codes.Internal, err.Error())
	}

	return rent, nil
}

func (r *repository) CreateRent(_ context.Context, req models.CreateRentReq) (err error) {
	query := "INSERT INTO rents(rent_uuid, car_uuid, phone_number,passport_number,card_credentials, rent_start,rent_end)" +
		"VALUES ($1,$2,$3,$4,$5,$6,$7)"

	_, err = r.db.
		Exec(query, req.RentUUID, req.CarUUID, req.PhoneNumber, req.PassportNumber, req.CardCredentials, req.RentStart, req.RentEnd)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
