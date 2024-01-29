package converter

import (
	"github.com/alserov/rently/api/internal/models"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/proto/gen/user"
)

type ServerConverter interface {
	CreateCarReqToPb(req models.CreateCarReq) *carsharing.CreateCarReq
	DeleteCarReqToPb(uuid string) *carsharing.DeleteCarReq
	UpdateCarPriceToPb(req models.UpdateCarPriceReq) *carsharing.UpdateCarPriceReq
	GetCarByUUIDReqToPb(uuid string) *carsharing.GetCarByUUIDReq
	GetCarsParamsReqToPb(req models.GetCarsByParamsReq) *carsharing.GetCarsByParamsReq
	GetImage(bucket string, id string) *carsharing.GetImageReq
	CheckIfAuthorizedReqToPb(token string) *user.CheckIfAuthorizedReq
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s *serverConverter) CheckIfAuthorizedReqToPb(token string) *user.CheckIfAuthorizedReq {
	return &user.CheckIfAuthorizedReq{
		Token: token,
	}
}

func (s *serverConverter) GetImage(bucket string, id string) *carsharing.GetImageReq {
	return &carsharing.GetImageReq{
		Bucket: bucket,
		Id:     id,
	}
}

func (s *serverConverter) UpdateCarPriceToPb(req models.UpdateCarPriceReq) *carsharing.UpdateCarPriceReq {
	return &carsharing.UpdateCarPriceReq{
		CarUUID:     req.UUID,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) GetCarByUUIDReqToPb(uuid string) *carsharing.GetCarByUUIDReq {
	return &carsharing.GetCarByUUIDReq{
		UUID: uuid,
	}
}

func (s *serverConverter) GetCarsParamsReqToPb(req models.GetCarsByParamsReq) *carsharing.GetCarsByParamsReq {
	return &carsharing.GetCarsByParamsReq{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) GetCarByUUIDReq(uuid string) *carsharing.GetCarByUUIDReq {
	return &carsharing.GetCarByUUIDReq{
		UUID: uuid,
	}
}

func (s *serverConverter) UpdateCarPrice(req models.UpdateCarPriceReq) *carsharing.UpdateCarPriceReq {
	return &carsharing.UpdateCarPriceReq{
		CarUUID:     req.UUID,
		PricePerDay: req.PricePerDay,
	}
}

func (s *serverConverter) DeleteCarReqToPb(uuid string) *carsharing.DeleteCarReq {
	return &carsharing.DeleteCarReq{
		CarUUID: uuid,
	}
}

func (s *serverConverter) CreateCarReqToPb(req models.CreateCarReq) *carsharing.CreateCarReq {
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
