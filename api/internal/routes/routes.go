package routes

import (
	"github.com/alserov/rently/api/internal/middleware"
	"github.com/alserov/rently/api/internal/server"
	"github.com/gofiber/fiber/v2"
)

const (
	RENT = "/rent"
	INFO = "/info"
	AUTH = "/auth"
)

func Setup(c *fiber.App, s *server.Server) {
	rent := c.Group(RENT)
	rent.Post("carsharing/", middleware.CheckIfAuthorized, s.Carsharing.CreateCar)
	rent.Delete("carsharing/:car_uuid", middleware.CheckIfAuthorized, s.Carsharing.DeleteCar)
	rent.Patch("carsharing/", middleware.CheckIfAuthorized, s.Carsharing.UpdateCarPrice)

	info := c.Group(INFO)
	info.Get("carsharing/car/image/:bucket/:id", s.Carsharing.GetImage)
	info.Get("carsharing/car/:car_uuid", s.Carsharing.GetCarByUUID)
	info.Get("carsharing/filter", s.Carsharing.GetCarsByParams)

	auth := c.Group(AUTH)
	auth.Post("register/", s.User.Register)
	auth.Get("login/", s.User.Login)
}
