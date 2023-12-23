package convertation

import (
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/proto/gen/car"
)

type ServerConverter interface {
	CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s serverConverter) CreateRentReqToService(req *car.CreateRentReq) models.CreateRentReq {
	//TODO implement me
	panic("implement me")
}
