package clients

import (
	"context"
	"github.com/alserov/rently/carsharing/internal/models"
	"github.com/alserov/rently/proto/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type UserClient interface {
	GetPassportAndPhone(ctx context.Context, token string) (models.UserInfo, error)
}

func NewUserClient(cl user.UserClient) UserClient {
	return &userClient{cl: cl}
}

type userClient struct {
	cl user.UserClient
}

func (u userClient) GetPassportAndPhone(ctx context.Context, token string) (models.UserInfo, error) {
	info, err := u.cl.GetInfoForRent(ctx, &user.GetRentInfoReq{Token: token})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return models.UserInfo{}, &models.Error{
					Msg:    st.Message(),
					Status: http.StatusBadRequest,
				}
			default:
				return models.UserInfo{}, &models.Error{
					Msg:    st.Message(),
					Status: http.StatusInternalServerError,
				}
			}
		}
	}

	return models.UserInfo{
		PassportNumber: info.PassportNumber,
		PhoneNumber:    info.PhoneNumber,
	}, nil
}
