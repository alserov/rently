package service

import (
	"context"
	repomock "github.com/alserov/rently/car/internal/db/mocks"
	metrmock "github.com/alserov/rently/car/internal/metrics/mocks"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
		PaymentSource:  "tok_visa",
		CarPricePerDay: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, rentUUID)
}
