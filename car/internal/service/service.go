package service

import (
	"context"
	"github.com/alserov/rently/car/internal/db"
	"github.com/alserov/rently/car/internal/metrics"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/google/uuid"
	"log/slog"
)

type Service interface {
	RentActions
	CarActions
}

type CarActions interface {
	GetCarByUUID(ctx context.Context, uuid string) (car models.Car, err error)
	GetCarsByParams(ctx context.Context, params models.CarParams) (cars []models.Car, err error)
	GetAvailableCars(ctx context.Context, period models.Period) (cars []models.Car, err error)
}

type RentActions interface {
	CreateRent(ctx context.Context, req models.CreateRentReq) (rentUUID string, err error)
	CancelRent(ctx context.Context, rentUUID string) (err error)
	CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error)
}

func NewService(repo db.Repository, metrics metrics.Metrics, log *slog.Logger) Service {
	return &service{
		log:     log,
		repo:    repo,
		metrics: metrics,
		convert: convertation.NewServiceConverter(),
	}
}

type service struct {
	log *slog.Logger

	repo db.Repository

	metrics metrics.Metrics

	convert convertation.ServiceConverter
}

func (s *service) GetCarByUUID(ctx context.Context, uuid string) (models.Car, error) {
	car, err := s.repo.GetCarByUUID(ctx, uuid)
	if err != nil {
		return models.Car{}, err
	}

	return s.convert.CarToService(car), nil
}

func (s *service) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.Car, error) {
	cars, err := s.repo.GetCarsByParams(ctx, s.convert.ParamsToRepo(params))
	if err != nil {
		return nil, err
	}

	return s.convert.CarsToService(cars), nil
}

func (s *service) GetAvailableCars(ctx context.Context, period models.Period) (availableCars []models.Car, err error) {
	cars, err := s.repo.GetAvailableCars(ctx, s.convert.PeriodToRepo(period))
	if err != nil {
		return nil, err
	}

	return s.convert.CarsToService(cars), nil
}

func (s *service) CancelRent(ctx context.Context, rentUUID string) (err error) {
	_, err = s.repo.CancelRent(ctx, rentUUID)
	if err != nil {
		return err
	}

	// TODO: money refund

	s.metrics.DecreaseActiveRentsAmount()
	return nil
}

func (s *service) CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error) {
	rent, err := s.repo.CheckRent(ctx, rentUUID)
	if err != nil {
		return models.Rent{}, err
	}

	return s.convert.CheckRentToService(rent), nil
}

func (s *service) CreateRent(ctx context.Context, req models.CreateRentReq) (rentUUID string, err error) {
	req.RentUUID = uuid.New().String()

	if err = s.repo.CreateRent(ctx, s.convert.CreateRentToRepo(req)); err != nil {
		return "", err
	}

	s.metrics.IncreaseActiveRentsAmount()
	return req.RentUUID, nil
}
