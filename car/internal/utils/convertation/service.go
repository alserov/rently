package convertation

import (
	repo "github.com/alserov/rently/car/internal/db/models"
	"github.com/alserov/rently/car/internal/service/models"
)

type ServiceConverter interface {
	ToRepo
	RepoToService
}

type ToRepo interface {
	ParamsToRepo(req models.CarParams) repo.CarParams
	PeriodToRepo(req models.Period) repo.Period
	CreateRentToRepo(req models.CreateRentReq) repo.CreateRentReq
	CheckIfCarAvailableToRepo(req models.CreateRentReq) repo.CheckIfCarAvailable
}

type RepoToService interface {
	CheckRentToService(res repo.Rent) models.Rent
	CarsToService(res []repo.Car) []models.Car
	CarToService(res repo.Car) models.Car
}

func NewServiceConverter() ServiceConverter {
	return &serviceConverter{}
}

type serviceConverter struct {
}

func (s serviceConverter) CheckIfCarAvailableToRepo(req models.CreateRentReq) repo.CheckIfCarAvailable {
	return repo.CheckIfCarAvailable{
		CarUUID:   req.CarUUID,
		RentStart: req.RentStart,
		RentEnd:   req.RentEnd,
	}
}

func (s serviceConverter) ParamsToRepo(req models.CarParams) repo.CarParams {
	return repo.CarParams{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
	}
}

func (s serviceConverter) CarToService(res repo.Car) models.Car {
	return models.Car{
		Brand:       res.Brand,
		Type:        res.Type,
		MaxSpeed:    res.MaxSpeed,
		Seats:       res.Seats,
		Category:    res.Category,
		PricePerDay: res.PricePerDay,
		UUID:        res.UUID,
	}
}

func (s serviceConverter) CarsToService(res []repo.Car) []models.Car {
	var cars []models.Car

	for _, c := range res {
		car := models.Car{
			Brand:       c.Brand,
			Type:        c.Type,
			MaxSpeed:    c.MaxSpeed,
			Seats:       c.Seats,
			Category:    c.Category,
			PricePerDay: c.PricePerDay,
			UUID:        c.UUID,
		}
		cars = append(cars, car)
	}

	return cars
}

func (s serviceConverter) PeriodToRepo(req models.Period) repo.Period {
	return repo.Period{
		Start: req.Start,
		End:   req.End,
	}
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
	return repo.CreateRentReq{
		RentUUID:       req.RentUUID,
		CarUUID:        req.CarUUID,
		PhoneNumber:    req.PhoneNumber,
		PassportNumber: req.PassportNumber,
		ChargeID:       "",
		RentPrice:      100,
		RentStart:      req.RentStart,
		RentEnd:        req.RentEnd,
	}
}
