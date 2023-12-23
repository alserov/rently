package convertation

import (
	repo "github.com/alserov/rently/car/internal/db/models"
	"github.com/alserov/rently/car/internal/service/models"
)

type ServiceConverter interface {
	CreateRentToRepo(req models.CreateRentReq) repo.CreateRentReq
	CheckRentToService(res repo.Rent) models.Rent
}

func NewServiceConverter() ServiceConverter {
	return &serviceConverter{}
}

type serviceConverter struct {
}

func (s serviceConverter) CheckRentToService(res repo.Rent) models.Rent {
	return models.Rent{
		CarUUID:   res.CarUUID,
		RentPrice: res.RentPrice,
		RentStart: res.RentStart,
		RentEnd:   res.RentEnd,
	}
}

func (s serviceConverter) CreateRentToRepo(req models.CreateRentReq) repo.CreateRentReq {
	//TODO implement me
	panic("implement me")
}
