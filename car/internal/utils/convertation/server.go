package convertation

import (
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/proto/gen/car"
)

type ServerConverter interface {
	CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq

	CheckRentToPb(res models.Rent) *car.CheckRentRes
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s serverConverter) CheckRentToPb(res models.Rent) *car.CheckRentRes {
	//TODO implement me
	panic("implement me")
}

func (s serverConverter) CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq {
	//TODO implement me
	panic("implement me")
}
