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

		validation: validation.NewServerValidator(),
		convert: convertation.NewServerConverter(),
	})
}

type server struct {
	service service.Service

	validation validation.ServerValidator
	convert    convertation.ServerConverter
}

func (s *server) CreateRent(ctx context.Context, req) () {
	if err := s.validation.ValidateCreateRentReq(req); err != nil {

	}

	rentUUID, err := s.CreateRent(ctx, s.convert.CreateOrderReqToService(req))
	if err != nil {

	}

	return , nil
}
