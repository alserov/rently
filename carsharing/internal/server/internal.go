package server

import (
	"context"
	"errors"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

const internalError = "internal error"

func (s *server) handleError(err error) error {
	e := &models.Error{}
	ok := errors.As(err, &e)

	if ok {
		switch e.Status {
		case http.StatusInternalServerError:
			s.log.Error(e.Msg)
			return status.Error(codes.Internal, internalError)
		case http.StatusBadRequest:
			return status.Error(codes.InvalidArgument, e.Msg)
		case http.StatusNotFound:
			return status.Error(codes.NotFound, e.Msg)
		}
	}

	s.log.Error("unexpected error", slog.String("error", err.Error()))
	return e
}

func (s *server) ctxWithID(ctx context.Context) context.Context {
	return context.WithValue(ctx, models.ID, uuid.New().String())
}
