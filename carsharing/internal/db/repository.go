package db

import (
	"context"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository interface {
	RentRepository
	CarRepository
	AdminRepository
}

type AdminRepository interface {
	CreateCar(ctx context.Context, car models.Car) error
	DeleteCar(ctx context.Context, uuid string) error
	UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error
}

type CarRepository interface {
	GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.CarMainInfo, error)
	GetCarByUUID(ctx context.Context, uuid string) (models.Car, error)
	GetAvailableCars(ctx context.Context, period models.Period) ([]models.CarMainInfo, error)
}

type RentRepository interface {
	CheckIfCarAvailableInPeriod(ctx context.Context, carUUID string, from, to time.Time) (bool, error)
	CheckRent(ctx context.Context, rentUUID string) (rent models.Rent, err error)
	GetRentsWhatStartsOnDate(ctx context.Context, date time.Time) ([]models.RentStartData, error)
	Tx
}

type Tx interface {
	StartTx(ctx context.Context) (SqlTx, error)

	CreateRentTx(ctx context.Context, tx SqlTx, req models.CreateRentReq) (price float32, err error)
	CancelRentTx(ctx context.Context, tx SqlTx, rentUUID string) (rentInfo models.CancelRentInfo, err error)
	CreateChargeTx(ctx context.Context, tx SqlTx, req models.Charge) error
	RefundChargeTx(ctx context.Context, tx SqlTx, chargeUUID string) error
}

type SqlTx struct {
	*sqlx.Tx
}
