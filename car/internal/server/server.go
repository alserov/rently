package server

import (
	"context"
	"github.com/alserov/rently/car/internal/service"
	"github.com/alserov/rently/car/internal/utils/convertation"
	"github.com/alserov/rently/car/internal/utils/validation"
	"github.com/alserov/rently/proto/gen/car"
	"google.golang.org/grpc"
)

func RegisterGRPCServer(gRPCServer *grpc.Server, service service.Service) {
	car.RegisterCarsServer(gRPCServer, server{
		service: service,

		validation: validation.NewValidator(),
		convert: convertation.NewServerConverter(),
	})
}

type server struct {
	service service.Service

	validation validation.Validator
	convert    convertation.ServerConverter
}

func (s *server) CreateRent(ctx context.Context, req) () {
	if err := s.validation.ValidateCreateRentReq(req); err != nil {

	}

	if err = s.validation.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return "", err
	}

	if err = s.validation.ValidateCardCredentials(req.CardCredentials); err != nil {
		return "", err
	}

	if err = s.validation.ValidatePassportNumber(req.PassportNumber); err != nil {
		return "", err
	}

	rentUUID, err := s.CreateRent(ctx, s.convert.CreateOrderReqToService(req))
	if err != nil {

	}

	return , nil
}
