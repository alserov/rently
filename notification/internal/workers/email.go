package workers

import (
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/notifications/internal/email"
	"github.com/alserov/rently/notifications/internal/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
)

func NewEmailWorker(smtpPort int, smtpHost string, chReg <-chan amqp.Delivery) Worker {
	return &worker{
		log: log.GetLogger(),
		email: email.NewEmailer(email.Params{
			From:     os.Getenv("EMAIL"),
			Password: os.Getenv("EMAIL_PASSWORD"),
			Sender:   os.Getenv("EMAIL_SENDER"),
			SmtpHost: smtpHost,
			SmtpPort: smtpPort,
		}),
		registerCh: chReg,
	}
}

const (
	registerWorkersAmount    = 2
	readMessageWorkersAmount = 3
)

type worker struct {
	log        log.Logger
	email      email.Emailer
	registerCh <-chan amqp.Delivery
}

const (
	REGISTER_ID = "0"
)

func (w *worker) Start() {
	w.log.Debug("starting email worker", slog.Int("registerWorkers", registerWorkersAmount), slog.Int("readMessageWorkers", readMessageWorkersAmount))

	chErr := make(chan error, 1)
	chRegister := make(chan amqp.Delivery, registerWorkersAmount)

	for i := 0; i < readMessageWorkersAmount; i++ {
		go func() {
			for msg := range w.registerCh {
				switch msg.CorrelationId {
				case REGISTER_ID:
					w.log.Debug("received message", slog.Any("value", msg))
					chRegister <- msg
				default:
					w.log.Warn("received message with unknown id", slog.String("id", msg.CorrelationId))
				}
			}
		}()
	}

	for i := 0; i < registerWorkersAmount; i++ {
		go func() {
			for msg := range chRegister {
				var email string
				if err := json.Unmarshal(msg.Body, &email); err != nil {
					chErr <- fmt.Errorf("failed to unmarshal message: %w", err)
				}

				if err := w.register(email); err != nil {
					chErr <- fmt.Errorf("failed to notify about registration: %w", err)
				}
			}
		}()
	}

	for e := range chErr {
		w.log.Error("failed to send email", slog.String("error", e.Error()))
	}
}

func (w *worker) register(email string) error {
	if err := w.email.Registration(email); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
