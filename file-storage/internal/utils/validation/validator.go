package validation

import (
	"fmt"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validator interface {
	ValidateGetLinksReq(req *fstorage.GetLinksReq) error
}

func NewValidator() Validator {
	return &validator{}
}

type validator struct {
}

const (
	ERR_EMPTY = "can not be empty"
)

func (v validator) ValidateGetLinksReq(req *fstorage.GetLinksReq) error {
	if req.GetUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("uuid %s", ERR_EMPTY))
	}

	return nil
}
