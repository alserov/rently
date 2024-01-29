package domains

import (
	"bytes"
	"context"
	grpcbreaker "github.com/alserov/circuit_breaker/grpc"
	"github.com/alserov/rently/api/internal/log"
	"github.com/alserov/rently/api/internal/middleware"
	"github.com/alserov/rently/api/internal/models"
	"github.com/alserov/rently/api/internal/utils/converter"
	carsh "github.com/alserov/rently/proto/gen/carsharing"
	usr "github.com/alserov/rently/proto/gen/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"image"
	"image/jpeg"
	"log/slog"
	"net/http"
	"time"
)

const (
	maxServiceErrors = 10
	breakerTimeout   = time.Second * 5
)

func NewCarsharing(p Params[carsh.CarsClient]) Carsharing {
	return &carsharing{
		log:              log.GetLogger(),
		carsharingClient: p.Client,
		readTimeout:      p.ReadTimeout,
		writeTimeout:     p.WriteTimeout,
		valid:            validator.New(),
		breaker:          grpcbreaker.NewBreaker(maxServiceErrors, breakerTimeout),
		convert:          converter.NewServerConverter(),
	}
}

type Carsharing interface {
	CreateCar(c *fiber.Ctx) error
	DeleteCar(c *fiber.Ctx) error
	UpdateCarPrice(c *fiber.Ctx) error

	GetAvailableCars(c *fiber.Ctx) error
	GetCarsByParams(c *fiber.Ctx) error
	GetCarByUUID(c *fiber.Ctx) error
	GetImage(c *fiber.Ctx) error
}

type carsharing struct {
	log log.Logger

	readTimeout  time.Duration
	writeTimeout time.Duration

	valid *validator.Validate

	convert converter.ServerConverter

	carsharingClient carsh.CarsClient
	userClient       usr.UserClient

	breaker *grpcbreaker.Breaker
}

func (csh *carsharing) GetImage(c *fiber.Ctx) error {
	bucket := c.Params("bucket")
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	res, err := grpcbreaker.Execute(ctx, csh.carsharingClient.GetImage, csh.convert.GetImage(bucket, id), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader((*res).File))
	if err != nil {
		csh.log.Error("failed to decode bytes to img", slog.String("error", err.Error()))
		return err
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		csh.log.Error("failed to encode img to bytes", slog.String("error", err.Error()))
		return err
	}

	c.Status(http.StatusOK)
	c.Response().Header.Set("Content-Type", "image/jpeg")
	handleResponseError(c.Send(buffer.Bytes()))
	return nil
}

func (csh *carsharing) GetAvailableCars(c *fiber.Ctx) error {
	var req models.GetAvailableCarsReq
	if err := parseQueryParams(c, &req); err != nil {
		c.Status(http.StatusBadRequest)
		handleResponseError(c.Send(marshal(models.Error{Err: err.Error()})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	cars, err := grpcbreaker.Execute(ctx, csh.carsharingClient.GetAvailableCars, nil, csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	c.Status(http.StatusOK)
	handleResponseError(c.Send(marshal(cars)))
	return nil
}

func (csh *carsharing) GetCarsByParams(c *fiber.Ctx) error {
	var req models.GetCarsByParamsReq
	if err := parseQueryParams(c, &req); err != nil {
		c.Status(http.StatusBadRequest)
		handleResponseError(c.Send(marshal(models.Error{Err: err.Error()})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	cars, err := grpcbreaker.Execute(ctx, csh.carsharingClient.GetCarsByParams, csh.convert.GetCarsParamsReqToPb(req), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	res := csh.transformImageLinks((*cars).Cars)

	c.Status(http.StatusOK)
	handleResponseError(c.Send(marshal(res)))
	return nil
}

func (csh *carsharing) GetCarByUUID(c *fiber.Ctx) error {
	carUUID := c.Params("car_uuid")
	if len(carUUID) < 10 {
		c.Status(http.StatusBadRequest)
		handleResponseError(c.Send(marshal(models.Error{Err: "carsharing uuid is not valid"})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	cars, err := grpcbreaker.Execute(ctx, csh.carsharingClient.GetCarByUUID, csh.convert.GetCarByUUIDReqToPb(carUUID), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	c.Status(http.StatusOK)
	handleResponseError(c.Send(marshal(cars)))
	return nil
}

func (csh *carsharing) DeleteCar(c *fiber.Ctx) error {
	carUUID := c.Params("car_uuid")
	if len(carUUID) < 10 {
		c.Status(http.StatusBadRequest)
		handleResponseError(c.Send(marshal(models.Error{Err: "carsharing uuid is not valid"})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	_, err := grpcbreaker.Execute(ctx, csh.carsharingClient.DeleteCar, csh.convert.DeleteCarReqToPb(carUUID), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	c.Status(http.StatusOK)
	return nil
}

func (csh *carsharing) UpdateCarPrice(c *fiber.Ctx) error {
	var req models.UpdateCarPriceReq
	if err := decode(c.Request().Body(), &req, csh.valid); err != nil {
		c.Status(http.StatusBadRequest)
		handleResponseError(c.Send(marshal(models.Error{Err: err.Error()})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	_, err := grpcbreaker.Execute(ctx, csh.carsharingClient.UpdateCarPrice, csh.convert.UpdateCarPriceToPb(req), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)

		return nil
	}

	c.Status(http.StatusOK)
	return nil
}

func (csh *carsharing) CreateCar(c *fiber.Ctx) error {
	var req models.CreateCarReq
	if err := parseForm(c, &req); err != nil {
		c.Status(http.StatusBadRequest)
		handleResponseError(c.Send(marshal(models.Error{Err: err.Error()})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(csh.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	if err := csh.checkIfAuthorized(ctx, c); err != nil {
		handleResponseError(c.Status(http.StatusMethodNotAllowed).Send(marshal(err)))
		return nil
	}

	_, err := grpcbreaker.Execute(ctx, csh.carsharingClient.CreateCar, csh.convert.CreateCarReqToPb(req), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	c.Status(http.StatusCreated)
	return nil
}

const path = "http://localhost:3001/info/carsharing/car/image"

func (csh *carsharing) transformImageLinks(cars []*carsh.CarMainInfo) []*carsh.CarMainInfo {
	c := make([]*carsh.CarMainInfo, 0, len(cars))
	for _, car := range cars {
		car.Image = transformImageInfoToLink(car.UUID, car.Image)
		c = append(c, car)
	}
	return c
}

func (csh *carsharing) checkIfAuthorized(ctx context.Context, c *fiber.Ctx) error {
	token := c.Context().Value(middleware.AUTH_TOKEN).(string)

	res, err := grpcbreaker.Execute(ctx, csh.userClient.CheckIfAuthorized, csh.convert.CheckIfAuthorizedReqToPb(token), csh.breaker)
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	if !(*res).IsAuthorized {
		return &models.Error{
			Err: middleware.ERR_NOT_AUTHORIZED,
		}
	}

	if (*res).Role != "admin" {
		return &models.Error{
			Err: middleware.ERR_NOT_ALLOWED,
		}
	}

	return nil
}
