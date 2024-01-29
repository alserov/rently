package service

import (
	"context"
	repomock "github.com/alserov/rently/user/internal/db/mocks"
	"github.com/alserov/rently/user/internal/models"
	notmock "github.com/alserov/rently/user/internal/notifications/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestService_Register(t *testing.T) {
	crtl := gomock.NewController(t)
	defer crtl.Finish()

	req := models.RegisterReq{
		Username:       "user",
		Password:       "password",
		Email:          "email",
		PassportNumber: "123414124",
		PaymentSource:  "23123 13123 13123 13123",
		PhoneNumber:    "23424234234",
	}

	os.Setenv("SECRET_KEY", "secret")

	repo := repomock.NewMockRepository(crtl)
	repo.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	notif := notmock.NewMockNotifier(crtl)
	notif.EXPECT().
		Registration(req.Email).
		Return(nil).
		Times(1)

	s := NewService(Params{
		Repo:     repo,
		Notifier: notif,
	}, nil)

	res, err := s.Register(context.Background(), req)
	require.NoError(t, err)

	require.NotEmpty(t, res.UUID)
	require.NotEmpty(t, res.Token)
}

func TestService_Login(t *testing.T) {

}
