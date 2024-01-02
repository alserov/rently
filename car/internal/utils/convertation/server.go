package convertation

import (
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/proto/gen/car"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type ServerConverter interface {
	PbToService
	ToPb
}

type PbToService interface {
	GetAvailableCarsReqToService(req *car.GetAvailableCarsReq) models.Period
	CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq
	GetCarsByParamsReqToService(req *car.GetCarsByParamsReq) models.CarParams
}

type ToPb interface {
	CreateRentToPb(res models.CreateRentRes) *car.CreateRentRes
	CheckRentToPb(res models.Rent) *car.CheckRentRes
	CarsToPb(res []models.Car) *car.GetCarsRes
	CarToPb(res models.Car) *car.Car
}

func NewServerConverter() ServerConverter {
	return &serverConverter{
		wg: sync.WaitGroup{},
	}
}

type serverConverter struct {
	wg sync.WaitGroup
}

func (s *serverConverter) CreateRentToPb(res models.CreateRentRes) *car.CreateRentRes {
	return &car.CreateRentRes{
		RentUUID: res.RentUUID,
		ChargeID: res.ChargeID,
	}
}

func (s *serverConverter) GetCarsByParamsReqToService(req *car.GetCarsByParamsReq) models.CarParams {
	//TODO implement me
	panic("implement me")
}

func (s *serverConverter) CarToPb(res models.Car) *car.Car {
	return &car.Car{
		Brand:       res.Brand,
		Type:        res.Type,
		MaxSpeed:    res.MaxSpeed,
		Seats:       res.Seats,
		Category:    res.Category,
		PricePerDay: res.PricePerDay,
		UUID:        res.UUID,
	}
}

func (s *serverConverter) CarsToPb(res []models.Car) *car.GetCarsRes {
	var cars *car.GetCarsRes

	for _, c := range res {
		car := &car.Car{
			Brand:       c.Brand,
			Type:        c.Type,
			MaxSpeed:    c.MaxSpeed,
			Seats:       c.Seats,
			Category:    c.Category,
			PricePerDay: c.PricePerDay,
			UUID:        c.UUID,
		}
		cars.Cars = append(cars.Cars, car)
	}

	return cars
}

func (s *serverConverter) GetAvailableCarsReqToService(req *car.GetAvailableCarsReq) models.Period {
	from := req.From.AsTime()
	to := req.To.AsTime()
	return models.Period{
		Start: &from,
		End:   &to,
	}
}

func (s *serverConverter) CheckRentToPb(res models.Rent) *car.CheckRentRes {
	return &car.CheckRentRes{
		CarUUID:   res.CarUUID,
		RentPrice: res.RentPrice,
		RentStart: s.timeToTimestampPb(res.RentStart),
		RentEnd:   s.timeToTimestampPb(res.RentEnd),
	}
}

func (s *serverConverter) CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq {
	start := req.RentStart.AsTime()
	end := req.RentEnd.AsTime()

	return models.CreateRentReq{
		CarUUID:        req.CarUUID,
		PhoneNumber:    req.PhoneNumber,
		PassportNumber: req.PassportNumber,
		RentStart:      &start,
		RentEnd:        &end,
	}
}

func (s *serverConverter) timeToTimestampPb(time *time.Time) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}
