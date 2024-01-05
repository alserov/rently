package server

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

const internalError = "internal error"

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}

func handleError(err error, log *slog.Logger) error {
	var error Error
	ok := errors.As(err, &error)

	if ok {
		switch error.Code {
		case http.StatusInternalServerError:
			return status.Error(codes.Internal, internalError)
		case http.StatusBadRequest:
			return status.Error(codes.InvalidArgument, error.Msg)
		case http.StatusNotFound:
			return status.Error(codes.NotFound, error.Msg)
		}
	}

	log.Error("unexpected error", slog.String("error", err.Error()))
	return error
}
