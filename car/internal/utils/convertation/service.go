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
	UpdateCarPriceToRepo(req models.UpdateCarPriceReq) repo.UpdateCarPriceReq
	CarToRepo(req models.Car[[]byte]) repo.Car
	ParamsToRepo(req models.CarParams) repo.CarParams
	PeriodToRepo(req models.Period) repo.Period
	CreateRentToRepo(req models.CreateRentReq, chargeID string, pricePerDay float32) repo.CreateRentReq
	CheckIfCarAvailableToRepo(req models.CreateRentReq) repo.CheckIfCarAvailable
}

type RepoToService interface {
	CarToCarWithImages(res repo.Car, links []string) models.Car[string]
	CheckRentToService(res repo.Rent) models.Rent
	CarsToService(res []repo.Car) []models.Car[string]
	CarToService(res repo.Car) models.Car[string]
}

func NewServiceConverter() ServiceConverter {
	return &serviceConverter{}
}

type serviceConverter struct {
}

func (s serviceConverter) CarToCarWithImages(res repo.Car, images []string) models.Car[string] {
	return models.Car[string]{
		Brand:       res.Brand,
		Type:        res.Type,
		MaxSpeed:    res.MaxSpeed,
		Seats:       res.Seats,
		Category:    res.Category,
		PricePerDay: res.PricePerDay,
		UUID:        res.UUID,
		Images:      images,
	}
}

func (s serviceConverter) UpdateCarPriceToRepo(req models.UpdateCarPriceReq) repo.UpdateCarPriceReq {
	return repo.UpdateCarPriceReq{
		Price:   req.Price,
		CarUUID: req.CarUUID,
	}
}

func (s serviceConverter) CarToRepo(req models.Car[[]byte]) repo.Car {
	return repo.Car{
		Brand:       req.Brand,
		Type:        req.Type,
		MaxSpeed:    req.MaxSpeed,
		Seats:       req.Seats,
		Category:    req.Category,
		PricePerDay: req.PricePerDay,
		UUID:        req.UUID,
	}
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

func (s serviceConverter) CarToService(res repo.Car) models.Car[string] {
	return models.Car[string]{
		Brand:       res.Brand,
		Type:        res.Type,
		MaxSpeed:    res.MaxSpeed,
		Seats:       res.Seats,
		Category:    res.Category,
		PricePerDay: res.PricePerDay,
		UUID:        res.UUID,
	}
}

func (s serviceConverter) CarsToService(res []repo.Car) []models.Car[string] {
	var cars []models.Car[string]

	for _, c := range res {
		car := models.Car[string]{
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

func (s serviceConverter) CreateRentToRepo(req models.CreateRentReq, chargeID string, pricePerDay float32) repo.CreateRentReq {
	return repo.CreateRentReq{
		RentUUID:       req.RentUUID,
		CarUUID:        req.CarUUID,
		PhoneNumber:    req.PhoneNumber,
		PassportNumber: req.PassportNumber,
		ChargeID:       chargeID,
		RentPrice:      pricePerDay,
		RentStart:      req.RentStart,
		RentEnd:        req.RentEnd,
	}
}
