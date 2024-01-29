package clients

import (
	"fmt"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/proto/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	CarAddr  string
	UserAddr string
}

type Clients struct {
	CarClient  carsharing.CarsClient
	UserClient user.UserClient

	Conns []*grpc.ClientConn
}

func (c *Clients) CloseConns() error {
	var e error
	for _, conn := range c.Conns {
		if err := conn.Close(); err != nil {
			e = fmt.Errorf("%v: %w", err, e)
		}
	}
	return nil
}

func DialServices(s Services) *Clients {
	var clients Clients

	clients.Conns = append(clients.Conns, []*grpc.ClientConn{dial(s.CarAddr), dial(s.UserAddr)}...)
	clients.CarClient = carsharing.NewCarsClient(clients.Conns[len(clients.Conns)-1])
	clients.UserClient = user.NewUserClient(clients.Conns[len(clients.Conns)-1])

	return &clients
}

func dial(addr string) *grpc.ClientConn {
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("failed to dial: " + addr)
	}
	return cc
}
