package workers

import (
	"context"
	"fmt"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/alserov/rently/user/internal/config"
	"github.com/alserov/rently/user/internal/db"
	"github.com/alserov/rently/user/internal/utils/broker"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type NotifierParams struct {
	CarsharingClient carsharing.CarsClient
	Repo             db.Repository
	Producer         broker.Producer
	Topics           config.Topics
}

func NewRentNotifier(p NotifierParams) Actor {
	return &rentReminder{
		carsharingClient: p.CarsharingClient,
		repo:             p.Repo,
		producer:         p.Producer,
		topics:           p.Topics,
	}
}

const (
	RENT_NOTIFIER_ID = "10"
)

type rentReminder struct {
	carsharingClient carsharing.CarsClient
	repo             db.Repository
	producer         broker.Producer
	topics           config.Topics
}

func (r *rentReminder) Action() error {
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
		return fmt.Errorf("failed to get rents from carsharing service: %w", err)
	}

	for _, rent := range rents.RentsInfo {
		var user db.EmailNotificationsInfo
		user, err = r.repo.GetUserByUUID(ctx, rent.UserUUID)
		if err != nil {
			return fmt.Errorf("failed to get user from db: %w", err)
		}

		err = r.producer.Produce(ctx, ProducerMessage{
			Username:  user.Username,
			Email:     user.Email,
			RentStart: rent.RentStart.AsTime(),
			RentEnd:   rent.RentEnd.AsTime(),
		}, RENT_NOTIFIER_ID, r.topics.Email)
		if err != nil {
			return fmt.Errorf("failed to produce notifying message: %w", err)
		}
	}

	return nil
}

type ProducerMessage struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RentStart time.Time `json:"rentStart"`
	RentEnd   time.Time `json:"rentEnd"`
}
