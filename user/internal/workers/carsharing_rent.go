package workers

import (
	"context"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/user/internal/db"
	"github.com/alserov/rently/user/internal/log"
	"github.com/alserov/rently/user/internal/utils/broker"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

func NewRentReminder() Actor {
	return &rentReminder{}
}

type rentReminder struct {
	log log.Logger
	carsharingClient carsharing.CarsClient
	repo             db.Repository
	producer         broker.Producer
}

func (r rentReminder) Notify() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	tomorrow := time.Now().Add(time.Hour * 24)

	rents, err := r.carsharingClient.GetRentStartingTomorrow(ctx, &carsharing.GetRentStartingTomorrowReq{
		StartingOn: &timestamppb.Timestamp{
			Seconds: tomorrow.Unix(),
			Nanos:   int32(tomorrow.Nanosecond()),
		},
	})
	if err != nil {
		r.log.Error("failed to get rents from carsharing service", slog.String("error", err.Error()))
	}

	for _, rent := range rents.RentsInfo {
		var user db.EmailInfo
		user, err = r.repo.GetUserByUUID(ctx, rent.UserUUID)
		if err != nil {
			r.log.Error("failed to get user from db", slog.String("error", err.Error()))
		}
		if err = r.producer.Produce(ProducerMessage{
			Username:  user.Username,
			Email:     user.Email,
			RentStart: rent.RentStart.AsTime(),
			RentEnd:   rent.RentEnd.AsTime(),
		}); err != nil {
			r.log.Error("failed to produce notifying message", slog.String("error", err.Error()))
		}
	}

	cancel()
}

type ProducerMessage struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RentStart time.Time `json:"rentStart"`
	RentEnd   time.Time `json:"rentEnd"`
}