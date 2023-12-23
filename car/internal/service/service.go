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

func (s *service) CancelRent(ctx context.Context, rentUUID string) (err error) {
	rentInfo, err = s.repo.CancelRent(ctx, rentUUID)
	if err != nil {
		return err
	}

	// TODO: money refund

	if err = s.metrics.DecreaseActiveRentsAmount(); err != nil {
		s.log.Error("failed to decrease active rents amount (metrics): ", err.Error())
	}

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

	if err = s.metrics.IncreaseActiveRentsAmount(); err != nil {
		s.log.Error("failed to increase active rents amount (metrics): ", err.Error())
	}

	return req.RentUUID, nil
}
