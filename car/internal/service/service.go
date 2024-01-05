package service

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"

	"github.com/alserov/rently/car/internal/db"
	"github.com/alserov/rently/car/internal/metrics"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/car/internal/utils/broker"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/alserov/rently/car/internal/utils/payment"

	"github.com/google/uuid"
	"log/slog"
	"sync"
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

func NewService(repo db.Repository, metrics metrics.Metrics, log *slog.Logger) Service {
	consumerConfig := sarama.NewConfig()

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Errors = true
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{}, producerConfig)
	if err != nil {
		panic("failed to init async producer: " + err.Error())
	}

	return &service{
		log:            log,
		repo:           repo,
		metrics:        metrics,
		convert:        convertation.NewServiceConverter(),
		payment:        payment.NewPayer("sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh"),
		producer:       producer,
		consumerConfig: consumerConfig,
	}
}

type service struct {
	log *slog.Logger

	repo db.Repository

	metrics metrics.Metrics

	convert convertation.ServiceConverter

	payment payment.Payer

	consumerConfig *sarama.Config
	producer       sarama.AsyncProducer
	topics         broker.Topics
}

func (s *service) CreateCar(ctx context.Context, car models.Car[[]byte]) error {
	car.UUID = uuid.New().String()

	var (
		chErr = make(chan error)
		wg    = sync.WaitGroup{}
	)

	wg.Add(len(car.Images))

	for idx, img := range car.Images {
		go func(img []byte, idx int, wg *sync.WaitGroup) {
			defer wg.Done()
			b, err := json.Marshal(broker.SaveImageMessage{
				Value: img,
				UUID:  car.UUID,
				Idx:   idx,
			})
			if err != nil {
				chErr <- err
			}

			m := &sarama.ProducerMessage{
				Value: sarama.StringEncoder(b),
				Topic: s.topics.Images.Save,
			}

			s.producer.Input() <- m

			if err != nil {
				chErr <- err
			}
		}(img, idx, &wg)
	}

	go func() {
		wg.Wait()
		close(chErr)
	}()

	if err := s.repo.CreateCar(ctx, s.convert.CarToRepo(car)); err != nil {
		return err
	}

	if err := <-chErr; err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteCar(ctx context.Context, uuid string) error {
	chErr := make(chan error)
	go func() {
		defer close(chErr)

		m := &sarama.ProducerMessage{
			Value: sarama.StringEncoder(uuid),
			Topic: s.topics.Images.Delete,
		}

		s.producer.Input() <- m
	}()

	if err := s.repo.DeleteCar(ctx, uuid); err != nil {
		return err
	}

	if err := <-chErr; err != nil {
		return err
	}

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

	s.metrics.NotifyBrandDemand(params.Brand)

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

func (s *service) CreateRent(ctx context.Context, req models.CreateRentReq) (models.CreateRentRes, error) {
	req.RentUUID = uuid.New().String()

	if err := s.repo.CheckIfCarAvailable(ctx, s.convert.CheckIfCarAvailableToRepo(req)); err != nil {
		return models.CreateRentRes{}, err
	}

	totalPrice := s.payment.CountPrice(req.CarPricePerDay, &req)

	if err := s.repo.CreateRent(ctx, s.convert.CreateRentToRepo(req)); err != nil {
		return models.CreateRentRes{}, err
	}

	chargeID, err := s.payment.Debit(req.PaymentSource, totalPrice)
	if err != nil {
		return models.CreateRentRes{}, err
	}

	s.metrics.IncreaseActiveRentsAmount()
	return models.CreateRentRes{
		RentUUID: req.RentUUID,
		ChargeID: chargeID,
	}, nil
}
