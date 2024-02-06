package workers

import (
	"context"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/user/internal/config"
	"github.com/alserov/rently/user/internal/db"
	"github.com/alserov/rently/user/internal/log"
	"github.com/alserov/rently/user/internal/utils/broker"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

func NewRentNotifier() Actor {
	return &rentReminder{}
}

const (
	RENT_NOTIFIER_ID = "10"
)

type rentReminder struct {
	log              log.Logger
	carsharingClient carsharing.CarsClient
	repo             db.Repository
	producer         broker.Producer
	topics           config.Topics
}

func (r *rentReminder) Notify() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tomorrow := time.Now().Add(time.Hour * 24)

	rents, err := r.carsharingClient.GetRentStartingOnDate(ctx, &carsharing.GetRentStartingOnDateReq{
		StartingOn: &timestamppb.Timestamp{
			Seconds: tomorrow.Unix(),
			Nanos:   int32(tomorrow.Nanosecond()),
		},
	})
	if err != nil {
		r.log.Error("failed to get rents from carsharing service", slog.String("error", err.Error()))
	}

	for _, rent := range rents.RentsInfo {
		var user db.EmailNotificationsInfo
		user, err = r.repo.GetUserByUUID(ctx, rent.UserUUID)
		if err != nil {
			r.log.Error("failed to get user from db", slog.String("error", err.Error()))
		}

		err = r.producer.Produce(ctx, ProducerMessage{
			Username:  user.Username,
			Email:     user.Email,
			RentStart: rent.RentStart.AsTime(),
			RentEnd:   rent.RentEnd.AsTime(),
		}, RENT_NOTIFIER_ID, r.topics.Email)
		if err != nil {
			r.log.Error("failed to produce notifying message", slog.String("error", err.Error()))
		}
	}
}

type ProducerMessage struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RentStart time.Time `json:"rentStart"`
	RentEnd   time.Time `json:"rentEnd"`
}
