package clients

import (
	"fmt"
	"github.com/alserov/rently/proto/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	UserAddr string
}

type Clients struct {
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
	return e
}

func DialServices(s Services) *Clients {
	var clients Clients

	clients.Conns = append(clients.Conns, []*grpc.ClientConn{dial(s.UserAddr)}...)
	clients.UserClient = user.NewUserClient(clients.Conns[0])

	return &clients
}

func dial(addr string) *grpc.ClientConn {
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("failed to dial: " + addr)
	}
	return cc
}
