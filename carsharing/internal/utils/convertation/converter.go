package convertation

import (
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/alserov/rently/proto/gen/carsharing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ServerConverter interface {
	PbToModel
	ModelToPb
}

type PbToModel interface {
	CreateCarReqToService(req *carsharing.CreateCarReq) models.Car
	CreateRentReqToService(req *carsharing.CreateRentReq) models.CreateRentReq

	GetCarsByParamsReqToService(req *carsharing.GetCarsByParamsReq) models.CarParams
	GetAvailableCarsReqToService(req *carsharing.GetAvailableCarsReq) models.Period

	UpdateCarPriceReqToService(req *carsharing.UpdateCarPriceReq) models.UpdateCarPriceReq
}

type ModelToPb interface {
	CreateRentToPb(res models.CreateRentRes) *carsharing.CreateRentRes
	CheckRentToPb(res models.Rent) *carsharing.CheckRentRes
	RentsStartOnDateToPb(res []models.RentStartData) *carsharing.GetRentStartingOnDateRes
	CarsToPb(res []models.CarMainInfo) *carsharing.GetCarsRes
	CarToPb(res models.Car) *carsharing.Car
	GetImageResToPb(res []byte) *carsharing.GetImageRes
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s *serverConverter) RentsStartOnDateToPb(res []models.RentStartData) *carsharing.GetRentStartingOnDateRes {
	r := &carsharing.GetRentStartingOnDateRes{}
	for _, rent := range res {
		r.RentsInfo = append(r.RentsInfo, &carsharing.CheckRentRes{
			CarUUID:   rent.CarUUID,
			UserUUID:  rent.UserUUID,
			RentStart: s.timeToTimestampPb(rent.RentStart),
			RentEnd:   s.timeToTimestampPb(rent.RentEnd),
		})
	}

	return r
}

func (s *serverConverter) GetImageResToPb(res []byte) *carsharing.GetImageRes {
	return &carsharing.GetImageRes{
		File: res,
	}
}

func (s *serverConverter) UpdateCarPriceReqToService(req *carsharing.UpdateCarPriceReq) models.UpdateCarPriceReq {
	return models.UpdateCarPriceReq{
		CarUUID: req.CarUUID,
		Price:   req.PricePerDay,
	}
}

func (s *serverConverter) CreateCarReqToService(req *carsharing.CreateCarReq) models.Car {
	return models.Car{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) CreateRentToPb(res models.CreateRentRes) *carsharing.CreateRentRes {
	return &carsharing.CreateRentRes{
		RentUUID: res.RentUUID,
	}
}

func (s *serverConverter) GetCarsByParamsReqToService(req *carsharing.GetCarsByParamsReq) models.CarParams {
	return models.CarParams{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) CarToPb(res models.Car) *carsharing.Car {
	return &carsharing.Car{
		Brand:       res.Brand,
		Type:        res.Type,
		MaxSpeed:    res.MaxSpeed,
		Seats:       res.Seats,
		Category:    res.Category,
		PricePerDay: res.PricePerDay,
		UUID:        res.UUID,
		Images:      res.Images,
	}
}

func (s *serverConverter) CarsToPb(res []models.CarMainInfo) *carsharing.GetCarsRes {
	var cars carsharing.GetCarsRes

	for _, v := range res {
		c := &carsharing.CarMainInfo{
			Brand:       v.Brand,
			Type:        v.Type,
			Category:    v.Category,
			PricePerDay: v.PricePerDay,
			UUID:        v.UUID,
			Image:       v.Image,
		}
		cars.Cars = append(cars.Cars, c)
	}

	return &cars
}

func (s *serverConverter) GetAvailableCarsReqToService(req *carsharing.GetAvailableCarsReq) models.Period {
	from := req.Start.AsTime()
	to := req.End.AsTime()
	return models.Period{
		Start: &from,
		End:   &to,
	}
}

func (s *serverConverter) CheckRentToPb(res models.Rent) *carsharing.CheckRentRes {
	return &carsharing.CheckRentRes{
		CarUUID:   res.CarUUID,
		RentPrice: res.RentPrice,
		RentStart: s.timeToTimestampPb(res.RentStart),
		RentEnd:   s.timeToTimestampPb(res.RentEnd),
	}
}

const UUID_IF_NOT_AUTHORIZED = "not authorized"

func (s *serverConverter) CreateRentReqToService(req *carsharing.CreateRentReq) models.CreateRentReq {
	return models.CreateRentReq{
		CarUUID:        req.CarUUID,
		UserUUID:       UUID_IF_NOT_AUTHORIZED,
		PhoneNumber:    req.PhoneNumber,
		PassportNumber: req.PassportNumber,
		PaymentSource:  req.PaymentSource,
		Token:          req.Token,
		RentStart:      req.RentStart.AsTime(),
		RentEnd:        req.RentEnd.AsTime(),
	}
}

func (s *serverConverter) timeToTimestampPb(time time.Time) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}
