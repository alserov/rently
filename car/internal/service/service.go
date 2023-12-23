package service

import (
	"context"
	"github.com/alserov/rently/car/internal/db"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/google/uuid"
)

type Service interface {
	CreateRent(ctx context.Context, req models.CreateRentReq) (rentUUID string, err error)
	CancelRent(ctx context.Context, rentUUID string) (err error)
	CheckRent(ctx context.Context, rentUUID string) (res models.Rent, err error)
}

func NewService(repo db.Repository) Service {
	return &service{
		repo: repo,

		convert: convertation.NewServiceConverter(),
	}
}

type service struct {
	repo db.Repository

	convert convertation.ServiceConverter
}

func (s *service) CancelRent(ctx context.Context, rentUUID string) (err error) {
	if err := s.repo.CancelRent(ctx, rentUUID); err != nil {
		return err
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

	return req.RentUUID, nil
}
