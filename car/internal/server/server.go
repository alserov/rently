package server

import (
	"context"
	"errors"

	"github.com/alserov/rently/car/internal/metrics"
	"github.com/alserov/rently/car/internal/models"
	"github.com/alserov/rently/car/internal/service"
	"github.com/alserov/rently/car/internal/utils/clients"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/alserov/rently/car/internal/utils/validation"

	"github.com/alserov/rently/proto/gen/car"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"net/http"
)

type Server struct {
	Metrics metrics.Metrics

	Service service.Service

	Clients *clients.Clients
}

func RegisterGRPCServer(gRPCServer *grpc.Server, server Server, log *slog.Logger) {
	car.RegisterCarsServer(gRPCServer, newServer(server, log))
}

func newServer(serv Server, log *slog.Logger) car.CarsServer {
	return &server{
		log:        log,
		service:    serv.Service,
		clients:    serv.Clients,
		metrics:    serv.Metrics,
		validation: validation.NewValidator(),
		convert:    convertation.NewServerConverter(),
	}
}

type server struct {
	car.UnimplementedCarsServer

	log *slog.Logger

	metrics metrics.Metrics

	service service.Service

	clients *clients.Clients

	validation validation.Validator

	convert convertation.ServerConverter
}

func (s *server) CreateCar(ctx context.Context, req *car.CreateCarReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateCreateCarReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CreateCar(ctx, s.convert.CreateCarReqToService(req)); err != nil {
		return nil, handleError(err, s.log)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteCar(ctx context.Context, req *car.DeleteCarReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateDeleteCarReq(req); err != nil {
		return nil, err
	}

	if err := s.service.DeleteCar(ctx, req.CarUUID); err != nil {
		return nil, handleError(err, s.log)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) UpdateCarPrice(ctx context.Context, req *car.UpdateCarPriceReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateUpdateCarPriceReq(req); err != nil {
		return nil, err
	}

	if err := s.service.UpdateCarPrice(ctx, s.convert.UpdateCarPriceReqToService(req)); err != nil {
		return nil, handleError(err, s.log)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateRent(ctx context.Context, req *car.CreateRentReq) (*car.CreateRentRes, error) {
	if err := s.validation.ValidateCreateRentReq(req); err != nil {
		return nil, err
	}

	res, err := s.service.CreateRent(ctx, s.convert.CreateRentReqToService(req))
	if err != nil {
		return nil, handleError(err, s.log)
	}

	s.metrics.IncreaseActiveRentsAmount()

	return s.convert.CreateRentToPb(res), nil
}

func (s *server) CancelRent(ctx context.Context, req *car.CancelRentReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateCancelRentReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CancelRent(ctx, req.RentUUID); err != nil {
		return nil, err
	}

	s.metrics.DecreaseActiveRentsAmount()

	return &emptypb.Empty{}, nil
}

func (s *server) CheckRent(ctx context.Context, req *car.CheckRentReq) (*car.CheckRentRes, error) {
	if err := s.validation.ValidateCheckRentReq(req); err != nil {
		return nil, err
	}

	rent, err := s.service.CheckRent(ctx, req.RentUUID)
	if err != nil {
		return nil, handleError(err, s.log)
	}

	return s.convert.CheckRentToPb(rent), nil
}

func (s *server) GetAvailableCars(ctx context.Context, req *car.GetAvailableCarsReq) (*car.GetCarsRes, error) {
	if err := s.validation.ValidateGetAvailableCarsReq(req); err != nil {
		return nil, err
	}

	cars, err := s.service.GetAvailableCars(ctx, s.convert.GetAvailableCarsReqToService(req))
	if err != nil {
		return nil, err
	}

	return s.convert.CarsToPb(cars), nil
}

func (s *server) GetCarsByParams(ctx context.Context, req *car.GetCarsByParamsReq) (*car.GetCarsRes, error) {
	if err := s.validation.ValidateGetCarsByParamsReq(req); err != nil {
		return nil, err
	}

	cars, err := s.service.GetCarsByParams(ctx, s.convert.GetCarsByParamsReqToService(req))
	if err != nil {
		return nil, err
	}

	s.metrics.NotifyBrandDemand(req.Brand)

	return s.convert.CarsToPb(cars), nil
}

func (s *server) GetCarByUUID(ctx context.Context, req *car.GetCarByUUIDReq) (*car.Car, error) {
	if err := s.validation.ValidateGetCarByUUID(req); err != nil {
		return nil, err
	}

	car, err := s.service.GetCarByUUID(ctx, req.UUID)
	if err != nil {
		return nil, err
	}

	return s.convert.CarToPb(car), nil
}

const internalError = "internal error"

func handleError(err error, log *slog.Logger) error {
	e := &models.Error{}
	ok := errors.As(err, &e)

	if ok {
		switch e.Code {
		case http.StatusInternalServerError:
			log.Error(e.Msg)
			return status.Error(codes.Internal, internalError)
		case http.StatusBadRequest:
			return status.Error(codes.InvalidArgument, e.Msg)
		case http.StatusNotFound:
			return status.Error(codes.NotFound, e.Msg)
		}
	}

	log.Error("unexpected error", slog.String("error", err.Error()))
	return e
}
