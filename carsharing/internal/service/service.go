package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alserov/rently/carsharing/internal/clients"
	"github.com/alserov/rently/carsharing/internal/db"
	"github.com/alserov/rently/carsharing/internal/log"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/alserov/rently/carsharing/internal/notifications"
	"github.com/alserov/rently/carsharing/internal/payment"
	"github.com/alserov/rently/carsharing/internal/storage"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Service interface {
	RentActions
	CarActions
	AdminActions
}

type AdminActions interface {
	CreateCar(ctx context.Context, car models.Car, imageFiles [][]byte, mainImage []byte) error
	DeleteCar(ctx context.Context, uuid string) error
	UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error
}

type CarActions interface {
	GetCarByUUID(ctx context.Context, uuid string) (car models.Car, err error)
	GetCarsByParams(ctx context.Context, params models.CarParams) (cars []models.CarMainInfo, err error)
	GetAvailableCars(ctx context.Context, period models.Period) (cars []models.CarMainInfo, err error)
	GetImage(ctx context.Context, imageId string) ([]byte, error)
}

type RentActions interface {
	CreateRent(ctx context.Context, req models.CreateRentReq) (res models.CreateRentRes, err error)
	CancelRent(ctx context.Context, rentUUID string) error
	CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error)
	GetRentsWhatStartsOnDate(ctx context.Context, startingOn time.Time) ([]models.RentStartData, error)
}

type Params struct {
	Payment       payment.Payer
	ImageStorage  storage.ImageStorage
	Notifications notifications.Notifier
	UserClient    clients.UserClient
	Repo          db.Repository
}

func NewService(p Params) Service {
	return &service{
		log:          log.GetLogger(),
		repo:         p.Repo,
		imageStorage: p.ImageStorage,
		payment:      p.Payment,
		notification: p.Notifications,
		userClient:   p.UserClient,
	}
}

type service struct {
	log log.Logger

	repo db.Repository

	payment payment.Payer

	userClient clients.UserClient

	imageStorage storage.ImageStorage

	notification notifications.Notifier
}

func (s *service) GetRentsWhatStartsOnDate(ctx context.Context, startingOn time.Time) ([]models.RentStartData, error) {
	rents, err := s.repo.GetRentsWhatStartsOnDate(ctx, startingOn)
	if err != nil {
		return nil, err
	}

	return rents, nil
}

func (s *service) GetImage(ctx context.Context, imageId string) ([]byte, error) {
	file, err := s.imageStorage.Get(ctx, imageId)
	if err != nil {
		return nil, &models.Error{
			Msg:    fmt.Sprintf("image not found: %v", err),
			Status: http.StatusNotFound,
		}
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, &models.Error{
			Msg:    fmt.Sprintf("failed to read file: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return b, err
}

func (s *service) CreateCar(ctx context.Context, car models.Car, imageFiles [][]byte, mainImage []byte) error {
	car.UUID = uuid.New().String()

	var (
		mu    = sync.Mutex{}
		chErr = make(chan error, 3)
		wg    = sync.WaitGroup{}
	)

	wg.Add(len(imageFiles) + 1)

	go func() {
		defer wg.Done()

		id, err := s.imageStorage.Save(ctx, car.UUID, bytes.NewReader(mainImage))
		if err != nil {
			chErr <- err
		}

		car.MainImage = id
	}()

	for _, f := range imageFiles {
		go func(f []byte, wg *sync.WaitGroup) {
			defer wg.Done()

			id, err := s.imageStorage.Save(ctx, car.UUID, bytes.NewReader(f))
			if err != nil {
				chErr <- err
			}

			mu.Lock()
			defer mu.Unlock()
			car.Images = append(car.Images, id)
		}(f, &wg)
	}

	go func() {
		wg.Wait()
		close(chErr)
	}()

	var errCounter int
	for err := range chErr {
		errCounter++
		if errCounter >= len(imageFiles)/3 {
			return fmt.Errorf("failed to upload %d images from %d: car creation canceled", len(imageFiles)/3, len(imageFiles))
		}
		s.log.Error("failed to upload image to storage", slog.String("error", err.Error()))
	}

	if err := s.repo.CreateCar(ctx, car); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}

func (s *service) DeleteCar(ctx context.Context, uuid string) error {
	if err := s.repo.DeleteCar(ctx, uuid); err != nil {
		return err
	}

	if err := s.imageStorage.Delete(ctx, uuid); err != nil {
		s.log.Error("failed to delete image from storage", slog.String("error", err.Error()))
	}

	return nil
}

func (s *service) UpdateCarPrice(ctx context.Context, req models.UpdateCarPriceReq) error {
	if err := s.repo.UpdateCarPrice(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *service) GetCarByUUID(ctx context.Context, uuid string) (models.Car, error) {
	car, err := s.repo.GetCarByUUID(ctx, uuid)
	if err != nil {
		return models.Car{}, err
	}

	return car, nil
}

func (s *service) GetCarsByParams(ctx context.Context, params models.CarParams) ([]models.CarMainInfo, error) {
	cars, err := s.repo.GetCarsByParams(ctx, params)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

func (s *service) GetAvailableCars(ctx context.Context, period models.Period) ([]models.CarMainInfo, error) {
	cars, err := s.repo.GetAvailableCars(ctx, period)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

func (s *service) CancelRent(ctx context.Context, rentUUID string) error {
	rent, tx, err := s.repo.CancelRentTx(ctx, rentUUID)
	defer func() {
		if err = tx.Rollback(); err != nil {
			s.log.Err("failed to rollback tx", err, "op: cancel rent")
		}
	}()
	if err != nil {
		return err
	}

	if err = s.payment.Refund(rent.ChargeID, rent.RentPrice); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return &models.Error{
			Msg:    fmt.Sprintf("failed to commit tx: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *service) CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error) {
	rent, err := s.repo.CheckRent(ctx, rentUUID)
	if err != nil {
		return models.Rent{}, fmt.Errorf("repository error: %w", err)
	}

	return rent, nil
}

func (s *service) CreateRent(ctx context.Context, req models.CreateRentReq) (models.CreateRentRes, error) {
	req.RentUUID = uuid.New().String()

	available, err := s.repo.CheckIfCarAvailableInPeriod(ctx, req.CarUUID, req.RentStart, req.RentEnd)
	if err != nil {
		return models.CreateRentRes{}, err
	}
	if !available {
		return models.CreateRentRes{}, &models.Error{
			Msg:    "this car is already rented in this period",
			Status: http.StatusBadRequest,
		}
	}

	if req.Token != "" {
		info, err := s.userClient.GetPassportAndPhone(ctx, req.Token)
		if err != nil {
			return models.CreateRentRes{}, err
		}

		req.PhoneNumber = info.PhoneNumber
		req.PassportNumber = info.PassportNumber
	}

	carPricePerDay, tx, err := s.repo.CreateRentTx(ctx, req)
	defer func() {
		if err = tx.Rollback(); err != nil {
			s.log.Err("failed to commit tx", err, "op: create rent")
		}
	}()
	if err != nil {
		return models.CreateRentRes{}, fmt.Errorf("repository error: %w", err)
	}

	rentPrice := carPricePerDay * float32(req.RentEnd.Sub(req.RentStart).Hours()/24) * 100

	chargeID, err := s.payment.Debit(req.PaymentSource, rentPrice)
	if err != nil {
		return models.CreateRentRes{}, fmt.Errorf("payment error: %w", err)
	}

	if err = s.repo.CreateCharge(ctx, models.Charge{ChargeUUID: chargeID, RentUUID: req.RentUUID, ChargeAmount: rentPrice}); err != nil {
		if err = s.payment.Refund(chargeID, rentPrice); err != nil {
			return models.CreateRentRes{}, fmt.Errorf("payment error: %w", err)
		}
		return models.CreateRentRes{}, fmt.Errorf("repository error: %w", err)
	}

	if err = tx.Commit(); err != nil {
		if err = s.payment.Refund(chargeID, rentPrice); err != nil {
			return models.CreateRentRes{}, fmt.Errorf("payment error: %w", err)
		}
		return models.CreateRentRes{}, &models.Error{
			Msg:    fmt.Sprintf("failed to commit tx: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return models.CreateRentRes{
		RentUUID: req.RentUUID,
	}, nil
}
