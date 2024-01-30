package domains

import (
	"context"
	"github.com/alserov/rently/api/internal/middleware"
	"github.com/alserov/rently/api/internal/models"
	"github.com/alserov/rently/api/internal/utils/converter"
	usr "github.com/alserov/rently/proto/gen/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type User interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

func NewUser(p Params[usr.UserClient]) User {
	return &user{}
}

type user struct {
	userClient usr.UserClient

	valid *validator.Validate

	convert converter.Converter

	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (u *user) Register(c *fiber.Ctx) error {
	var req models.RegisterReq
	if err := decode(c.Request().Body(), &req, u.valid); err != nil {
		handleResponseError(c.Status(http.StatusBadRequest).Send(marshal(models.Error{Err: err.Error()})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(u.writeTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	res, err := u.userClient.Register(ctx, u.convert.RegisterReqToPb(req))
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	c.Cookie(&fiber.Cookie{
		Name:  middleware.AUTH_TOKEN,
		Value: res.Token,
	})

	handleResponseError(c.Send(marshal(models.RegisterRes{UUID: res.UUID})))
	return nil
}

func (u *user) Login(c *fiber.Ctx) error {
	var req models.LoginReq
	if err := decode(c.Request().Body(), &req, u.valid); err != nil {
		handleResponseError(c.Status(http.StatusBadRequest).Send(marshal(models.Error{Err: err.Error()})))
		return nil
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Duration(u.readTimeout.Seconds()*0.80*float64(time.Second)))
	defer cancel()

	res, err := u.userClient.Login(ctx, u.convert.LoginReqToPb(req))
	if err != nil {
		handleServiceError(c.Response(), err)
		return nil
	}

	c.Cookie(&fiber.Cookie{
		Name:  middleware.AUTH_TOKEN,
		Value: res.Token,
	})

	c.Status(http.StatusOK)
	return nil
}
