package server

import (
	"context"
	"errors"
	repomock "github.com/alserov/rently/car/internal/db/mocks"
	"github.com/alserov/rently/car/internal/db/models"
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

func TestServer_GetAvailableCars(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cars := []models.Car{
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
		Return(nil, &Error{Code: http.StatusInternalServerError, Msg: "server error"}).
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

	e := &Error{}
	ok := errors.As(err, &e)

	require.True(t, ok)
	require.Error(t, err)
	require.Empty(t, res)
	require.Equal(t, e.Code, http.StatusInternalServerError)
}

func TestServer_GetAvailableCars_with_client_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cars := []models.Car{
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
