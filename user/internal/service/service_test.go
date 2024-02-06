package service

import (
	"context"
	"errors"
	"github.com/alserov/rently/user/internal/db"
	repomock "github.com/alserov/rently/user/internal/db/mocks"
	"github.com/alserov/rently/user/internal/models"
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
	repo.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(db.LoginInfo{}, errors.New("user not found")).
		Times(1)

	//notif := notmock.NewMockNotifier(crtl)
	//notif.EXPECT().
	//	Registration(req.Email).
	//	Return(nil).
	//	Times(1)

	s := NewService(Params{
		Repo: repo,
		//Notifier: notif,
	})

	res, err := s.Register(context.Background(), req)
	require.NoError(t, err)

	require.NotEmpty(t, res.UUID)
	require.NotEmpty(t, res.Token)
}

func TestService_Login(t *testing.T) {
	crtl := gomock.NewController(t)
	defer crtl.Finish()

	req := models.LoginReq{
		Password: "password",
		Email:    "email",
	}

	hashedPassword, err := hash(req.Password)
	require.NoError(t, err)

	repo := repomock.NewMockRepository(crtl)
	repo.EXPECT().
		Login(gomock.Any(), req.Email).
		Return(db.LoginInfo{
			UUID:     "uuid",
			Email:    "email@mail.com",
			Password: hashedPassword,
		}, nil).
		Times(1)

	//notif := notmock.NewMockNotifier(crtl)

	s := NewService(Params{
		Repo: repo,
		//Notifier: notif,
	})

	token, err := s.Login(context.Background(), req)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
