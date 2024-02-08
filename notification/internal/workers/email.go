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
		msgs: chReg,
	}
}

const (
	sendMailWorkersAmount = 3
)

type worker struct {
	log   log.Logger
	email email.Emailer
	msgs  <-chan amqp.Delivery
}

func (w *worker) Start() {
	w.log.Debug("starting email worker", slog.Int("sendMailWorkers", sendMailWorkersAmount))

	chErr := make(chan error, 1)

	for i := 0; i < sendMailWorkersAmount; i++ {
		go func() {
			for msg := range w.msgs {
				var mail string
				if err := json.Unmarshal(msg.Body, &mail); err != nil {
					chErr <- fmt.Errorf("failed to unmarshal message: %w", err)
				}

				if err := w.email.Send(email.MessageTypeStringToInt(msg.CorrelationId), mail); err != nil {
					chErr <- fmt.Errorf("failed to send email: %v", err)
				}
			}
		}()
	}

	for e := range chErr {
		w.log.Error("failed to send email", slog.String("error", e.Error()))
	}
}
