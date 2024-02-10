package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/carsharing/internal/cache"
	"github.com/alserov/rently/carsharing/internal/log"
	"github.com/alserov/rently/carsharing/internal/metrics"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/alserov/rently/carsharing/internal/service"
	"github.com/alserov/rently/carsharing/internal/utils/convertation"
	"github.com/alserov/rently/carsharing/internal/utils/validation"
	"net/http"
	"time"

	"github.com/alserov/rently/proto/gen/carsharing"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type Params struct {
	Service service.Service

	Metrics metrics.Metrics

	Cache cache.Cache
}

func RegisterGRPCServer(gRPCServer *grpc.Server, server Params) {
	carsharing.RegisterCarsServer(gRPCServer, newServer(server))
}

func newServer(p Params) carsharing.CarsServer {
	return &server{
		log:     log.GetLogger(),
		service: p.Service,
		cache:   p.Cache,
		metrics: p.Metrics,
		valid:   validation.NewValidator(),
		convert: convertation.NewServerConverter(),
	}
}

type server struct {
	carsharing.UnimplementedCarsServer

	log log.Logger

	cache cache.Cache

	metrics metrics.Metrics

	service service.Service

	valid validation.Validator

	convert convertation.ServerConverter
}

const (
	OP_CREATE_CAR         = "CREATE CAR"
	OP_GET_CARS_BY_PARAMS = "GET CARS BY PARAMS"
)

func (s *server) GetRentStartingOnDate(ctx context.Context, req *carsharing.GetRentStartingOnDateReq) (*carsharing.GetRentStartingOnDateRes, error) {
	ctx = s.ctxWithID(ctx)

	if err := s.valid.ValidateGetRentStartingTomorrowReq(req); err != nil {
		return nil, err
	}

	rents, err := s.service.GetRentsWhatStartsOnDate(ctx, req.StartingOn.AsTime())
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.RentsStartOnDateToPb(rents), nil
}

func (s *server) GetImage(ctx context.Context, req *carsharing.GetImageReq) (*carsharing.GetImageRes, error) {
	ctx = s.ctxWithID(ctx)

	if err := s.valid.ValidateGetCarImageReq(req); err != nil {
		return nil, err
	}

	image, err := s.service.GetImage(ctx, fmt.Sprintf("%s/%s", req.Bucket, req.Id))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.GetImageResToPb(image), nil
}

func (s *server) CreateCar(ctx context.Context, req *carsharing.CreateCarReq) (*emptypb.Empty, error) {
	ctx = s.ctxWithID(ctx)
	start := time.Now()

	if err := s.valid.ValidateCreateCarReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CreateCar(ctx, s.convert.CreateCarReqToService(req), req.Images, req.MainImage); err != nil {
		return nil, s.handleError(err)
	}

	s.metrics.ResponseTime(time.Since(start), OP_CREATE_CAR, http.MethodPost)

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteCar(ctx context.Context, req *carsharing.DeleteCarReq) (*emptypb.Empty, error) {
	ctx = s.ctxWithID(ctx)
	if err := s.valid.ValidateDeleteCarReq(req); err != nil {
		return nil, err
	}

	if err := s.service.DeleteCar(ctx, req.CarUUID); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) UpdateCarPrice(ctx context.Context, req *carsharing.UpdateCarPriceReq) (*emptypb.Empty, error) {
	ctx = s.ctxWithID(ctx)
	if err := s.valid.ValidateUpdateCarPriceReq(req); err != nil {
		return nil, err
	}

	if err := s.service.UpdateCarPrice(ctx, s.convert.UpdateCarPriceReqToService(req)); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateRent(ctx context.Context, req *carsharing.CreateRentReq) (*carsharing.CreateRentRes, error) {
	ctx = s.ctxWithID(ctx)

	s.log.Debug("received create rent request", slog.String("id", ctx.Value(models.ID).(string)))

	if err := s.valid.ValidateCreateRentReq(req); err != nil {
		return nil, err
	}

	res, err := s.service.CreateRent(ctx, s.convert.CreateRentReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CreateRentToPb(res), nil
}

func (s *server) CancelRent(ctx context.Context, req *carsharing.CancelRentReq) (*emptypb.Empty, error) {
	ctx = s.ctxWithID(ctx)
	if err := s.valid.ValidateCancelRentReq(req); err != nil {
		return nil, err
	}

	if err := s.service.CancelRent(ctx, req.RentUUID); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CheckRent(ctx context.Context, req *carsharing.CheckRentReq) (*carsharing.CheckRentRes, error) {
	ctx = s.ctxWithID(ctx)
	if err := s.valid.ValidateCheckRentReq(req); err != nil {
		return nil, err
	}

	rent, err := s.service.CheckRent(ctx, req.RentUUID)
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CheckRentToPb(rent), nil
}

func (s *server) GetAvailableCars(ctx context.Context, req *carsharing.GetAvailableCarsReq) (*carsharing.GetCarsRes, error) {
	ctx = s.ctxWithID(ctx)
	if err := s.valid.ValidateGetAvailableCarsReq(req); err != nil {
		return nil, err
	}

	cars, err := s.service.GetAvailableCars(ctx, s.convert.GetAvailableCarsReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CarsToPb(cars), nil
}

func (s *server) GetCarsByParams(ctx context.Context, req *carsharing.GetCarsByParamsReq) (*carsharing.GetCarsRes, error) {
	ctx = s.ctxWithID(ctx)
	start := time.Now()

	if err := s.valid.ValidateGetCarsByParamsReq(req); err != nil {
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

	s.metrics.ResponseTime(time.Since(start), OP_GET_CARS_BY_PARAMS, http.MethodGet)

	return s.convert.CarsToPb(cars), nil
}

func (s *server) GetCarByUUID(ctx context.Context, req *carsharing.GetCarByUUIDReq) (*carsharing.Car, error) {
	ctx = s.ctxWithID(ctx)
	if err := s.valid.ValidateGetCarByUUID(req); err != nil {
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
