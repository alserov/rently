package service

import (
	"context"
	"fmt"
	"github.com/alserov/rently/car/internal/db"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/alserov/rently/car/internal/utils/files"
	"github.com/alserov/rently/car/internal/utils/payment"

	"github.com/google/uuid"
	"log/slog"
)

type Service interface {
	RentActions
	CarActions
	AdminActions
}

type AdminActions interface {
	CreateCar(ctx context.Context, car models.Car[[]byte]) error
	DeleteCar(ctx context.Context, uuid string) error
	UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error
}

type CarActions interface {
	GetCarByUUID(ctx context.Context, uuid string) (car models.Car[string], err error)
	GetCarsByParams(ctx context.Context, params models.CarParams) (cars []models.Car[string], err error)
	GetAvailableCars(ctx context.Context, period models.Period) (cars []models.Car[string], err error)
}

type RentActions interface {
	CreateRent(ctx context.Context, req models.CreateRentReq) (res models.CreateRentRes, err error)
	CancelRent(ctx context.Context, rentUUID string) error
	CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error)
}

func NewService(repo db.Repository, imager files.Imager, log *slog.Logger) Service {
	return &service{
		log:     log,
		repo:    repo,
		images:  imager,
		convert: convertation.NewServiceConverter(),
		payment: payment.NewPayer("sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh"),
	}
}

type service struct {
	log *slog.Logger

	repo db.Repository

	convert convertation.ServiceConverter

	payment payment.Payer

	images files.Imager
}

func (s *service) SetLogger(log *slog.Logger) {
	s.log = log
}

func (s *service) CreateCar(ctx context.Context, car models.Car[[]byte]) error {
	car.UUID = uuid.New().String()

	if err := s.repo.CreateCar(ctx, s.convert.CarToRepo(car)); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	if err := s.images.Save(car.Images, car.UUID); err != nil {
		return fmt.Errorf("filer error: %w", err)
	}

	return nil
}

func (s *service) DeleteCar(ctx context.Context, uuid string) error {
	if err := s.repo.DeleteCar(ctx, uuid); err != nil {
		return err
	}

	s.images.Delete(uuid)

	return nil
}

func (s *service) UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error {
	if err := s.repo.UpdateCarPrice(ctx, s.convert.UpdateCarPriceToRepo(req)); err != nil {
		return err
	}

	return nil
}

func (s *service) GetCarByUUID(ctx context.Context, uuid string) (models.Car[string], error) {
	car, err := s.repo.GetCarByUUID(ctx, uuid)
	if err != nil {
		return models.Car[string]{}, err
	}

	return s.convert.CarToService(car), nil
}

func (s *service) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.Car[string], error) {
	cars, err := s.repo.GetCarsByParams(ctx, s.convert.ParamsToRepo(params))
	if err != nil {
		return nil, err
	}

	return s.convert.CarsToService(cars), nil
}

func (s *service) GetAvailableCars(ctx context.Context, period models.Period) ([]models.Car[string], error) {
	cars, err := s.repo.GetAvailableCars(ctx, s.convert.PeriodToRepo(period))
	if err != nil {
		return nil, err
	}

	return s.convert.CarsToService(cars), nil
}

func (s *service) CancelRent(ctx context.Context, rentUUID string) error {
	rent, err := s.repo.CancelRent(ctx, rentUUID)
	if err != nil {
		return err
	}

	if err = s.payment.Refund(rent.ChargeID, rent.RentPrice); err != nil {
		return err
	}

	return nil
}

func (s *service) CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error) {
	rent, err := s.repo.CheckRent(ctx, rentUUID)
	if err != nil {
		return models.Rent{}, fmt.Errorf("repository error: %w", err)
	}

	return s.convert.CheckRentToService(rent), nil
}

func (s *service) CreateRent(ctx context.Context, req models.CreateRentReq) (models.CreateRentRes, error) {
	req.RentUUID = uuid.New().String()

	var err error
	pricePerDay, err := s.repo.PrepareCreateRent(ctx, s.convert.CheckIfCarAvailableToRepo(req))
	if err != nil {
		return models.CreateRentRes{}, fmt.Errorf("repository error: %w", err)
	}

	rentPrice := s.payment.CountPrice(pricePerDay, &req)
	chargeID, err := s.payment.Debit(req.PaymentSource, rentPrice)
	if err != nil {
		return models.CreateRentRes{}, fmt.Errorf("payment error: %w", err)
	}

	if err = s.repo.CreateRent(ctx, s.convert.CreateRentToRepo(req, chargeID, rentPrice)); err != nil {
		return models.CreateRentRes{}, fmt.Errorf("repository error: %w", err)
	}

	return models.CreateRentRes{
		RentUUID: req.RentUUID,
		ChargeID: chargeID,
	}, nil
}
