package service

import (
	"context"
	"errors"
	"github.com/alserov/rently/car/internal/utils/payment"
	"io/ioutil"
	"os"

	repomock "github.com/alserov/rently/car/internal/db/mocks"
	repository "github.com/alserov/rently/car/internal/db/models"
	metrmock "github.com/alserov/rently/car/internal/metrics/mocks"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const paymentSource = "tok_visa"

func TestService_DeleteCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuid := "uuid"

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().DeleteCar(gomock.Eq(context.Background()), gomock.Eq(uuid)).Return(nil).Times(1)

	s := NewService(repo, nil, nil)

	err := os.MkdirAll("images/uuid", 0644)
	f, err := os.OpenFile("images/uuid/0.txt", os.O_CREATE, 0644)
	require.NoError(t, err)
	defer os.RemoveAll("images")

	_, err = f.Write([]byte("image"))
	require.NoError(t, err)
	f.Close()

	err = s.DeleteCar(context.Background(), uuid)
	require.NoError(t, err)
}

func TestService_DeleteCar_with_repo_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuid := "uuid"

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().DeleteCar(gomock.Eq(context.Background()), gomock.Eq(uuid)).Return(errors.New("not found")).Times(1)

	s := NewService(repo, nil, nil)

	err := os.MkdirAll("images/uuid", 0644)
	f, err := os.OpenFile("images/uuid/0.txt", os.O_CREATE, 0644)
	require.NoError(t, err)
	defer os.RemoveAll("images")

	_, err = f.Write([]byte("image"))
	require.NoError(t, err)
	f.Close()

	err = s.DeleteCar(context.Background(), uuid)
	require.Error(t, err)
}

func TestService_DeleteCar_with_file_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuid := "invalid_uuid"

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().DeleteCar(gomock.Eq(context.Background()), gomock.Eq(uuid)).Return(nil).AnyTimes()

	s := NewService(repo, nil, nil)

	err := os.MkdirAll("images/uuid", 0644)
	f, err := os.OpenFile("images/uuid/0.txt", os.O_CREATE, 0644)
	require.NoError(t, err)
	defer os.RemoveAll("images")

	_, err = f.Write([]byte("image"))
	require.NoError(t, err)
	f.Close()

	err = s.DeleteCar(context.Background(), uuid)
	require.Error(t, err)
}

func TestService_CreateCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().CreateCar(gomock.Eq(context.Background()), gomock.Any()).Return(nil).Times(1)

	s := NewService(repo, nil, nil)

	f, err := os.OpenFile("image", os.O_CREATE, 0644)
	require.NoError(t, err)
	defer os.Remove("image")
	defer os.RemoveAll("images")
	defer f.Close()

	_, err = f.Write([]byte("image"))
	require.NoError(t, err)

	b, err := ioutil.ReadFile("image")
	require.NoError(t, err)

	err = s.CreateCar(context.Background(), models.Car{
		UUID:   "uuid",
		Images: [][]byte{b, b},
	})
	require.NoError(t, err)
}

func TestService_CreateRent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().CreateRent(gomock.Eq(context.Background()), gomock.Any()).Return(nil).Times(1)
	repo.EXPECT().CheckIfCarAvailable(gomock.Eq(context.Background()), gomock.Any()).Return(nil).Times(1)

	metr := metrmock.NewMockMetrics(ctrl)
	metr.EXPECT().IncreaseActiveRentsAmount().Times(1)

	s := NewService(repo, metr, nil)

	rentStart := time.Now()
	rentEnd := rentStart.Add(time.Hour * 24 * 3)
	rentUUID, err := s.CreateRent(context.Background(), models.CreateRentReq{
		RentStart:      &rentStart,
		RentEnd:        &rentEnd,
		PaymentSource:  paymentSource,
		CarPricePerDay: 50,
	})
	require.NoError(t, err)
	require.NotEmpty(t, rentUUID)
}

func TestService_CancelRent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const rentPrice = 10000

	p := payment.NewPayer("sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh")
	chargeID, err := p.Debit(paymentSource, rentPrice)
	require.NoError(t, err)

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().
		CancelRent(gomock.Eq(context.Background()), gomock.Any()).
		Return(repository.CancelRentInfo{ChargeID: chargeID, RentPrice: rentPrice}, nil).
		Times(1)

	metr := metrmock.NewMockMetrics(ctrl)
	metr.EXPECT().DecreaseActiveRentsAmount().Times(1)

	s := NewService(repo, metr, nil)

	err = s.CancelRent(context.Background(), "uuid")
	require.NoError(t, err)
}
