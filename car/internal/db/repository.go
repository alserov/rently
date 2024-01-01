package db

import (
	"context"
	"github.com/alserov/rently/car/internal/db/models"
)

type Repository interface {
	RentRepository
	CarRepository
}

type CarRepository interface {
	GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.Car, error)
	GetCarByUUID(ctx context.Context, uuid string) (models.Car, error)
	GetAvailableCars(ctx context.Context, period models.Period) ([]models.Car, error)
}

type RentRepository interface {
	CheckIfCarAvailable(ctx context.Context, req models.CheckIfCarAvailable) error

	CreateRent(ctx context.Context, req models.CreateRentReq) (err error)
	CancelRent(ctx context.Context, rentUUID string) (rentInfo models.CancelRentInfo, err error)
	CheckRent(ctx context.Context, rentUUID string) (rent models.Rent, err error)
}
