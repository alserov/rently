package clients

import (
	"fmt"
	"github.com/alserov/rently/car/internal/server"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"
	"google.golang.org/grpc"
	"net/http"
)

type ClientAddresses struct {
	FileStorageAddr string
}

type Clients struct {
	FileStorage fstorage.FileStorageClient
}

type CloseFunc func([]*grpc.ClientConn)

const clientsAmount = 1

func SetupClients(addr ClientAddresses) (*Clients, []*grpc.ClientConn, CloseFunc) {
	clients := &Clients{}
	conns := make([]*grpc.ClientConn, clientsAmount)

	fStorageClient, cc, err := newFileStorageClient(addr.FileStorageAddr)
	if err != nil {
		panic(err)
	}
	clients.FileStorage = fStorageClient
	conns = append(conns, cc)

	return clients, conns, Close
}

func Close(conns []*grpc.ClientConn) {
	for _, c := range conns {
		c.Close()
	}
}

func newFileStorageClient(addr string) (fstorage.FileStorageClient, *grpc.ClientConn, error) {
	cc, err := dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial file storage: %w", err)
	}

	cl := fstorage.NewFileStorageClient(cc)

	return cl, cc, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	cc, err := grpc.Dial(addr)
	if err != nil {
		return nil, &server.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to dial: %v", err),
		}
	}

	return cc, nil
}
