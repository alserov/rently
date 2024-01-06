package server

import (
	"context"
	"errors"
	"fmt"

	repomock "github.com/alserov/rently/car/internal/db/mocks"
	repomodels "github.com/alserov/rently/car/internal/db/models"
	metrmock "github.com/alserov/rently/car/internal/metrics/mocks"
	models "github.com/alserov/rently/car/internal/models"
	clientmocks "github.com/alserov/rently/car/internal/server/mocks"
	"github.com/alserov/rently/car/internal/service"
	"github.com/alserov/rently/car/internal/utils/clients"
	"github.com/alserov/rently/car/internal/utils/log"
	"github.com/alserov/rently/proto/gen/car"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"testing"
	"time"
)

func TestServer_CheckRent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	start := time.Now()
	end := start.Add(time.Hour * 24 * 3)

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().
		CheckRent(gomock.Eq(context.Background()), gomock.Any()).
		Return(repomodels.Rent{
			CarUUID:   "uuid",
			RentPrice: 100,
			RentStart: &start,
			RentEnd:   &end,
		}, nil)

	s := service.NewService(service.Params{
		Repo: repo,
	})

	serv := newServer(s, nil, log.MustSetup("local"))

	res, err := serv.CheckRent(context.Background(), &car.CheckRentReq{
		RentUUID: "uuid",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.CarUUID)
}

func TestServer_CreateRent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().
		CreateRent(gomock.Eq(context.Background()), gomock.Any()).
		Return(nil).
		Times(1)
	repo.EXPECT().
		GetCarPriceIfAvailable(gomock.Eq(context.Background()), gomock.Any()).
		Return(float32(100), nil).
		Times(1)

	metr := metrmock.NewMockMetrics(ctrl)
	metr.EXPECT().
		IncreaseActiveRentsAmount().Times(1)

	s := service.NewService(service.Params{
		Repo:    repo,
		Metrics: metr,
	})

	serv := newServer(s, nil, log.MustSetup("local"))

	start := time.Now()
	end := start.Add(time.Hour * 24)
	fmt.Println(int32(end.Nanosecond()), end.Unix(), int32(start.Nanosecond()), start.Unix())

	res, err := serv.CreateRent(context.Background(), &car.CreateRentReq{
		CarUUID:        "bfe2f3b8-15dc-494e-a9cf-4d6fc538b75e",
		PassportNumber: "2342",
		PaymentSource:  "tok_visa",
		PhoneNumber:    "79787777777",
		RentEnd: &timestamppb.Timestamp{
			Nanos:   int32(end.Nanosecond()),
			Seconds: end.Unix(),
		},
		RentStart: &timestamppb.Timestamp{
			Nanos:   int32(start.Nanosecond()),
			Seconds: start.Unix(),
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, res.RentUUID)
}

func TestServer_GetAvailableCars(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cars := []repomodels.Car{
		{
			UUID: "uuid1",
		},
		{
			UUID: "uuid2",
		},
	}

	cl := clientmocks.NewMockFileStorageClient(ctrl)
	cl.EXPECT().
		GetLinks(gomock.Eq(context.Background()), gomock.Any()).
		Return(&fstorage.GetLinksRes{Links: []string{"link1", "link2"}}, nil).
		Times(len(cars))

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().
		GetAvailableCars(gomock.Eq(context.Background()), gomock.Any()).
		Return(cars, nil).
		Times(1)

	s := newServer(service.NewService(service.Params{
		BrokerConsumerConfig: nil,
		BrokerProducer:       nil,
		Repo:                 repo,
		Metrics:              nil,
	}), &clients.Clients{
		FileStorage: cl,
	}, nil)

	start, end := time.Now(), time.Now().Add(time.Hour*24)

	res, err := s.GetAvailableCars(context.Background(), &car.GetAvailableCarsReq{
		Start: &timestamppb.Timestamp{
			Seconds: int64(start.Second()),
			Nanos:   int32(start.Nanosecond()),
		},
		End: &timestamppb.Timestamp{
			Seconds: int64(end.Second()),
			Nanos:   int32(end.Nanosecond()),
		},
	})
	require.NoError(t, err)

	require.NotEmpty(t, res)
	for _, c := range res.Cars {
		require.NotEmpty(t, c.Images)
	}
}

func TestServer_GetAvailableCars_with_repo_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().
		GetAvailableCars(gomock.Eq(context.Background()), gomock.Any()).
		Return(nil, &models.Error{Code: http.StatusInternalServerError, Msg: "server error"}).
		Times(1)

	s := newServer(service.NewService(service.Params{
		BrokerConsumerConfig: nil,
		BrokerProducer:       nil,
		Repo:                 repo,
		Metrics:              nil,
	}), nil, nil)

	start, end := time.Now(), time.Now().Add(time.Hour*24)

	res, err := s.GetAvailableCars(context.Background(), &car.GetAvailableCarsReq{
		Start: &timestamppb.Timestamp{
			Seconds: int64(start.Second()),
			Nanos:   int32(start.Nanosecond()),
		},
		End: &timestamppb.Timestamp{
			Seconds: int64(end.Second()),
			Nanos:   int32(end.Nanosecond()),
		},
	})

	e := &models.Error{}
	ok := errors.As(err, &e)

	require.True(t, ok)
	require.Error(t, err)
	require.Empty(t, res)
	require.Equal(t, e.Code, http.StatusInternalServerError)
}

func TestServer_GetAvailableCars_with_client_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cars := []repomodels.Car{
		{
			UUID: "uuid1",
		},
		{
			UUID: "uuid2",
		},
	}

	cl := clientmocks.NewMockFileStorageClient(ctrl)
	cl.EXPECT().
		GetLinks(gomock.Eq(context.Background()), gomock.Any()).
		Return(nil, errors.New("other service error")).
		Times(len(cars))

	repo := repomock.NewMockRepository(ctrl)
	repo.EXPECT().
		GetAvailableCars(gomock.Eq(context.Background()), gomock.Any()).
		Return(cars, nil).
		Times(1)

	srvc := service.NewService(service.Params{
		BrokerConsumerConfig: nil,
		BrokerProducer:       nil,
		Repo:                 repo,
		Metrics:              nil,
	})

	s := newServer(srvc, &clients.Clients{
		FileStorage: cl,
	}, log.MustSetup("local"))

	start, end := time.Now(), time.Now().Add(time.Hour*24)

	res, err := s.GetAvailableCars(context.Background(), &car.GetAvailableCarsReq{
		Start: &timestamppb.Timestamp{
			Seconds: int64(start.Second()),
			Nanos:   int32(start.Nanosecond()),
		},
		End: &timestamppb.Timestamp{
			Seconds: int64(end.Second()),
			Nanos:   int32(end.Nanosecond()),
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, res)
	for _, c := range res.Cars {
		require.Empty(t, c.Images)
	}
}
