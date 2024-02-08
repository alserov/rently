package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alserov/rently/carsharing/internal/db"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
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

func (r *repository) CheckIfCarAvailableInPeriod(_ context.Context, carUUID string, from, to time.Time) (bool, error) {
	query := `SELECT count(*) FROM rents WHERE car_uuid = $1 AND $2 BETWEEN rent_start AND rent_end OR $3 BETWEEN rent_start AND rent_end`

	var found int
	if err := r.db.QueryRowx(query, carUUID, from, to).Scan(&found); err != nil {
		return false, &models.Error{
			Msg:    fmt.Sprintf("failed to check if car is available: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return found == 0, nil
}

func (r *repository) GetRentsWhatStartsOnDate(ctx context.Context, tomorrowDate time.Time) ([]models.RentStartData, error) {
	query := `SELECT car_uuid, user_uuid, rent_start, rent_end FROM rents WHERE user_uuid != '' AND rent_start = $1`

	rows, err := r.db.Queryx(query, tomorrowDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, &models.Error{
			Msg:    fmt.Sprintf("failed to select rents which start on %v: %v", tomorrowDate, err),
			Status: http.StatusInternalServerError,
		}
	}

	var rentData []models.RentStartData
	for rows.Next() {
		var rent models.RentStartData
		if err = rows.StructScan(&rent); err != nil {
			return nil, &models.Error{
				Msg:    fmt.Sprintf("failed to scan into struct: %v", err),
				Status: http.StatusInternalServerError,
			}
		}
		rentData = append(rentData, rent)
	}

	return rentData, nil
}

func (r *repository) CreateCharge(ctx context.Context, req models.Charge) error {
	query := `INSERT INTO charges (rent_uuid,charge_uuid,charge_amount) VALUES ($1,$2,$3)`

	if err := r.db.QueryRowx(query, req.RentUUID, req.ChargeUUID, req.ChargeAmount).Err(); err != nil {
		return &models.Error{
			Msg:    fmt.Sprintf("failed to insert charge: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r *repository) CancelRentTx(ctx context.Context, rentUUID string) (models.CancelRentInfo, db.Tx, error) {
	query := `DELETE FROM rents WHERE rent_uuid = $1 LIMIT 1 RETURNING rent_price, charge_id`

	tx, err := r.db.Beginx()
	if err != nil {
		return models.CancelRentInfo{}, nil, &models.Error{
			Msg:    fmt.Sprintf("failed to start tx: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	row := r.db.QueryRowx(query, rentUUID)

	var rentInfo models.CancelRentInfo
	if err = row.Scan(&rentInfo); err != nil {
		return models.CancelRentInfo{}, tx, status.Error(codes.Internal, fmt.Sprintf("failed to cancel rent: %v", err))
	}

	return rentInfo, tx, nil
}

func (r *repository) CreateCar(ctx context.Context, car models.Car) error {
	query := `INSERT INTO cars (uuid, brand, type,max_speed,seats,category,price_per_day, image_uuid)
				VALUES ($1,$2,$3,$4,$5,$6,$7, $8)`

	_, err := r.db.Query(query, car.UUID, car.Brand, car.Type, car.MaxSpeed, car.Seats, car.Category, car.PricePerDay, car.MainImage)
	if err != nil {
		return &models.Error{
			Status: http.StatusInternalServerError,
			Msg:    fmt.Sprintf("failed to execute query: %v", err),
		}
	}

	query = `INSERT INTO images (uuid, car_uuid)
				VALUES ($1,$2)`

	_, err = r.db.Query(query, car.MainImage, car.UUID)
	if err != nil {
		return &models.Error{
			Status: http.StatusInternalServerError,
			Msg:    fmt.Sprintf("failed to execute query: %v", err),
		}
	}

	return nil
}

func (r *repository) DeleteCar(ctx context.Context, uuid string) error {
	query := `DELETE FROM cars WHERE uuid = $1`

	err := r.db.QueryRow(query, uuid).Err()
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Error{
			Msg:    fmt.Sprintf("car with uuid: %s not found", uuid),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return &models.Error{
			Msg:    fmt.Sprintf("failed to execute delete query: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r *repository) UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error {
	query := `UPDATE cars SET price_per_day = $1 WHERE uuid = $2`

	_, err := r.db.Exec(query, req.Price, req.CarUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Error{
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("carsharing with uuid: %s not found", req.CarUUID),
		}
	}
	if err != nil {
		return &models.Error{
			Status: http.StatusInternalServerError,
			Msg:    fmt.Sprintf("failed to execute delete query: %v", err),
		}
	}

	return nil
}

func (r *repository) CreateRentTx(_ context.Context, req models.CreateRentReq) (float32, db.Tx, error) {
	query := `INSERT INTO rents(rent_uuid,car_uuid, user_uuid,phone_number,passport_number,rent_start,rent_end)
				VALUES ($1,$2,$3,$4,$5,$6, $7) RETURNING (SELECT price_per_day FROM cars WHERE uuid = $2 LIMIT 1)`

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, tx, &models.Error{
			Msg:    fmt.Sprintf("failed to start tx: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	var carPricePerDay float32
	err = tx.QueryRowx(query, req.RentUUID, req.CarUUID, req.UserUUID, req.PhoneNumber, req.PassportNumber, req.RentStart, req.RentEnd).Scan(&carPricePerDay)
	if err != nil {
		return 0, tx, &models.Error{
			Status: http.StatusInternalServerError,
			Msg:    fmt.Sprintf("failed to create rent: %v", err),
		}
	}

	return carPricePerDay, tx, nil
}

func (r *repository) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.CarMainInfo, error) {
	query := `SELECT cars.uuid,brand, type,category, price_per_day, images.uuid as image FROM cars LEFT JOIN images ON cars.image_uuid = images.uuid
                WHERE 
                    (LOWER(brand) = LOWER($1) OR $1 = '') AND
                    (LOWER(type) = LOWER($2) OR $2 = '') AND
                	(max_speed > $3 OR $3 = 0) AND
                	(seats = $4 OR $4 = 0) AND
                	(LOWER(category) = LOWER($5) OR $5 = '') AND
                	(price_per_day < $6 OR $6 = 0)
                    `

	rows, err := r.db.Queryx(query, params.Brand, params.Type, params.MaxSpeed, params.Seats, params.Category, params.PricePerDay)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cars []models.CarMainInfo
	for rows.Next() {
		var car models.CarMainInfo
		if err = rows.StructScan(&car); err != nil {
			return nil, &models.Error{
				Msg:    fmt.Sprintf("failed to scan: %v", err),
				Status: http.StatusInternalServerError,
			}
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
		return models.Car{}, &models.Error{
			Msg:    fmt.Sprintf("car not found: %s", uuid),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return models.Car{}, &models.Error{
			Msg:    fmt.Sprintf("failed to get car: %s", uuid),
			Status: http.StatusNotFound,
		}
	}

	query = `SELECT uuid FROM images WHERE car_uuid = $1`

	err = r.db.Get(&car.Images, query, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Car{}, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return models.Car{}, status.Error(codes.Internal, err.Error())
	}

	return car, nil
}

func (r *repository) GetAvailableCars(ctx context.Context, period models.Period) ([]models.CarMainInfo, error) {
	query := `SELECT uuid,brand, type, max_speed, seats, category, price_per_day, images FROM cars 
            	WHERE uuid NOT IN (SELECT car_uuid FROM rents WHERE rent_start > $1 OR rent_end < $2)`

	rows, err := r.db.Queryx(query, period.Start, period.End)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, ERR_NO_ROWS)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cars []models.CarMainInfo
	for rows.Next() {
		var car models.CarMainInfo
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
		return models.Rent{}, &models.Error{
			Status: http.StatusNotFound,
			Msg:    fmt.Sprintf("rent with uuid: %s not found", rentUUID),
		}
	}
	if err != nil {
		return models.Rent{}, &models.Error{
			Status: http.StatusInternalServerError,
			Msg:    fmt.Sprintf("failed to check rent: %v", err),
		}
	}

	return rent, nil
}
