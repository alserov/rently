package convertation

import (
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/proto/gen/car"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ServerConverter interface {
	PbToService
	ToPb
}

type PbToService interface {
	CreateCarReqToService(req *car.CreateCarReq) models.Car[[]byte]
	CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq

	GetCarsByParamsReqToService(req *car.GetCarsByParamsReq) models.CarParams
	GetAvailableCarsReqToService(req *car.GetAvailableCarsReq) models.Period

	UpdateCarPriceReqToService(req *car.UpdateCarPriceReq) models.UpdateCarPriceReq
}

type ToPb interface {
	GetLinksReqToPb(req models.Car[string]) *fstorage.GetLinksReq
	CreateRentToPb(res models.CreateRentRes) *car.CreateRentRes
	CheckRentToPb(res models.Rent) *car.CheckRentRes
	CarsToPb(res []models.Car[string]) *car.GetCarsRes
	CarToPb(res models.Car[string]) *car.Car
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s *serverConverter) GetLinksReqToPb(req models.Car[string]) *fstorage.GetLinksReq {
	return &fstorage.GetLinksReq{
		UUID: req.UUID,
	}
}

func (s *serverConverter) UpdateCarPriceReqToService(req *car.UpdateCarPriceReq) models.UpdateCarPriceReq {
	return models.UpdateCarPriceReq{
		CarUUID: req.CarUUID,
		Price:   req.PricePerDay,
	}
}

func (s *serverConverter) CreateCarReqToService(req *car.CreateCarReq) models.Car[[]byte] {
	return models.Car[[]byte]{
		Images:      req.Images,
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) CreateRentToPb(res models.CreateRentRes) *car.CreateRentRes {
	return &car.CreateRentRes{
		RentUUID: res.RentUUID,
		ChargeID: res.ChargeID,
	}
}

func (s *serverConverter) GetCarsByParamsReqToService(req *car.GetCarsByParamsReq) models.CarParams {
	return models.CarParams{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) CarToPb(res models.Car[string]) *car.Car {
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

func (s *serverConverter) CarsToPb(res []models.Car[string]) *car.GetCarsRes {
	var cars car.GetCarsRes

	for _, v := range res {
		c := &car.Car{
			Brand:       v.Brand,
			Type:        v.Type,
			MaxSpeed:    v.MaxSpeed,
			Seats:       v.Seats,
			Category:    v.Category,
			PricePerDay: v.PricePerDay,
			UUID:        v.UUID,
			Images:      v.Images,
		}
		cars.Cars = append(cars.Cars, c)
	}

	return &cars
}

func (s *serverConverter) GetAvailableCarsReqToService(req *car.GetAvailableCarsReq) models.Period {
	from := req.Start.AsTime()
	to := req.End.AsTime()
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
		PaymentSource:  req.PaymentSource,
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
