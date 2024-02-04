package service

import (
	"context"
	repomock "github.com/alserov/rently/carsharing/internal/db/mocks"
	"github.com/alserov/rently/carsharing/internal/models"
	storagemock "github.com/alserov/rently/carsharing/internal/storage/mocks"
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
	repo.EXPECT().DeleteCar(gomock.Eq(context.Background()), gomock.Eq(uuid)).Return([]string{"image_url"}, nil).Times(1)

	storage := storagemock.NewMockImageStorage(ctrl)
	storage.EXPECT().Delete(gomock.Eq(context.Background()), gomock.Eq("image_url")).Return(nil).Times(1)

	s := NewService(Params{
		Repo:         repo,
		ImageStorage: storage,
	})

	err := s.DeleteCar(context.Background(), uuid)
	require.NoError(t, err)
}

func TestNewService_GetCarByUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().GetCarByUUID(gomock.Eq(context.Background()), "uuid").Return(models.Car{}, nil).Times(1)

	s := NewService(Params{
		Repo: repo,
	})

	_, err := s.GetCarByUUID(context.Background(), "uuid")
	require.NoError(t, err)
}

func TestService_CreateCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().CreateCar(gomock.Eq(context.Background()), gomock.Any()).Return(nil).Times(1)

	storage := storagemock.NewMockImageStorage(ctrl)
	storage.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	s := NewService(Params{
		Repo:         repo,
		ImageStorage: storage,
	})

	err := s.CreateCar(context.Background(), models.Car{
		UUID: "uuid",
	}, [][]byte{[]byte("bytes")}, []byte("bytes"))
	require.NoError(t, err)
}

func TestService_CreateRent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().CreateRentTx(gomock.Eq(context.Background()), gomock.Any()).Return(100, nil, nil).Times(1)

	s := NewService(Params{})

	rentStart := time.Now()
	rentEnd := rentStart.Add(time.Hour * 24 * 3)
	rentUUID, err := s.CreateRent(context.Background(), models.CreateRentReq{
		RentStart:      rentStart,
		RentEnd:        rentEnd,
		PaymentSource:  paymentSource,
		CarUUID:        "uuid",
		PassportNumber: "passport_number",
		PhoneNumber:    "458934534524",
	})
	require.NoError(t, err)
	require.NotEmpty(t, rentUUID)
}
