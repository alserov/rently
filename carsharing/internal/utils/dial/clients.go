package dial

import (
	"google.golang.org/grpc"
)

type Services struct {
	UserAddr string
}

type Connections struct {
	UserConn *grpc.ClientConn
}

