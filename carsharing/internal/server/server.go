package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/carsharing/internal/cache"
	"github.com/alserov/rently/carsharing/internal/log"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/alserov/rently/carsharing/internal/service"
	"github.com/alserov/rently/carsharing/internal/utils/convertation"
	"github.com/alserov/rently/carsharing/internal/utils/validation"
	"time"

	"github.com/alserov/rently/proto/gen/carsharing"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type Server struct {
	Service service.Service

	Cache cache.Cache
}

func RegisterGRPCServer(gRPCServer *grpc.Server, server Server) {
	carsharing.RegisterCarsServer(gRPCServer, newServer(server))
}

func newServer(serv Server) carsharing.CarsServer {
	return &server{
		log:        log.GetLogger(),
		service:    serv.Service,
		cache:      serv.Cache,
		validation: validation.NewValidator(),
		convert:    convertation.NewServerConverter(),
	}
}

type server struct {
	carsharing.UnimplementedCarsServer

	log log.Logger

	cache cache.Cache

	service service.Service

	validation validation.Validator

	convert convertation.ServerConverter
}

func (s *server) GetImage(ctx context.Context, req *carsharing.GetImageReq) (*carsharing.GetImageRes, error) {
	if err := s.validation.ValidateGetCarImageReq(req); err != nil {
		return nil, err
	}

	image, err := s.service.GetImage(ctx, fmt.Sprintf("%s/%s", req.Bucket, req.Id))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.GetImageResToPb(image), nil
}

func (s *server) CreateCar(ctx context.Context, req *carsharing.CreateCarReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateCreateCarReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CreateCar(ctx, s.convert.CreateCarReqToService(req), req.Images, req.MainImage); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteCar(ctx context.Context, req *carsharing.DeleteCarReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateDeleteCarReq(req); err != nil {
		return nil, err
	}

	if err := s.service.DeleteCar(ctx, req.CarUUID); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) UpdateCarPrice(ctx context.Context, req *carsharing.UpdateCarPriceReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateUpdateCarPriceReq(req); err != nil {
		return nil, err
	}

	if err := s.service.UpdateCarPrice(ctx, s.convert.UpdateCarPriceReqToService(req)); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateRent(ctx context.Context, req *carsharing.CreateRentReq) (*carsharing.CreateRentRes, error) {
	if err := s.validation.ValidateCreateRentReq(req); err != nil {
		return nil, err
	}

	fmt.Println("req")

	res, err := s.service.CreateRent(ctx, s.convert.CreateRentReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CreateRentToPb(res), nil
}

func (s *server) CancelRent(ctx context.Context, req *carsharing.CancelRentReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateCancelRentReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CancelRent(ctx, req.RentUUID); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CheckRent(ctx context.Context, req *carsharing.CheckRentReq) (*carsharing.CheckRentRes, error) {
	if err := s.validation.ValidateCheckRentReq(req); err != nil {
		return nil, err
	}

	rent, err := s.service.CheckRent(ctx, req.RentUUID)
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CheckRentToPb(rent), nil
}

func (s *server) GetAvailableCars(ctx context.Context, req *carsharing.GetAvailableCarsReq) (*carsharing.GetCarsRes, error) {
	if err := s.validation.ValidateGetAvailableCarsReq(req); err != nil {
		return nil, err
	}

	cars, err := s.service.GetAvailableCars(ctx, s.convert.GetAvailableCarsReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CarsToPb(cars), nil
}

func (s *server) GetCarsByParams(ctx context.Context, req *carsharing.GetCarsByParamsReq) (*carsharing.GetCarsRes, error) {
	if err := s.validation.ValidateGetCarsByParamsReq(req); err != nil {
		return nil, err
	}

	marshaledValue, err := json.Marshal(req)
	if err == nil {
		var res []models.CarMainInfo
		if err = s.cache.Get(ctx, string(marshaledValue), &res); err == nil {
			return s.convert.CarsToPb(res), nil
		}
	}

	cars, err := s.service.GetCarsByParams(ctx, s.convert.GetCarsByParamsReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	if err = s.cache.Set(ctx, string(marshaledValue), cars, time.Hour*12); err != nil {
		s.log.Error("failed to cache value", slog.String("error", err.Error()))
	}

	return s.convert.CarsToPb(cars), nil
}

func (s *server) GetCarByUUID(ctx context.Context, req *carsharing.GetCarByUUIDReq) (*carsharing.Car, error) {
	if err := s.validation.ValidateGetCarByUUID(req); err != nil {
		return nil, err
	}

	var car models.Car
	if err := s.cache.Get(ctx, req.UUID, car); err != nil {
		return s.convert.CarToPb(car), nil
	}

	car, err := s.service.GetCarByUUID(ctx, req.UUID)
	if err != nil {
		return nil, s.handleError(err)
	}

	if err = s.cache.Set(ctx, req.UUID, car, time.Hour*12); err != nil {
		s.log.Error("failed to set cache", slog.String("error", err.Error()))
	}

	return s.convert.CarToPb(car), nil
}
