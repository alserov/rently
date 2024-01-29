package email

import (
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/notifications/internal/workers"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
)

func NewEmailWorker(port int, host string) workers.Worker {
	return &worker{
		email: NewEmailer(Params{
			From:     os.Getenv("EMAIL"),
			Password: os.Getenv("EMAIL_PASSWORD"),
			SmtpHost: host,
			SmtpPort: port,
		}),
	}
}

const (
	workersAmount = 5
)

type worker struct {
	log *slog.Logger

	email Emailer

	msgs <-chan amqp.Delivery
}

func (w *worker) Start() {
	go w.register(w.msgs)

	select {}
}

func (w *worker) register(msgs <-chan amqp.Delivery) {
	chErr := make(chan error, workersAmount)

	for i := 0; i < workersAmount; i++ {
		go func(msgs <-chan amqp.Delivery) {
			for m := range msgs {
				var email string
				if err := json.Unmarshal(m.Body, &email); err != nil {
					chErr <- fmt.Errorf("failed to unmarshal message: %v", err)
				}
				if err := w.email.Registration(email); err != nil {
					chErr <- fmt.Errorf("failed to send email: %v", err)
				}
			}
		}(msgs)
	}

	go func() {
		for e := range chErr {
			w.log.Error("failed to send email", slog.String("error", e.Error()))
		}
	}()
}
