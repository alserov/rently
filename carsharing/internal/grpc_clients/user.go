package grpc_clients

import (
	"context"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/alserov/rently/proto/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type UserClient interface {
	GetInfo(ctx context.Context, uuid string) (models.UserInfo, error)
}

func NewUserClient(cl user.UserClient) UserClient {
	return &userClient{}
}

type userClient struct {
	cl user.UserClient
}

func (u userClient) GetInfo(ctx context.Context, uuid string) (models.UserInfo, error) {
	info, err := u.cl.GetInfoForRent(ctx, &user.GetInfoReq{UUID: uuid})
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Internal:
			return models.UserInfo{}, &models.Error{
				Msg:    st.Message(),
				Status: http.StatusInternalServerError,
			}
		case codes.InvalidArgument:
			return models.UserInfo{}, &models.Error{
				Msg:    st.Message(),
				Status: http.StatusBadRequest,
			}
		}
	}

	return models.UserInfo{
		PassportNumber: info.PassportNumber,
		PhoneNumber:    info.PhoneNumber,
	}, nil
}
