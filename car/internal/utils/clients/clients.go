package clients

import (
	"fmt"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"
	"google.golang.org/grpc"
)

type Services struct {
	FileStorageAddr string
}

type Clients struct {
	FileStorage fstorage.FileStorageClient
}

type CloseFunc func([]*grpc.ClientConn)

const clientsAmount = 1

func SetupClients(addr Services) (*Clients, []*grpc.ClientConn, CloseFunc) {
	clients := &Clients{}
	conns := make([]*grpc.ClientConn, 0, clientsAmount)

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
		err := c.Close()
		if err != nil {
			panic("failed to close connection: " + err.Error())
		}
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
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return cc, nil
}
