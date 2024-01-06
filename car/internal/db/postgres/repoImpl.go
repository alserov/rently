package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alserov/rently/car/internal/db"
	"github.com/alserov/rently/car/internal/db/models"
	global "github.com/alserov/rently/car/internal/models"
	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func NewRepo(db *sqlx.DB) db.Repository {
	return &repository{
		db: db,
	}
}

const (
	ERR_NO_ROWS = "nothing found"
)

type repository struct {
	db *sqlx.DB
}

func (r *repository) CancelRent(ctx context.Context, rentUUID string) (models.CancelRentInfo, error) {
	query := `DELETE FROM rents WHERE rent_uuid = $1 LIMIT 1 RETURNING rent_price, charge_id`

	row := r.db.QueryRowx(query, rentUUID)
	var rentInfo models.CancelRentInfo
	if err := row.Scan(&rentInfo); err != nil {
		return models.CancelRentInfo{}, status.Error(codes.Internal, fmt.Sprintf("failed to cancel rent: %v", err))
	}

	return rentInfo, nil
}

func (r *repository) CreateCar(ctx context.Context, car models.Car) error {
	query := `INSERT INTO cars (brand, type,max_speed,seats,category,price_per_day,uuid)
				VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := r.db.Exec(query, car.Brand, car.Type, car.MaxSpeed, car.Seats, car.Category, car.PricePerDay, car.UUID)
	if err != nil {
		return &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to execute query: %v", err),
		}
	}

	return nil
}

func (r *repository) DeleteCar(ctx context.Context, uuid string) error {
	query := `DELETE FROM cars WHERE uuid = $1`

	_, err := r.db.Exec(query, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to execute delete query: %v", err))
	}

	return nil
}

func (r *repository) UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error {
	query := `UPDATE cars SET price_per_day = $1 WHERE uuid = $2`

	_, err := r.db.Exec(query, req.Price, req.CarUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return &global.Error{
			Code: http.StatusNotFound,
			Msg:  fmt.Sprintf("car with uuid: %s not found", req.CarUUID),
		}
	}
	if err != nil {
		return &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to execute delete query: %v", err),
		}
	}

	return nil
}

func (r *repository) PrepareCreateRent(ctx context.Context, req models.CheckIfCarAvailable) (float32, error) {
	query := `SELECT price_per_day FROM cars WHERE uuid = $1`

	var pricePerDay float32
	err := r.db.Get(&pricePerDay, query, req.CarUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, &global.Error{
			Code: http.StatusNotFound,
			Msg:  fmt.Sprintf("car with uuid: %s not found", req.CarUUID),
		}
	}

	query = `SELECT 1 
				FROM rents 
				WHERE car_uuid = $3 
				AND $1 BETWEEN rent_start AND rent_end
				OR $2 BETWEEN rent_start AND rent_end
				LIMIT 1;`

	var engaged bool
	err = r.db.Get(&engaged, query, req.RentStart, req.RentEnd, req.CarUUID)
	if engaged {
		return 0, &global.Error{
			Code: http.StatusNotFound,
			Msg:  fmt.Sprintf("this car is not available from %v to %v", req.RentStart, req.RentEnd),
		}
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to get available cars: %v", err),
		}
	}

	return pricePerDay, nil
}

func (r *repository) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.Car, error) {
	query := `SELECT brand, type, max_speed, seats, category, price_per_day, uuid FROM cars as c WHERE 
                  lower(c.brand) LIKE $1  AND
			      lower(c.type) LIKE $2	  AND
			      c.max_speed > $3 		  AND
				  c.seats  > $4 		  AND
				  lower(c.category) LIKE $5 AND
				  c.price_per_day  < $6`

	rows, err := r.db.Queryx(query, params.Brand, params.Type, params.MaxSpeed, params.Seats, params.Category, params.PricePerDay)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err = rows.StructScan(&car); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *repository) GetCarByUUID(ctx context.Context, uuid string) (models.Car, error) {
	query := `SELECT * FROM cars WHERE uuid = $1`

	var car models.Car
	err := r.db.Get(&car, query, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Car{}, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return models.Car{}, status.Error(codes.Internal, err.Error())
	}

	return car, nil
}

func (r *repository) GetAvailableCars(ctx context.Context, period models.Period) ([]models.Car, error) {
	query := `SELECT brand, type, max_speed, seats, category, price_per_day, uuid FROM cars 
            	WHERE uuid NOT IN (SELECT car_uuid FROM rents WHERE rent_start > $1 OR rent_end < $2)`

	rows, err := r.db.Queryx(query, period.Start, period.End)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err = rows.StructScan(&car); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *repository) CheckRent(_ context.Context, rentUUID string) (models.Rent, error) {
	query := `SELECT rent_price, car_uuid, rent_start, rent_end FROM rents WHERE rent_uuid = $1`

	var rent models.Rent
	err := r.db.Get(&rent, query, rentUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Rent{}, &global.Error{
			Code: http.StatusNotFound,
			Msg:  fmt.Sprintf("rent with uuid: %s not found", rentUUID),
		}
	}
	if err != nil {
		return models.Rent{}, &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to check rent: %v", err),
		}
	}

	return rent, nil
}

func (r *repository) CreateRent(_ context.Context, req models.CreateRentReq) (err error) {
	query := `INSERT INTO rents(rent_uuid,rent_price, car_uuid, phone_number,passport_number,charge_id, rent_start,rent_end)
				VALUES ($1,$2,$3,$4,$5,$6,$7, $8)`

	_, err = r.db.
		Exec(query, req.RentUUID, req.RentPrice, req.CarUUID, req.PhoneNumber, req.PassportNumber, req.ChargeID, req.RentStart, req.RentEnd)
	if err != nil {
		return &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to create rent: %v", err),
		}
	}

	return nil
}
