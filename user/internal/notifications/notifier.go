package notifications

import (
	"context"
	"github.com/alserov/rently/user/internal/config"
	"github.com/alserov/rently/user/internal/utils/broker"
)

type Notifier interface {
	Registration(ctx context.Context, email string) error
	Login(ctx context.Context, email string) error
}

func NewNotifier(producer broker.Producer, topics config.Topics) Notifier {
	return &notifier{producer: producer, topics: topics}
}

const (
	REGISTRATION_ID = "0"
	LOGIN_ID        = "1"
)

type notifier struct {
	producer broker.Producer

	topics config.Topics
}

func (n notifier) Login(ctx context.Context, email string) error {
	return n.producer.Produce(ctx, email, LOGIN_ID, n.topics.Email)
}

func (n notifier) Registration(ctx context.Context, email string) error {
	return n.producer.Produce(ctx, email, REGISTRATION_ID, n.topics.Email)
}
