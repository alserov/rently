package notifications

import (
	"context"
	"github.com/alserov/rently/user/internal/config"
	"github.com/alserov/rently/user/internal/utils/broker"
)

type Notifier interface {
	Registration(ctx context.Context, email string) error
	Login(email string) error
}

func NewNotifier(producer broker.Producer, topics config.Topics) Notifier {
	return &notifier{producer: producer, topics: topics}
}

const (
	REGISTRATION_ID = "0"
)

type notifier struct {
	producer broker.Producer

	topics config.Topics
}

func (n notifier) Login(email string) error {
	//TODO implement me
	panic("implement me")
}

func (n notifier) Registration(ctx context.Context, email string) error {
	return n.producer.Produce(ctx, email, REGISTRATION_ID, n.topics.Email)
}
