package service

import (
	"context"
	"github.com/alserov/rently/car/internal/utils/payment"

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
	rentEnd := rentStart.Add(time.Hour * 24)
	rentUUID, err := s.CreateRent(context.Background(), models.CreateRentReq{
		RentStart:      &rentStart,
		RentEnd:        &rentEnd,
		PaymentSource:  paymentSource,
		CarPricePerDay: 10,
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
