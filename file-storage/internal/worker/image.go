package worker

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/alserov/file-storage/internal/utils/broker"
	"github.com/alserov/file-storage/internal/utils/files"
	"log/slog"
)

func NewImageWorker(brokerAddr string, topics broker.Topics, log *slog.Logger) Worker {
	return &imageWorker{
		log:            log,
		brokerAddr:     brokerAddr,
		topics:         topics,
		consumerConfig: sarama.NewConfig(),
		filer:          files.NewFiler(files.RelativeImageDir),
	}
}

type uuidMessage struct {
	UUID string `json:"uuid"`
}

type fileMessage struct {
	Value []byte `json:"value"`
	UUID  string `json:"uuid"`
	Idx   int    `json:"idx"`
}

type imageWorker struct {
	log *slog.Logger

	topics     broker.Topics
	brokerAddr string

	consumerConfig *sarama.Config

	filer files.Filer
}

func (iw *imageWorker) MustStart(ctx context.Context) {
	go iw.save(ctx, iw.topics.SaveImages)
	go iw.delete(ctx, iw.topics.DeleteImages)

	<-ctx.Done()
}

const workersAmount = 3

func (iw *imageWorker) delete(ctx context.Context, topic string) {
	messages := broker.Consume[uuidMessage](workersAmount, iw.brokerAddr, topic, iw.consumerConfig, iw.log)

	for i := 0; i < workersAmount; i++ {
		go func() {
			for msg := range messages {
				err := iw.filer.Delete(msg.UUID)
				if err != nil {
					iw.log.Error("failed to delete image", slog.String("error", err.Error()))
				}
			}
		}()
	}

	<-ctx.Done()
}

func (iw *imageWorker) save(ctx context.Context, topic string) {
	messages := broker.Consume[fileMessage](workersAmount, iw.brokerAddr, topic, iw.consumerConfig, iw.log)
	for i := 0; i < workersAmount; i++ {
		go func() {
			for msg := range messages {
				err := iw.filer.Save(msg.Value, msg.UUID, msg.Idx)
				if err != nil {
					iw.log.Error("failed to save image", slog.String("error", err.Error()))
				}
			}
		}()
	}

	<-ctx.Done()
}
