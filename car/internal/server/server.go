package server

import (
	"context"
	"github.com/alserov/rently/car/internal/service"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/alserov/rently/car/internal/utils/validation"
	"github.com/alserov/rently/proto/gen/car"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func RegisterGRPCServer(gRPCServer *grpc.Server, service service.Service) {
	car.RegisterCarsServer(gRPCServer, &server{
		service:    service,
		validation: validation.NewValidator(),
		convert:    convertation.NewServerConverter(),
	})
}

type server struct {
	car.UnimplementedCarsServer

	service service.Service

	validation validation.Validator

	convert convertation.ServerConverter
}

func (s *server) CreateRent(ctx context.Context, req *car.CreateRentReq) (*car.CreateRentRes, error) {
	if err := s.validation.ValidateCreateRentReq(req); err != nil {
		return nil, err
	}

	res, err := s.service.CreateRent(ctx, s.convert.CreateRentReqToService(req))
	if err != nil {
		return nil, err
	}

	return s.convert.CreateRentToPb(res), nil
}

func (s *server) CancelRent(ctx context.Context, req *car.CancelRentReq) (*emptypb.Empty, error) {
	if err := s.validation.ValidateCancelRentReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CancelRent(ctx, req.RentUUID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CheckRent(ctx context.Context, req *car.CheckRentReq) (*car.CheckRentRes, error) {
	if err := s.validation.ValidateCheckRentReq(req); err != nil {
		return nil, err
	}

	rent, err := s.service.CheckRent(ctx, req.RentUUID)
	if err != nil {

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
