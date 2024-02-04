package db

import (
	"context"
	"github.com/alserov/rently/carsharing/internal/models"
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
	CreateRentTx(ctx context.Context, req models.CreateRentReq) (price float32, tx Tx, err error)
	CancelRentTx(ctx context.Context, rentUUID string) (rentInfo models.CancelRentInfo, tx Tx, err error)
	CheckRent(ctx context.Context, rentUUID string) (rent models.Rent, err error)
	CreateCharge(ctx context.Context, req models.Charge) error
	GetRentsWhatStartsOnDate(ctx context.Context, date time.Time) ([]models.RentStartData, error)
}

type Tx interface {
	Commit() error
	Rollback() error
}
