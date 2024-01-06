package db

import (
	"context"
	"github.com/alserov/rently/car/internal/db/models"
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
	GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.Car, error)
	GetCarByUUID(ctx context.Context, uuid string) (models.Car, error)
	GetAvailableCars(ctx context.Context, period models.Period) ([]models.Car, error)
}

type RentRepository interface {
	PrepareCreateRent(ctx context.Context, req models.CheckIfCarAvailable) (float32, error)

	CreateRent(ctx context.Context, req models.CreateRentReq) (err error)
	CancelRent(ctx context.Context, rentUUID string) (rentInfo models.CancelRentInfo, err error)
	CheckRent(ctx context.Context, rentUUID string) (rent models.Rent, err error)
}
