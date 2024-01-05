package server

import (
	"context"
	"github.com/alserov/file-storage/internal/service"
	"github.com/alserov/file-storage/internal/utils/convertation"
	"github.com/alserov/file-storage/internal/utils/validation"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"
	"google.golang.org/grpc"
	"log/slog"
)

func RegisterGRPCServer(gRPCServer *grpc.Server, service service.Service) {
	fstorage.RegisterFileStorageServer(gRPCServer, &server{})
}

type server struct {
	fstorage.UnimplementedFileStorageServer

	log *slog.Logger

	service service.Service

	valid validation.Validator

	convert convertation.ServerConverter
}

func (s server) GetLinks(ctx context.Context, req *fstorage.GetLinksReq) (*fstorage.GetLinksRes, error) {
	if err := s.valid.ValidateGetLinksReq(req); err != nil {
		return nil, err
	}

	links, err := s.service.GetLinks(ctx, req.UUID)
	if err != nil {
		return nil, handleError(err, s.log)
	}

	return s.convert.LinksToPb(links), nil
}

func (s server) GetImage(ctx context.Context, req *fstorage.GetImageReq) (*fstorage.GetImageRes, error) {
	// validation

	imageBytes, err := s.service.GetImage(ctx, req.Link)
	if err != nil {
		return nil, handleError(err, s.log)
	}

	return s.convert.ImageBytesToPb(imageBytes), nil
}
