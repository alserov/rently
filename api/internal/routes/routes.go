package routes

import (
	"github.com/alserov/rently/api/internal/middleware"
	"github.com/alserov/rently/api/internal/server"
	"github.com/gofiber/fiber/v2"
)

const (
	RENT  = "/rent"
	INFO  = "/info"
	AUTH  = "/auth"
	ADMIN = "/admin"
	USER  = "/user"
)

func Setup(c *fiber.App, s *server.Server) {
	admin := c.Group(ADMIN)
	admin.Post("carsharing/", middleware.CheckIfAuthorized, s.Carsharing.CreateCar)
	admin.Delete("carsharing/:car_uuid", middleware.CheckIfAuthorized, s.Carsharing.DeleteCar)
	admin.Patch("carsharing/", middleware.CheckIfAuthorized, s.Carsharing.UpdateCarPrice)

	info := c.Group(INFO)
	info.Get("carsharing/car/image/:bucket/:id", s.Carsharing.GetImage)
	info.Get("carsharing/car/:car_uuid", s.Carsharing.GetCarByUUID)
	info.Get("carsharing/filter", s.Carsharing.GetCarsByParams)

	auth := c.Group(AUTH)
	auth.Post("register/", s.User.Register)
	auth.Get("login/", s.User.Login)

	rent := c.Group(RENT)
	rent.Post("/carsharing/new", s.Carsharing.CreateRent)

	user := c.Group(USER)
	user.Patch("/reset-password", s.User.ResetPassword)
}
