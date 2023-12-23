package convertation

import "github.com/alserov/rently/car/internal/service/models"

type ServerConverter interface {
	CreateOrderReqToService() models.CreateRentReq
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s serverConverter) CreateOrderReqToService() models.CreateRentReq {
	//TODO implement me
	panic("implement me")
}
