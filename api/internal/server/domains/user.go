package domains

import (
	usr "github.com/alserov/rently/proto/gen/user"
	"github.com/gofiber/fiber/v2"
)

type User interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

func NewUser(p Params[usr.UserClient]) User {
	return &user{}
}

type user struct {
	//userClient usr.UserClient
}

func (u user) Register(c *fiber.Ctx) error {

	//res, err := u.userClient.Register()
	//if err != nil {
	//	return nil
	//}
	//
	//c.Cookie(&fiber.Cookie{
	//	Name:  middleware.AUTH_TOKEN,
	//	Value: res.Token,
	//})
	panic("")
}

func (u user) Login(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
