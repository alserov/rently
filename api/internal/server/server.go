package server

import (
	"github.com/alserov/rently/api/internal/server/domains"
	"github.com/alserov/rently/api/internal/utils/clients"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/proto/gen/user"
	"time"
)

type Server struct {
	Carsharing domains.Carsharing
	User       domains.User
}

type Params struct {
	Clients *clients.Clients

	ReadTimeout time.Duration

	WriteTimeout time.Duration
}

func NewServer(p Params) *Server {
	return &Server{
		Carsharing: domains.NewCarsharing(domains.Params[carsharing.CarsClient]{
			Client:       p.Clients.CarClient,
			ReadTimeout:  p.ReadTimeout,
			WriteTimeout: p.WriteTimeout,
		}),
		User: domains.NewUser(domains.Params[user.UserClient]{}),
	}
}
