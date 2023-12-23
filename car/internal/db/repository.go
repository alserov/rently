package db

import (
	"context"
	"github.com/alserov/rently/car/internal/db/models"
)

type Repository interface {
	CreateRent(ctx context.Context, req models.CreateRentReq) (err error)
	CancelRent(ctx context.Context, rentUUID string) (rentInfo models.CancelRentInfo, err error)
	CheckRent(ctx context.Context, rentUUID string) (rent models.Rent, err error)
}
