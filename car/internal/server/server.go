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
	convert    convertation.ServerConverter
}

func (s *server) CreateRent(ctx context.Context, req *car.CreateRentReq) (*car.CreateRentRes, error) {
	if err := s.validation.ValidateCreateRentReq(req); err != nil {
		return nil, err
	}

	if err := s.validation.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return nil, err
	}

	if err := s.validation.ValidateCardCredentials(req.CardCredentials); err != nil {
		return nil, err
	}

	if err := s.validation.ValidatePassportNumber(req.PassportNumber); err != nil {
		return nil, err
	}

	rentUUID, err := s.service.CreateRent(ctx, s.convert.CreateRentReqToService(req))
	if err != nil {
		return nil, err
	}

	return &car.CreateRentRes{RentUUID: rentUUID}, nil
}

func (s *server) CancelRent(ctx context.Context, req *car.CancelRentReq) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) CheckRent(ctx context.Context, req *car.CheckRentReq) (*car.CheckRentRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetAvailableCars(ctx context.Context, req *car.GetAvailableCarsReq) (*car.GetCarsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetCarsByParams(ctx context.Context, req *car.GetCarsByParamsReq) (*car.GetCarsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s *server) GetCarByUUID(ctx context.Context, req *car.GetCarByUUIDReq) (*car.Car, error) {
	//TODO implement me
	panic("implement me")
}
