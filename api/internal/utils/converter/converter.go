package converter

import (
	"github.com/alserov/rently/api/internal/models"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/proto/gen/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Converter interface {
	CreateCarReqToPb(req models.CreateCarReq) *carsharing.CreateCarReq
	DeleteCarReqToPb(uuid string) *carsharing.DeleteCarReq
	UpdateCarPriceToPb(req models.UpdateCarPriceReq) *carsharing.UpdateCarPriceReq
	GetCarByUUIDReqToPb(uuid string) *carsharing.GetCarByUUIDReq
	GetCarsParamsReqToPb(req models.GetCarsByParamsReq) *carsharing.GetCarsByParamsReq
	GetImage(bucket string, id string) *carsharing.GetImageReq
	CheckIfAuthorizedReqToPb(token string) *user.CheckIfAuthorizedReq
	RegisterReqToPb(req models.RegisterReq) *user.RegisterReq
	LoginReqToPb(req models.LoginReq) *user.LoginReq
	CreateRentReqToPb(req models.CreateRentReq, token string) *carsharing.CreateRentReq
}

func NewConverter() Converter {
	return &converter{}
}

type converter struct {
}

func (s *converter) CreateRentReqToPb(req models.CreateRentReq, token string) *carsharing.CreateRentReq {
	return &carsharing.CreateRentReq{
		CarUUID:        req.CarUUID,
		PhoneNumber:    req.PhoneNumber,
		PassportNumber: req.PassportNumber,
		PaymentSource:  req.PaymentSource,
		Token:          token,
		RentStart: &timestamppb.Timestamp{
			Seconds: time.Unix(req.RentStart, 0).Unix(),
			Nanos:   int32(time.Unix(req.RentStart, 0).Nanosecond()),
		},
		RentEnd: &timestamppb.Timestamp{
			Seconds: time.Unix(req.RentEnd, 0).Unix(),
			Nanos:   int32(time.Unix(req.RentEnd, 0).Nanosecond()),
		},
	}
}

func (s *converter) RegisterReqToPb(req models.RegisterReq) *user.RegisterReq {
	return &user.RegisterReq{
		Username:       req.Username,
		Password:       req.Password,
		Email:          req.Email,
		PassportNumber: req.PassportNumber,
		PaymentSource:  req.PaymentSource,
		PhoneNumber:    req.PhoneNumber,
	}
}

func (s *converter) LoginReqToPb(req models.LoginReq) *user.LoginReq {
	return &user.LoginReq{
		Password: req.Password,
		Email:    req.Email,
	}
}

func (s *converter) CheckIfAuthorizedReqToPb(token string) *user.CheckIfAuthorizedReq {
	return &user.CheckIfAuthorizedReq{
		Token: token,
	}
}

func (s *converter) GetImage(bucket string, id string) *carsharing.GetImageReq {
	return &carsharing.GetImageReq{
		Bucket: bucket,
		Id:     id,
	}
}

func (s *converter) UpdateCarPriceToPb(req models.UpdateCarPriceReq) *carsharing.UpdateCarPriceReq {
	return &carsharing.UpdateCarPriceReq{
		CarUUID:     req.UUID,
		PricePerDay: req.PricePerDay,
	}
}

func (s *converter) GetCarByUUIDReqToPb(uuid string) *carsharing.GetCarByUUIDReq {
	return &carsharing.GetCarByUUIDReq{
		UUID: uuid,
	}
}

func (s *converter) GetCarsParamsReqToPb(req models.GetCarsByParamsReq) *carsharing.GetCarsByParamsReq {
	return &carsharing.GetCarsByParamsReq{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s *converter) GetCarByUUIDReq(uuid string) *carsharing.GetCarByUUIDReq {
	return &carsharing.GetCarByUUIDReq{
		UUID: uuid,
	}
}

func (s *converter) UpdateCarPrice(req models.UpdateCarPriceReq) *carsharing.UpdateCarPriceReq {
	return &carsharing.UpdateCarPriceReq{
		CarUUID:     req.UUID,
		PricePerDay: req.PricePerDay,
	}
}

func (s *converter) DeleteCarReqToPb(uuid string) *carsharing.DeleteCarReq {
	return &carsharing.DeleteCarReq{
		CarUUID: uuid,
	}
}

func (s *converter) CreateCarReqToPb(req models.CreateCarReq) *carsharing.CreateCarReq {
	return &carsharing.CreateCarReq{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
		MainImage:   req.MainImage,
		Images:      req.Images,
	}
}
